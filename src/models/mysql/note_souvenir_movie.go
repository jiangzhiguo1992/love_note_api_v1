package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirMovie
func AddSouvenirMovie(sm *entity.SouvenirMovie) (*entity.SouvenirMovie, error) {
	sm.Status = entity.STATUS_VISIBLE
	sm.CreateAt = time.Now().Unix()
	sm.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_MOVIE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,movie_id=?,year=?").
		Exec(sm.Status, sm.CreateAt, sm.UpdateAt, sm.UserId, sm.CoupleId, sm.SouvenirId, sm.MovieId, sm.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sm.Id, _ = db.Result().LastInsertId()
	return sm, nil
}

// DelSouvenirMovie
func DelSouvenirMovie(sm *entity.SouvenirMovie) error {
	sm.Status = entity.STATUS_DELETE
	sm.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_MOVIE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sm.Status, sm.UpdateAt, sm.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirMovieById
func GetSouvenirMovieById(smid int64) (*entity.SouvenirMovie, error) {
	var sm entity.SouvenirMovie
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,movie_id,year").
		Form(TABLE_SOUVENIR_MOVIE).
		Where("id=?").
		Query(smid).
		NextScan(&sm.Id, &sm.Status, &sm.CreateAt, &sm.UpdateAt, &sm.UserId, &sm.CoupleId, &sm.SouvenirId, &sm.MovieId, &sm.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sm.Id <= 0 {
		return nil, nil
	} else if sm.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sm, nil
}

// GetSouvenirMovieByCoupleSouvenirMovie
func GetSouvenirMovieByCoupleSouvenirMovie(cid, sid, mid int64) (*entity.SouvenirMovie, error) {
	var sm entity.SouvenirMovie
	sm.CoupleId = cid
	sm.SouvenirId = sid
	sm.MovieId = mid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_MOVIE).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND movie_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, mid).
		NextScan(&sm.Id, &sm.CreateAt, &sm.UpdateAt, &sm.UserId, &sm.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sm.Id <= 0 {
		return nil, nil
	}
	return &sm, nil
}

// GetSouvenirMovieListByCoupleSouvenir
func GetSouvenirMovieListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirMovie, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,movie_id,year").
		Form(TABLE_SOUVENIR_MOVIE).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirMovie, 0)
	for db.Next() {
		var sm entity.SouvenirMovie
		sm.CoupleId = cid
		sm.SouvenirId = sid
		db.Scan(&sm.Id, &sm.CreateAt, &sm.UpdateAt, &sm.UserId, &sm.MovieId, &sm.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sm)
	}
	return list, nil
}
