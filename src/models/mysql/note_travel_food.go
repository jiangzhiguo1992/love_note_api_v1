package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelFood
func AddTravelFood(tf *entity.TravelFood) (*entity.TravelFood, error) {
	tf.Status = entity.STATUS_VISIBLE
	tf.CreateAt = time.Now().Unix()
	tf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_FOOD).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,food_id=?").
		Exec(tf.Status, tf.CreateAt, tf.UpdateAt, tf.UserId, tf.CoupleId, tf.TravelId, tf.FoodId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	tf.Id, _ = db.Result().LastInsertId()
	return tf, nil
}

// DelTravelFood
func DelTravelFood(tf *entity.TravelFood) error {
	tf.Status = entity.STATUS_DELETE
	tf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_FOOD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(tf.Status, tf.UpdateAt, tf.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelFoodById
func GetTravelFoodById(tfid int64) (*entity.TravelFood, error) {
	var tf entity.TravelFood
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,food_id").
		Form(TABLE_TRAVEL_FOOD).
		Where("id=?").
		Query(tfid).
		NextScan(&tf.Id, &tf.Status, &tf.CreateAt, &tf.UpdateAt, &tf.UserId, &tf.CoupleId, &tf.TravelId, &tf.FoodId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tf.Id <= 0 {
		return nil, nil
	} else if tf.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &tf, nil
}

// GetTravelFoodByCoupleTravelFood
func GetTravelFoodByCoupleTravelFood(cid, tid, fid int64) (*entity.TravelFood, error) {
	var tf entity.TravelFood
	tf.CoupleId = cid
	tf.TravelId = tid
	tf.FoodId = fid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_TRAVEL_FOOD).
		Where("status>=? AND couple_id=? AND travel_id=? AND food_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, tid, fid).
		NextScan(&tf.Id, &tf.CreateAt, &tf.UpdateAt, &tf.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tf.Id <= 0 {
		return nil, nil
	}
	return &tf, nil
}

// GetTravelFoodListByCoupleTravel
func GetTravelFoodListByCoupleTravel(cid, tid int64) ([]*entity.TravelFood, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,food_id").
		Form(TABLE_TRAVEL_FOOD).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelFood, 0)
	for db.Next() {
		var tf entity.TravelFood
		tf.CoupleId = cid
		tf.TravelId = tid
		db.Scan(&tf.Id, &tf.CreateAt, &tf.UpdateAt, &tf.UserId, &tf.FoodId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &tf)
	}
	return list, nil
}

// GetTravelFoodTotalByCoupleTravel
func GetTravelFoodTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_FOOD).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
