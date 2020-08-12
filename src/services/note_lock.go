package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddLock
func AddLock(uid, cid int64, l *entity.Lock) (*entity.Lock, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if l == nil {
		return nil, errors.New("nil_lock")
	} else if len([]rune(l.Password)) <= 0 {
		return nil, errors.New("user_pwd_nil")
	}
	// 重复性检查
	old, _ := mysql.GetLockByUserCouple(uid, cid)
	if old != nil {
		return old, errors.New("lock_repeat")
	}
	// mysql
	l.UserId = uid
	l.CoupleId = cid
	l, err := mysql.AddLock(l)
	if l == nil || err != nil {
		return nil, err
	}
	return l, err
}

// UpdateLock
func UpdateLock(l *entity.Lock) (*entity.Lock, error) {
	if l == nil {
		return nil, errors.New("nil_lock")
	} else if len([]rune(l.Password)) <= 0 {
		return nil, errors.New("user_pwd_nil")
	}
	// mysql
	l, err := mysql.UpdateLock(l)
	return l, err
}

// GetLockByUserCouple
func GetLockByUserCouple(uid, cid int64) (*entity.Lock, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	l, err := mysql.GetLockByUserCouple(uid, cid)
	return l, err
}
