package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddPromise
func AddPromise(uid, cid int64, p *entity.Promise) (*entity.Promise, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if p == nil {
		return nil, errors.New("nil_promise")
	} else if p.HappenId == 0 {
		return nil, errors.New("promise_nil_happen_id")
	} else if p.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(p.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(p.ContentText)) > GetLimit().PromiseContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// mysql
	p.UserId = uid
	p.CoupleId = cid
	p, err := mysql.AddPromise(p)
	if p == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_PROMISE, p.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, p.Id, "push_title_note_update", p.ContentText, entity.PUSH_TYPE_NOTE_PROMISE)
	}()
	return p, err
}

// DelPromise
func DelPromise(uid, cid, pid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if pid <= 0 {
		return errors.New("nil_promise")
	}
	// 旧数据检查
	p, err := mysql.GetPromiseById(pid)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("nil_promise")
	} else if p.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelPromise(p)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_PROMISE, pid)
		AddTrends(trends)
	}()
	return err
}

// UpdatePromise
func UpdatePromise(uid, cid int64, p *entity.Promise, self bool) (*entity.Promise, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if p == nil || p.Id <= 0 {
		return nil, errors.New("nil_promise")
	} else if p.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(p.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(p.ContentText)) > GetLimit().PromiseContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 旧数据检查
	old, err := mysql.GetPromiseById(p.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_promise")
	} else if self && old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.HappenAt = p.HappenAt
	old.ContentText = p.ContentText
	old.BreakCount = int(mysql.GetPromiseBreakTotalByCouplePromise(cid, p.Id))
	p, err = mysql.UpdatePromise(old)
	if p == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_PROMISE, p.Id)
		AddTrends(trends)
	}()
	return p, err
}

// GetPromiseById
func GetPromiseById(uid, cid, pid int64) (*entity.Promise, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pid <= 0 {
		return nil, errors.New("nil_promise")
	}
	// mysql
	p, err := mysql.GetPromiseById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_promise")
	} else if p.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_PROMISE, pid)
		AddTrends(trends)
	}()
	return p, err
}

// GetPromiseListByUserCouple
func GetPromiseListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Promise, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Promise
	offset := page * limit
	var list []*entity.Promise
	var err error
	if suid > 0 {
		list, err = mysql.GetPromiseListByCoupleHappenUser(cid, suid, offset, limit)
	} else {
		list, err = mysql.GetPromiseListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_promise")
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
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_PROMISE)
		AddTrends(trends)
	}()
	return list, err
}
