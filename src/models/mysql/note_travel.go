package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravel
func AddTravel(t *entity.Travel) (*entity.Travel, error) {
	t.Status = entity.STATUS_VISIBLE
	t.CreateAt = time.Now().Unix()
	t.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?").
		Exec(t.Status, t.CreateAt, t.UpdateAt, t.UserId, t.CoupleId, t.HappenAt, t.Title)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	t.Id, _ = db.Result().LastInsertId()
	return t, nil
}

// DelTravel
func DelTravel(t *entity.Travel) error {
	t.Status = entity.STATUS_DELETE
	t.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(t.Status, t.UpdateAt, t.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateTravel
func UpdateTravel(t *entity.Travel) (*entity.Travel, error) {
	t.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL).
		Set("update_at=?,happen_at=?,title=?").
		Where("id=?").
		Exec(t.UpdateAt, t.HappenAt, t.Title, t.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return t, nil
}

// GetTravelById
func GetTravelById(tid int64) (*entity.Travel, error) {
	var t entity.Travel
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title").
		Form(TABLE_TRAVEL).
		Where("id=?").
		Query(tid).
		NextScan(&t.Id, &t.Status, &t.CreateAt, &t.UpdateAt, &t.UserId, &t.CoupleId, &t.HappenAt, &t.Title)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if t.Id <= 0 {
		return nil, nil
	} else if t.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &t, nil
}

// GetTravelListByCouple
func GetTravelListByCouple(cid int64, offset, limit int) ([]*entity.Travel, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title").
		Form(TABLE_TRAVEL).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Travel, 0)
	for db.Next() {
		var t entity.Travel
		t.CoupleId = cid
		db.Scan(&t.Id, &t.CreateAt, &t.UpdateAt, &t.UserId, &t.HappenAt, &t.Title)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &t)
	}
	return list, nil
}

// GetTravelTotalByCouple
func GetTravelTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
