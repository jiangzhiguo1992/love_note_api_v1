package mysql

import (
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddAlbum
func AddAlbum(a *entity.Album) (*entity.Album, error) {
	a.Status = entity.STATUS_VISIBLE
	a.CreateAt = time.Now().Unix()
	a.UpdateAt = time.Now().Unix()
	a.Cover = strings.TrimSpace(a.Cover)
	a.StartAt = 0
	a.EndAt = 0
	a.PictureCount = 0
	db := mysqlDB().
		Insert(TABLE_ALBUM).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,title=?,cover=?,start_at=?,end_at=?,picture_count=?").
		Exec(a.Status, a.CreateAt, a.UpdateAt, a.UserId, a.CoupleId, a.Title, a.Cover, a.StartAt, a.EndAt, a.PictureCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	a.Id, _ = db.Result().LastInsertId()
	return a, nil
}

// DelAlbum
func DelAlbum(a *entity.Album) error {
	a.Status = entity.STATUS_DELETE
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_ALBUM).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(a.Status, a.UpdateAt, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateAlbum
func UpdateAlbum(a *entity.Album) (*entity.Album, error) {
	if a.PictureCount < 0 {
		a.PictureCount = 0
	}
	a.UpdateAt = time.Now().Unix()
	a.Cover = strings.TrimSpace(a.Cover)
	db := mysqlDB().
		Update(TABLE_ALBUM).
		Set("update_at=?,title=?,cover=?,start_at=?,end_at=?,picture_count=?").
		Where("id=?").
		Exec(a.UpdateAt, a.Title, a.Cover, a.StartAt, a.EndAt, a.PictureCount, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return a, nil
}

// GetAlbumById
func GetAlbumById(aid int64) (*entity.Album, error) {
	var a entity.Album
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,title,cover,start_at,end_at,picture_count").
		Form(TABLE_ALBUM).
		Where("id=?").
		Query(aid).
		NextScan(&a.Id, &a.Status, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.CoupleId, &a.Title, &a.Cover, &a.StartAt, &a.EndAt, &a.PictureCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if a.Id <= 0 {
		return nil, nil
	} else if a.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &a, nil
}

// GetAlbumListByCouple
func GetAlbumListByCouple(cid int64, offset, limit int) ([]*entity.Album, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,title,cover,start_at,end_at,picture_count").
		Form(TABLE_ALBUM).
		Where("status>=? AND couple_id=?").
		OrderDown("end_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Album, 0)
	for db.Next() {
		var a entity.Album
		a.CoupleId = cid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.Title, &a.Cover, &a.StartAt, &a.EndAt, &a.PictureCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAlbumTotalByCouple
func GetAlbumTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_ALBUM).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
