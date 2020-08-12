package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddDream
func AddDream(d *entity.Dream) (*entity.Dream, error) {
	d.Status = entity.STATUS_VISIBLE
	d.CreateAt = time.Now().Unix()
	d.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_DREAM).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,content_text=?").
		Exec(d.Status, d.CreateAt, d.UpdateAt, d.UserId, d.CoupleId, d.HappenAt, d.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	d.Id, _ = db.Result().LastInsertId()
	return d, nil
}

// DelDream
func DelDream(d *entity.Dream) error {
	d.Status = entity.STATUS_DELETE
	d.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_DREAM).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(d.Status, d.UpdateAt, d.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateDream
func UpdateDream(d *entity.Dream) (*entity.Dream, error) {
	d.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_DREAM).
		Set("update_at=?,happen_at=?,content_text=?").
		Where("id=?").
		Exec(d.UpdateAt, d.HappenAt, d.ContentText, d.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return d, nil
}

// GetDreamById
func GetDreamById(did int64) (*entity.Dream, error) {
	var d entity.Dream
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,content_text").
		Form(TABLE_DREAM).
		Where("id=?").
		Query(did).
		NextScan(&d.Id, &d.Status, &d.CreateAt, &d.UpdateAt, &d.UserId, &d.CoupleId, &d.HappenAt, &d.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if d.Id <= 0 {
		return nil, nil
	} else if d.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &d, nil
}

// GetDreamListByCouple
func GetDreamListByCouple(cid int64, offset, limit int) ([]*entity.Dream, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text").
		Form(TABLE_DREAM).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Dream, 0)
	for db.Next() {
		var d entity.Dream
		d.CoupleId = cid
		db.Scan(&d.Id, &d.CreateAt, &d.UpdateAt, &d.UserId, &d.HappenAt, &d.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &d)
	}
	return list, nil
}

// GetDreamListByUserCouple
func GetDreamListByUserCouple(uid, cid int64, offset, limit int) ([]*entity.Dream, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,happen_at,content_text").
		Form(TABLE_DREAM).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid)
	defer db.Close()
	list := make([]*entity.Dream, 0)
	for db.Next() {
		var d entity.Dream
		d.UserId = uid
		d.CoupleId = cid
		db.Scan(&d.Id, &d.CreateAt, &d.UpdateAt, &d.HappenAt, &d.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &d)
	}
	return list, nil
}

// GetDreamTotalByCouple
func GetDreamTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_DREAM).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
