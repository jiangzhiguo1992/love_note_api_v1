package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddPromiseBreak
func AddPromiseBreak(uid, cid int64, pb *entity.PromiseBreak) (*entity.PromiseBreak, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pb == nil {
		return nil, errors.New("nil_promise_break")
	} else if pb.PromiseId <= 0 {
		return nil, errors.New("nil_promise")
	} else if pb.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len([]rune(pb.ContentText)) > GetLimit().PromiseBreakContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 数据检查
	p, err := mysql.GetPromiseById(pb.PromiseId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_promise")
	} else if p.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// mysql
	pb.UserId = uid
	pb.CoupleId = cid
	pb, err = mysql.AddPromiseBreak(pb)
	if pb == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// promise
		UpdatePromise(uid, cid, p, false)
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_PROMISE, p.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, p.Id, "push_title_note_update", pb.ContentText, entity.PUSH_TYPE_NOTE_PROMISE_BREAK)
	}()
	return pb, err
}

// DelPromiseBreak
func DelPromiseBreak(uid, cid, pbid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if pbid <= 0 {
		return errors.New("nil_promise_break")
	}
	// 旧数据检查
	pb, err := mysql.GetPromiseBreakById(pbid)
	if err != nil {
		return err
	} else if pb == nil {
		return errors.New("nil_promise_break")
	} else if pb.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	p, err := mysql.GetPromiseById(pb.PromiseId)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("nil_promise")
	} else if p.CoupleId != cid {
		return errors.New("db_query_refuse")
	}
	// mysql
	err = mysql.DelPromiseBreak(pb)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		// promise
		UpdatePromise(uid, cid, p, false)
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_PROMISE, p.Id)
		AddTrends(trends)
	}()
	return err
}

// GetPromiseBreakListByCouplePromise
func GetPromiseBreakListByCouplePromise(uid, cid, pid int64, page int) ([]*entity.PromiseBreak, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pid <= 0 {
		return nil, errors.New("nil_promise")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().PromiseBreak
	offset := page * limit
	list, err := mysql.GetPromiseBreakListByCouplePromise(cid, pid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_promise_break")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 没有trends，promise里有了
	return list, err
}
