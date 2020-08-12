package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelPlace
func AddTravelPlace(tp *entity.TravelPlace) (*entity.TravelPlace, error) {
	tp.Status = entity.STATUS_VISIBLE
	tp.CreateAt = time.Now().Unix()
	tp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_PLACE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,happen_at=?,content_text=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(tp.Status, tp.CreateAt, tp.UpdateAt, tp.UserId, tp.CoupleId, tp.TravelId, tp.HappenAt, tp.ContentText, tp.Longitude, tp.Latitude, tp.Address, tp.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	tp.Id, _ = db.Result().LastInsertId()
	return tp, nil
}

// DelTravelPlace
func DelTravelPlace(tp *entity.TravelPlace) error {
	tp.Status = entity.STATUS_DELETE
	tp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_PLACE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(tp.Status, tp.UpdateAt, tp.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelPlaceById
func GetTravelPlaceById(tpid int64) (*entity.TravelPlace, error) {
	var tp entity.TravelPlace
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,happen_at,content_text,longitude,latitude,address,city_id").
		Form(TABLE_TRAVEL_PLACE).
		Where("id=?").
		Query(tpid).
		NextScan(&tp.Id, &tp.Status, &tp.CreateAt, &tp.UpdateAt, &tp.UserId, &tp.CoupleId, &tp.TravelId, &tp.HappenAt, &tp.ContentText, &tp.Longitude, &tp.Latitude, &tp.Address, &tp.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tp.Id <= 0 {
		return nil, nil
	} else if tp.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &tp, nil
}

// GetTravelPlaceListByCoupleTravel
func GetTravelPlaceListByCoupleTravel(cid, tid int64) ([]*entity.TravelPlace, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text,longitude,latitude,address,city_id").
		Form(TABLE_TRAVEL_PLACE).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("happen_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelPlace, 0)
	for db.Next() {
		var tp entity.TravelPlace
		tp.CoupleId = cid
		tp.TravelId = tid
		db.Scan(&tp.Id, &tp.CreateAt, &tp.UpdateAt, &tp.UserId, &tp.HappenAt, &tp.ContentText, &tp.Longitude, &tp.Latitude, &tp.Address, &tp.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &tp)
	}
	return list, nil
}

// GetTravelPlaceTotalByCoupleTravel
func GetTravelPlaceTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_PLACE).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
