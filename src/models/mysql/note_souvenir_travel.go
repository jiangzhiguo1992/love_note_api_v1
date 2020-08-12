package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirTravel
func AddSouvenirTravel(st *entity.SouvenirTravel) (*entity.SouvenirTravel, error) {
	st.Status = entity.STATUS_VISIBLE
	st.CreateAt = time.Now().Unix()
	st.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_TRAVEL).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,travel_id=?,year=?").
		Exec(st.Status, st.CreateAt, st.UpdateAt, st.UserId, st.CoupleId, st.SouvenirId, st.TravelId, st.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	st.Id, _ = db.Result().LastInsertId()
	return st, nil
}

// DelSouvenirTravel
func DelSouvenirTravel(st *entity.SouvenirTravel) error {
	st.Status = entity.STATUS_DELETE
	st.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_TRAVEL).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(st.Status, st.UpdateAt, st.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirTravelById
func GetSouvenirTravelById(stid int64) (*entity.SouvenirTravel, error) {
	var st entity.SouvenirTravel
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,travel_id,year").
		Form(TABLE_SOUVENIR_TRAVEL).
		Where("id=?").
		Query(stid).
		NextScan(&st.Id, &st.Status, &st.CreateAt, &st.UpdateAt, &st.UserId, &st.CoupleId, &st.SouvenirId, &st.TravelId, &st.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if st.Id <= 0 {
		return nil, nil
	} else if st.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &st, nil
}

// GetSouvenirTravelByCoupleSouvenirTravel
func GetSouvenirTravelByCoupleSouvenirTravel(cid, sid, tid int64) (*entity.SouvenirTravel, error) {
	var st entity.SouvenirTravel
	st.CoupleId = cid
	st.SouvenirId = sid
	st.TravelId = tid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_TRAVEL).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND travel_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, tid).
		NextScan(&st.Id, &st.CreateAt, &st.UpdateAt, &st.UserId, &st.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if st.Id <= 0 {
		return nil, nil
	}
	return &st, nil
}

// GetSouvenirTravelListByCoupleSouvenir
func GetSouvenirTravelListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirTravel, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,travel_id,year").
		Form(TABLE_SOUVENIR_TRAVEL).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirTravel, 0)
	for db.Next() {
		var st entity.SouvenirTravel
		st.CoupleId = cid
		st.SouvenirId = sid
		db.Scan(&st.Id, &st.CreateAt, &st.UpdateAt, &st.UserId, &st.TravelId, &st.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &st)
	}
	return list, nil
}
