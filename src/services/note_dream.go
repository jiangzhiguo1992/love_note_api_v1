package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddDream
func AddDream(uid, cid int64, d *entity.Dream) (*entity.Dream, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if d == nil {
		return nil, errors.New("nil_dream")
	} else if d.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(d.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(d.ContentText)) >= GetLimit().DreamContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// mysql
	d.UserId = uid
	d.CoupleId = cid
	d, err := mysql.AddDream(d)
	if d == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_DREAM, d.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, d.Id, "push_title_note_update", d.ContentText, entity.PUSH_TYPE_NOTE_DREAM)
	}()
	return d, err
}

// DelDream
func DelDream(uid, cid, did int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if did <= 0 {
		return errors.New("nil_dream")
	}
	// 旧数据检查
	d, err := mysql.GetDreamById(did)
	if err != nil {
		return err
	} else if d == nil {
		return errors.New("nil_dream")
	} else if d.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelDream(d)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_DREAM, did)
		AddTrends(trends)
	}()
	return err
}

// UpdateDream
func UpdateDream(uid, cid int64, d *entity.Dream) (*entity.Dream, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if d == nil || d.Id <= 0 {
		return nil, errors.New("nil_dream")
	} else if d.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(d.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(d.ContentText)) >= GetLimit().DreamContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 旧数据检查
	old, err := mysql.GetDreamById(d.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_dream")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.HappenAt = d.HappenAt
	old.ContentText = d.ContentText
	d, err = mysql.UpdateDream(old)
	if d == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_DREAM, d.Id)
		AddTrends(trends)
	}()
	return d, err
}

// GetDreamById
func GetDreamById(uid, cid, did int64) (*entity.Dream, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if did <= 0 {
		return nil, errors.New("nil_dream")
	}
	// mysql
	d, err := mysql.GetDreamById(did)
	if err != nil {
		return nil, err
	} else if d == nil {
		return nil, errors.New("nil_dream")
	} else if d.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_DREAM, did)
		AddTrends(trends)
	}()
	return d, err
}

// GetDreamListByUserCouple
func GetDreamListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Dream, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Dream
	offset := page * limit
	var list []*entity.Dream
	var err error
	if suid > 0 {
		list, err = mysql.GetDreamListByUserCouple(suid, cid, offset, limit)
	} else {
		list, err = mysql.GetDreamListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_dream")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_DREAM)
		AddTrends(trends)
	}()
	return list, err
}
