package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelMovie
func AddTravelMovie(tm *entity.TravelMovie) (*entity.TravelMovie, error) {
	tm.Status = entity.STATUS_VISIBLE
	tm.CreateAt = time.Now().Unix()
	tm.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_MOVIE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,movie_id=?").
		Exec(tm.Status, tm.CreateAt, tm.UpdateAt, tm.UserId, tm.CoupleId, tm.TravelId, tm.MovieId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	tm.Id, _ = db.Result().LastInsertId()
	return tm, nil
}

// DelTravelMovie
func DelTravelMovie(tm *entity.TravelMovie) error {
	tm.Status = entity.STATUS_DELETE
	tm.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_MOVIE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(tm.Status, tm.UpdateAt, tm.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelMovieById
func GetTravelMovieById(tmid int64) (*entity.TravelMovie, error) {
	var tm entity.TravelMovie
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,movie_id").
		Form(TABLE_TRAVEL_MOVIE).
		Where("id=?").
		Query(tmid).
		NextScan(&tm.Id, &tm.Status, &tm.CreateAt, &tm.UpdateAt, &tm.UserId, &tm.CoupleId, &tm.TravelId, &tm.MovieId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tm.Id <= 0 {
		return nil, nil
	} else if tm.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &tm, nil
}

// GetTravelMovieByCoupleTravelMovie
func GetTravelMovieByCoupleTravelMovie(cid, tid, mid int64) (*entity.TravelMovie, error) {
	var tm entity.TravelMovie
	tm.CoupleId = cid
	tm.TravelId = tid
	tm.MovieId = mid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_TRAVEL_MOVIE).
		Where("status>=? AND couple_id=? AND travel_id=? AND movie_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, tid, mid).
		NextScan(&tm.Id, &tm.CreateAt, &tm.UpdateAt, &tm.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tm.Id <= 0 {
		return nil, nil
	}
	return &tm, nil
}

// GetTravelMovieListByCoupleTravel
func GetTravelMovieListByCoupleTravel(cid, tid int64) ([]*entity.TravelMovie, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,movie_id").
		Form(TABLE_TRAVEL_MOVIE).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelMovie, 0)
	for db.Next() {
		var tm entity.TravelMovie
		tm.CoupleId = cid
		tm.TravelId = tid
		db.Scan(&tm.Id, &tm.CreateAt, &tm.UpdateAt, &tm.UserId, &tm.MovieId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &tm)
	}
	return list, nil
}

// GetTravelMovieTotalByCoupleTravel
func GetTravelMovieTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_MOVIE).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
