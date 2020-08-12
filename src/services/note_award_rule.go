package services

import (
	"errors"
	"math"
	"models/entity"
	"models/mysql"
)

// AddAwardRule
func AddAwardRule(uid, cid int64, ar *entity.AwardRule) (*entity.AwardRule, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if ar == nil {
		return nil, errors.New("nil_award_rule")
	} else if len(ar.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(ar.Title)) > GetLimit().AwardRuleTitleLength {
		return nil, errors.New("limit_title_over")
	} else if ar.Score == 0 {
		return nil, errors.New("award_score_nil")
	} else if int(math.Abs(float64(ar.Score))) > GetLimit().AwardRuleScoreMax {
		return nil, errors.New("award_score_over")
	}
	// mysql
	ar.UserId = uid
	ar.CoupleId = cid
	ar, err := mysql.AddAwardRule(ar)
	if ar == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_AWARD_RULE, ar.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, ar.Id, "push_title_note_update", ar.Title, entity.PUSH_TYPE_NOTE_AWARD_RULE)
	}()
	return ar, err
}

// DelAwardRule
func DelAwardRule(uid, cid, arid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if arid <= 0 {
		return errors.New("nil_award_rule")
	}
	// 旧数据检查
	ar, err := mysql.GetAwardRuleById(arid)
	if err != nil {
		return err
	} else if ar == nil {
		return errors.New("nil_award_rule")
	} else if ar.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelAwardRule(ar)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_AWARD_RULE, arid)
		AddTrends(trends)
	}()
	return err
}

// GetAwardRuleListByCouple
func GetAwardRuleListByCouple(uid, cid int64, page int) ([]*entity.AwardRule, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().AwardRule
	offset := page * limit
	list, err := mysql.GetAwardRuleListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_award_rule")
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
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_AWARD_RULE)
		AddTrends(trends)
	}()
	return list, err
}
