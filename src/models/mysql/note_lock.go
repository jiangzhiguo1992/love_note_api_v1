package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddLock
func AddLock(l *entity.Lock) (*entity.Lock, error) {
	l.IsLock = true
	l.Status = entity.STATUS_VISIBLE
	l.CreateAt = time.Now().Unix()
	l.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_LOCK).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,password=?,is_lock=?").
		Exec(l.Status, l.CreateAt, l.UpdateAt, l.UserId, l.CoupleId, l.Password, l.IsLock)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	l.Id, _ = db.Result().LastInsertId()
	return l, nil
}

// UpdateLock
func UpdateLock(l *entity.Lock) (*entity.Lock, error) {
	l.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_LOCK).
		Set("update_at=?,password=?,is_lock=?").
		Where("id=?").
		Exec(l.UpdateAt, l.Password, l.IsLock, l.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return l, nil
}

// GetLockById
func GetLockById(lid int64) (*entity.Lock, error) {
	var l entity.Lock
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,password,is_lock").
		Form(TABLE_LOCK).
		Where("id=?").
		Query(lid).
		NextScan(&l.Id, &l.Status, &l.CreateAt, &l.UpdateAt, &l.UserId, &l.CoupleId, &l.Password, &l.IsLock)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if l.Id <= 0 {
		return nil, nil
	} else if l.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &l, nil
}

// GetLockByUserCouple
func GetLockByUserCouple(uid, cid int64) (*entity.Lock, error) {
	var l entity.Lock
	l.UserId = uid
	l.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,password,is_lock").
		Form(TABLE_LOCK).
		Where("status>=? AND user_id=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&l.Id, &l.CreateAt, &l.UpdateAt, &l.Password, &l.IsLock)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if l.Id <= 0 {
		return nil, nil
	}
	return &l, nil
}
