package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelDiary
func AddTravelDiary(td *entity.TravelDiary) (*entity.TravelDiary, error) {
	td.Status = entity.STATUS_VISIBLE
	td.CreateAt = time.Now().Unix()
	td.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_DIARY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,diary_id=?").
		Exec(td.Status, td.CreateAt, td.UpdateAt, td.UserId, td.CoupleId, td.TravelId, td.DiaryId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	td.Id, _ = db.Result().LastInsertId()
	return td, nil
}

// DelTravelDiary
func DelTravelDiary(td *entity.TravelDiary) error {
	td.Status = entity.STATUS_DELETE
	td.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_DIARY).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(td.Status, td.UpdateAt, td.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelDiaryById
func GetTravelDiaryById(tdid int64) (*entity.TravelDiary, error) {
	var td entity.TravelDiary
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,diary_id").
		Form(TABLE_TRAVEL_DIARY).
		Where("id=?").
		Query(tdid).
		NextScan(&td.Id, &td.Status, &td.CreateAt, &td.UpdateAt, &td.UserId, &td.CoupleId, &td.TravelId, &td.DiaryId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if td.Id <= 0 {
		return nil, nil
	} else if td.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &td, nil
}

// GetTravelDiaryByCoupleTravelDiary
func GetTravelDiaryByCoupleTravelDiary(cid, tid, did int64) (*entity.TravelDiary, error) {
	var td entity.TravelDiary
	td.CoupleId = cid
	td.TravelId = tid
	td.DiaryId = did
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_TRAVEL_DIARY).
		Where("status>=? AND couple_id=? AND travel_id=? AND diary_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, tid, did).
		NextScan(&td.Id, &td.CreateAt, &td.UpdateAt, &td.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if td.Id <= 0 {
		return nil, nil
	}
	return &td, nil
}

// GetTravelDiaryListByCoupleTravel
func GetTravelDiaryListByCoupleTravel(cid, tid int64) ([]*entity.TravelDiary, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,diary_id").
		Form(TABLE_TRAVEL_DIARY).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelDiary, 0)
	for db.Next() {
		var td entity.TravelDiary
		td.CoupleId = cid
		td.TravelId = tid
		db.Scan(&td.Id, &td.CreateAt, &td.UpdateAt, &td.UserId, &td.DiaryId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &td)
	}
	return list, nil
}

// GetTravelDiaryTotalByCoupleTravel
func GetTravelDiaryTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_DIARY).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
