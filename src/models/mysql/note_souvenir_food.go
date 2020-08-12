package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirFood
func AddSouvenirFood(sf *entity.SouvenirFood) (*entity.SouvenirFood, error) {
	sf.Status = entity.STATUS_VISIBLE
	sf.CreateAt = time.Now().Unix()
	sf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_FOOD).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,food_id=?,year=?").
		Exec(sf.Status, sf.CreateAt, sf.UpdateAt, sf.UserId, sf.CoupleId, sf.SouvenirId, sf.FoodId, sf.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sf.Id, _ = db.Result().LastInsertId()
	return sf, nil
}

// DelSouvenirFood
func DelSouvenirFood(sf *entity.SouvenirFood) error {
	sf.Status = entity.STATUS_DELETE
	sf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_FOOD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sf.Status, sf.UpdateAt, sf.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirFoodById
func GetSouvenirFoodById(sfid int64) (*entity.SouvenirFood, error) {
	var sf entity.SouvenirFood
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,food_id,year").
		Form(TABLE_SOUVENIR_FOOD).
		Where("id=?").
		Query(sfid).
		NextScan(&sf.Id, &sf.Status, &sf.CreateAt, &sf.UpdateAt, &sf.UserId, &sf.CoupleId, &sf.SouvenirId, &sf.FoodId, &sf.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sf.Id <= 0 {
		return nil, nil
	} else if sf.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sf, nil
}

// GetSouvenirFoodByCoupleSouvenirFood
func GetSouvenirFoodByCoupleSouvenirFood(cid, sid, fid int64) (*entity.SouvenirFood, error) {
	var sf entity.SouvenirFood
	sf.CoupleId = cid
	sf.SouvenirId = sid
	sf.FoodId = fid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_FOOD).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND food_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, fid).
		NextScan(&sf.Id, &sf.CreateAt, &sf.UpdateAt, &sf.UserId, &sf.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sf.Id <= 0 {
		return nil, nil
	}
	return &sf, nil
}

// GetSouvenirFoodListByCoupleSouvenir
func GetSouvenirFoodListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirFood, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,food_id,year").
		Form(TABLE_SOUVENIR_FOOD).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirFood, 0)
	for db.Next() {
		var sf entity.SouvenirFood
		sf.CoupleId = cid
		sf.SouvenirId = sid
		db.Scan(&sf.Id, &sf.CreateAt, &sf.UpdateAt, &sf.UserId, &sf.FoodId, &sf.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sf)
	}
	return list, nil
}
