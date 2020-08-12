package mysql

import (
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddPicture
func AddPicture(p *entity.Picture) (*entity.Picture, error) {
	p.Status = entity.STATUS_VISIBLE
	p.CreateAt = time.Now().Unix()
	p.UpdateAt = time.Now().Unix()
	p.ContentImage = strings.TrimSpace(p.ContentImage)
	db := mysqlDB().
		Insert(TABLE_PICTURE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,album_id=?,happen_at=?,content_image=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(p.Status, p.CreateAt, p.UpdateAt, p.UserId, p.CoupleId, p.AlbumId, p.HappenAt, p.ContentImage, p.Longitude, p.Latitude, p.Address, p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	p.Id, _ = db.Result().LastInsertId()
	return p, nil
}

// DelPicture
func DelPicture(p *entity.Picture) error {
	p.Status = entity.STATUS_DELETE
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_PICTURE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(p.Status, p.UpdateAt, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdatePicture
func UpdatePicture(p *entity.Picture) (*entity.Picture, error) {
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_PICTURE).
		Set("update_at=?,album_id=?,happen_at=?,longitude=?,latitude=?,address=?,city_id=?").
		Where("id=?").
		Exec(p.UpdateAt, p.AlbumId, p.HappenAt, p.Longitude, p.Latitude, p.Address, p.CityId, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return p, nil
}

// GetPictureById
func GetPictureById(pid int64) (*entity.Picture, error) {
	var p entity.Picture
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,album_id,happen_at,content_image,longitude,latitude,address,city_id").
		Form(TABLE_PICTURE).
		Where("id=?").
		Query(pid).
		NextScan(&p.Id, &p.Status, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.AlbumId, &p.HappenAt, &p.ContentImage, &p.Longitude, &p.Latitude, &p.Address, &p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	} else if p.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &p, nil
}

// GetPictureStartByAlbum
func GetPictureStartByAlbum(aid int64) (*entity.Picture, error) {
	var p entity.Picture
	p.AlbumId = aid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,happen_at,content_image,longitude,latitude,address,city_id").
		Form(TABLE_PICTURE).
		Where("status>=? AND happen_at<>0 AND album_id=?").
		OrderUp("happen_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, aid).
		NextScan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.HappenAt, &p.ContentImage, &p.Longitude, &p.Latitude, &p.Address, &p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	}
	return &p, nil
}

// GetPictureEndByAlbum
func GetPictureEndByAlbum(aid int64) (*entity.Picture, error) {
	var p entity.Picture
	p.AlbumId = aid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,happen_at,content_image,longitude,latitude,address,city_id").
		Form(TABLE_PICTURE).
		Where("status>=? AND happen_at<>0 AND album_id=?").
		OrderDown("happen_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, aid).
		NextScan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.HappenAt, &p.ContentImage, &p.Longitude, &p.Latitude, &p.Address, &p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	}
	return &p, nil
}

// GetPictureListByCoupleAlbum
func GetPictureListByCoupleAlbum(cid, aid int64, offset, limit int) ([]*entity.Picture, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_image,longitude,latitude,address,city_id").
		Form(TABLE_PICTURE).
		Where("status>=? AND couple_id=? AND album_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, aid)
	defer db.Close()
	list := make([]*entity.Picture, 0)
	for db.Next() {
		var p entity.Picture
		p.CoupleId = cid
		p.AlbumId = aid
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.HappenAt, &p.ContentImage, &p.Longitude, &p.Latitude, &p.Address, &p.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &p)
	}
	return list, nil
}

// GetPictureTotalByAlbum
func GetPictureTotalByCoupleAlbum(cid, aid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_PICTURE).
		Where("status>=? AND couple_id=? AND album_id=?").
		Query(entity.STATUS_VISIBLE, cid, aid).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetPictureTotalByCouple
func GetPictureTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_PICTURE).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
