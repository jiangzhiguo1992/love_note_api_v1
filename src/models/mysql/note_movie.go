package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMovie
func AddMovie(m *entity.Movie) (*entity.Movie, error) {
	m.Status = entity.STATUS_VISIBLE
	m.CreateAt = time.Now().Unix()
	m.UpdateAt = time.Now().Unix()
	m.ContentImages = entity.JoinStrByColon(m.ContentImageList)
	db := mysqlDB().
		Insert(TABLE_MOVIE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?,content_images=?,content_text=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(m.Status, m.CreateAt, m.UpdateAt, m.UserId, m.CoupleId, m.HappenAt, m.Title, m.ContentImages, m.ContentText, m.Longitude, m.Latitude, m.Address, m.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	m.Id, _ = db.Result().LastInsertId()
	m.ContentImages = ""
	return m, nil
}

// DelMovie
func DelMovie(m *entity.Movie) error {
	m.Status = entity.STATUS_DELETE
	m.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MOVIE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(m.Status, m.UpdateAt, m.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateMovie
func UpdateMovie(m *entity.Movie) (*entity.Movie, error) {
	m.UpdateAt = time.Now().Unix()
	m.ContentImages = entity.JoinStrByColon(m.ContentImageList)
	db := mysqlDB().
		Update(TABLE_MOVIE).
		Set("update_at=?,happen_at=?,title=?,content_images=?,content_text=?,longitude=?,latitude=?,address=?,city_id=?").
		Where("id=?").
		Exec(m.UpdateAt, m.HappenAt, m.Title, m.ContentImages, m.ContentText, m.Longitude, m.Latitude, m.Address, m.CityId, m.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	m.ContentImages = ""
	return m, nil
}

// GetMovieById
func GetMovieById(mid int64) (*entity.Movie, error) {
	var m entity.Movie
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title,content_images,content_text,longitude,latitude,address,city_id").
		Form(TABLE_MOVIE).
		Where("id=?").
		Query(mid).
		NextScan(&m.Id, &m.Status, &m.CreateAt, &m.UpdateAt, &m.UserId, &m.CoupleId, &m.HappenAt, &m.Title, &m.ContentImages, &m.ContentText, &m.Longitude, &m.Latitude, &m.Address, &m.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	} else if m.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	m.ContentImageList = entity.SplitStrByColon(m.ContentImages)
	m.ContentImages = ""
	return &m, nil
}

// GetMovieListByCouple
func GetMovieListByCouple(cid int64, offset, limit int) ([]*entity.Movie, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,content_images,content_text,longitude,latitude,address,city_id").
		Form(TABLE_MOVIE).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Movie, 0)
	for db.Next() {
		var m entity.Movie
		m.CoupleId = cid
		db.Scan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.UserId, &m.HappenAt, &m.Title, &m.ContentImages, &m.ContentText, &m.Longitude, &m.Latitude, &m.Address, &m.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		m.ContentImageList = entity.SplitStrByColon(m.ContentImages)
		m.ContentImages = ""
		list = append(list, &m)
	}
	return list, nil
}

// GetMovieTotalByCouple
func GetMovieTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_MOVIE).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
