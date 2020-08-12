package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddAngry
func AddAngry(a *entity.Angry) (*entity.Angry, error) {
	a.Status = entity.STATUS_VISIBLE
	a.CreateAt = time.Now().Unix()
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_ANGRY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_id=?,happen_at=?,content_text=?,gift_id=?,promise_id=?").
		Exec(a.Status, a.CreateAt, a.UpdateAt, a.UserId, a.CoupleId, a.HappenId, a.HappenAt, a.ContentText, a.GiftId, a.PromiseId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	a.Id, _ = db.Result().LastInsertId()
	return a, nil
}

// DelAngry
func DelAngry(a *entity.Angry) error {
	a.Status = entity.STATUS_DELETE
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_ANGRY).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(a.Status, a.UpdateAt, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateAngry
func UpdateAngry(a *entity.Angry) (*entity.Angry, error) {
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_ANGRY).
		Set("update_at=?,gift_id=?,promise_id=?").
		Where("id=?").
		Exec(a.UpdateAt, a.GiftId, a.PromiseId, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return a, nil
}

// GetAngryById
func GetAngryById(aid int64) (*entity.Angry, error) {
	var a entity.Angry
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_id,happen_at,content_text,gift_id,promise_id").
		Form(TABLE_ANGRY).
		Where("id=?").
		Query(aid).
		NextScan(&a.Id, &a.Status, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.CoupleId, &a.HappenId, &a.HappenAt, &a.ContentText, &a.GiftId, &a.PromiseId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if a.Id <= 0 {
		return nil, nil
	} else if a.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &a, nil
}

// GetAngryListByCouple
func GetAngryListByCouple(cid int64, offset, limit int) ([]*entity.Angry, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_id,happen_at,content_text,gift_id,promise_id").
		Form(TABLE_ANGRY).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Angry, 0)
	for db.Next() {
		var a entity.Angry
		a.CoupleId = cid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.HappenId, &a.HappenAt, &a.ContentText, &a.GiftId, &a.PromiseId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAngryListByCoupleHappenUser
func GetAngryListByCoupleHappenUser(cid, hid int64, offset, limit int) ([]*entity.Angry, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text,gift_id,promise_id").
		Form(TABLE_ANGRY).
		Where("status>=? AND couple_id=? AND happen_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, hid)
	defer db.Close()
	list := make([]*entity.Angry, 0)
	for db.Next() {
		var a entity.Angry
		a.CoupleId = cid
		a.HappenId = hid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.HappenAt, &a.ContentText, &a.GiftId, &a.PromiseId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAngryTotalByCouple
func GetAngryTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_ANGRY).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
