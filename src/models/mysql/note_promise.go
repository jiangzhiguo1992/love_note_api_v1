package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPromise
func AddPromise(p *entity.Promise) (*entity.Promise, error) {
	p.Status = entity.STATUS_VISIBLE
	p.CreateAt = time.Now().Unix()
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_PROMISE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_id=?,happen_at=?,content_text=?,break_count=?").
		Exec(p.Status, p.CreateAt, p.UpdateAt, p.UserId, p.CoupleId, p.HappenId, p.HappenAt, p.ContentText, p.BreakCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	p.Id, _ = db.Result().LastInsertId()
	return p, nil
}

// DelPromise
func DelPromise(p *entity.Promise) error {
	p.Status = entity.STATUS_DELETE
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_PROMISE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(p.Status, p.UpdateAt, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdatePromise
func UpdatePromise(p *entity.Promise) (*entity.Promise, error) {
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_PROMISE).
		Set("update_at=?,happen_at=?,content_text=?,break_count=?").
		Where("id=?").
		Exec(p.UpdateAt, p.HappenAt, p.ContentText, p.BreakCount, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return p, nil
}

// GetPromiseById
func GetPromiseById(pid int64) (*entity.Promise, error) {
	var p entity.Promise
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_id,happen_at,content_text,break_count").
		Form(TABLE_PROMISE).
		Where("id=?").
		Query(pid).
		NextScan(&p.Id, &p.Status, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.HappenId, &p.HappenAt, &p.ContentText, &p.BreakCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	} else if p.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &p, nil
}

// GetPromiseListByCouple
func GetPromiseListByCouple(cid int64, offset, limit int) ([]*entity.Promise, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_id,happen_at,content_text,break_count").
		Form(TABLE_PROMISE).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Promise, 0)
	for db.Next() {
		var p entity.Promise
		p.CoupleId = cid
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.HappenId, &p.HappenAt, &p.ContentText, &p.BreakCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &p)
	}
	return list, nil
}

// GetPromiseListByCoupleHappenUser
func GetPromiseListByCoupleHappenUser(cid, hid int64, offset, limit int) ([]*entity.Promise, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text,break_count").
		Form(TABLE_PROMISE).
		Where("status>=? AND couple_id=? AND happen_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, hid)
	defer db.Close()
	list := make([]*entity.Promise, 0)
	for db.Next() {
		var p entity.Promise
		p.CoupleId = cid
		p.HappenId = hid
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.HappenAt, &p.ContentText, &p.BreakCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &p)
	}
	return list, nil
}

// GetPromiseTotalByCouple
func GetPromiseTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_PROMISE).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
