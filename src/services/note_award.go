package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddAward
func AddAward(uid, cid int64, a *entity.Award) (*entity.Award, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil {
		return nil, errors.New("nil_award")
	} else if a.AwardRuleId <= 0 {
		return nil, errors.New("nil_award_rule")
	} else if a.HappenId <= 0 {
		return nil, errors.New("award_nil_happen_id")
	} else if a.AwardRuleId <= 0 {
		return nil, errors.New("award_nil_rule_id")
	} else if a.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(a.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(a.ContentText)) > GetLimit().AwardContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 数据检查
	ar, _ := mysql.GetAwardRuleById(a.AwardRuleId)
	if ar == nil || ar.CoupleId != cid {
		return nil, errors.New("award_nil_rule_id")
	}
	// mysql
	a.UserId = uid
	a.CoupleId = cid
	a.ScoreChange = ar.Score
	a, err := mysql.AddAward(a)
	if a == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// rule
		ar.UseCount = ar.UseCount + 1
		mysql.UpdateAwardRule(ar)
		// score
		UpdateAwardScore(uid, cid, a, true)
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_AWARD, a.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, a.Id, "push_title_note_update", a.ContentText, entity.PUSH_TYPE_NOTE_AWARD)
	}()
	return a, err
}

// DelAward
func DelAward(uid, cid, aid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if aid <= 0 {
		return errors.New("nil_award")
	}
	// 旧数据检查
	a, err := mysql.GetAwardById(aid)
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("nil_award")
	} else if a.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelAward(a)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		// rule
		ar, _ := mysql.GetAwardRuleById(a.AwardRuleId)
		if ar != nil && ar.Id > 0 {
			ar.UseCount = ar.UseCount - 1
			mysql.UpdateAwardRule(ar)
		}
		// score
		UpdateAwardScore(uid, cid, a, false)
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_AWARD, aid)
		AddTrends(trends)
	}()
	return err
}

// GetAwardListByUserCouple
func GetAwardListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Award, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Award
	offset := page * limit
	var list []*entity.Award
	var err error
	if suid > 0 {
		list, err = mysql.GetAwardListByCoupleHappenUser(cid, suid, offset, limit)
	} else {
		list, err = mysql.GetAwardListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_award")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_AWARD)
		AddTrends(trends)
	}()
	return list, err
}
