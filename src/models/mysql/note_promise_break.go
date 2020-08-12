package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPromiseBreak
func AddPromiseBreak(pb *entity.PromiseBreak) (*entity.PromiseBreak, error) {
	if len(pb.ContentText) <= 0 {
		pb.ContentText = ""
	}
	pb.Status = entity.STATUS_VISIBLE
	pb.CreateAt = time.Now().Unix()
	pb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_PROMISE_BREAK).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,promise_id=?,happen_at=?,content_text=?").
		Exec(pb.Status, pb.CreateAt, pb.UpdateAt, pb.UserId, pb.CoupleId, pb.PromiseId, pb.HappenAt, pb.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	pb.Id, _ = db.Result().LastInsertId()
	return pb, nil
}

// DelPromiseBreak
func DelPromiseBreak(pb *entity.PromiseBreak) error {
	pb.Status = entity.STATUS_DELETE
	pb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_PROMISE_BREAK).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(pb.Status, pb.UpdateAt, pb.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetPromiseBreakById
func GetPromiseBreakById(pdId int64) (*entity.PromiseBreak, error) {
	var pb entity.PromiseBreak
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,promise_id,happen_at,content_text").
		Form(TABLE_PROMISE_BREAK).
		Where("id=?").
		Query(pdId).
		NextScan(&pb.Id, &pb.Status, &pb.CreateAt, &pb.UpdateAt, &pb.UserId, &pb.CoupleId, &pb.PromiseId, &pb.HappenAt, &pb.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pb.Id <= 0 {
		return nil, nil
	} else if pb.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &pb, nil
}

// GetPromiseBreakListByCouplePromise
func GetPromiseBreakListByCouplePromise(cid, pid int64, offset, limit int) ([]*entity.PromiseBreak, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text").
		Form(TABLE_PROMISE_BREAK).
		Where("status>=? AND couple_id=? AND promise_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, pid)
	defer db.Close()
	list := make([]*entity.PromiseBreak, 0)
	for db.Next() {
		var pb entity.PromiseBreak
		pb.CoupleId = cid
		pb.PromiseId = pid
		db.Scan(&pb.Id, &pb.CreateAt, &pb.UpdateAt, &pb.UserId, &pb.HappenAt, &pb.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pb)
	}
	return list, nil
}

// GetPromiseBreakTotalByCouplePromise
func GetPromiseBreakTotalByCouplePromise(cid, pid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_PROMISE_BREAK).
		Where("status>=? AND couple_id=? AND promise_id=?").
		Query(entity.STATUS_VISIBLE, cid, pid).
		NextScan(&total)
	defer db.Close()
	return total
}
