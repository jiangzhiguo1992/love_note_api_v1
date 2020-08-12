package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenir
func AddSouvenir(s *entity.Souvenir) (*entity.Souvenir, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?,done=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, s.UserId, s.CoupleId, s.HappenAt, s.Title, s.Done, s.Longitude, s.Latitude, s.Address, s.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// DelSouvenir
func DelSouvenir(s *entity.Souvenir) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(s.Status, s.UpdateAt, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateSouvenir
func UpdateSouvenir(s *entity.Souvenir) (*entity.Souvenir, error) {
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR).
		Set("update_at=?,happen_at=?,title=?,done=?,longitude=?,latitude=?,address=?,city_id=?").
		Where("id=?").
		Exec(s.UpdateAt, s.HappenAt, s.Title, s.Done, s.Longitude, s.Latitude, s.Address, s.CityId, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return s, nil
}

// GetSouvenirById 查看单个Souvenir
func GetSouvenirById(sid int64) (*entity.Souvenir, error) {
	var s entity.Souvenir
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title,done,longitude,latitude,address,city_id").
		Form(TABLE_SOUVENIR).
		Where("id=?").
		Query(sid).
		NextScan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.CoupleId, &s.HappenAt, &s.Title, &s.Done, &s.Longitude, &s.Latitude, &s.Address, &s.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	} else if s.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &s, nil
}

// GetSouvenirDoneListByCouple
func GetSouvenirDoneListByCouple(cid int64, offset, limit int) ([]*entity.Souvenir, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,longitude,latitude,address,city_id").
		Form(TABLE_SOUVENIR).
		Where("status>=? AND couple_id=? AND done=?").
		OrderUp("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, true)
	defer db.Close()
	list := make([]*entity.Souvenir, 0)
	for db.Next() {
		var s entity.Souvenir
		s.CoupleId = cid
		s.Done = true
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.HappenAt, &s.Title, &s.Longitude, &s.Latitude, &s.Address, &s.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetSouvenirWishListByUserCouple
func GetSouvenirWishListByUserCouple(cid int64, offset, limit int) ([]*entity.Souvenir, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,longitude,latitude,address,city_id").
		Form(TABLE_SOUVENIR).
		Where("status>=? AND couple_id=? AND done=?").
		OrderUp("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, false)
	defer db.Close()
	list := make([]*entity.Souvenir, 0)
	for db.Next() {
		var s entity.Souvenir
		s.CoupleId = cid
		s.Done = false
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.HappenAt, &s.Title, &s.Longitude, &s.Latitude, &s.Address, &s.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetSouvenirTotalByCouple
func GetSouvenirTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SOUVENIR).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
