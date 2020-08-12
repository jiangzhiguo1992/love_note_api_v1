package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirAlbum
func AddSouvenirAlbum(sa *entity.SouvenirAlbum) (*entity.SouvenirAlbum, error) {
	sa.Status = entity.STATUS_VISIBLE
	sa.CreateAt = time.Now().Unix()
	sa.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_ALBUM).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,album_id=?,year=?").
		Exec(sa.Status, sa.CreateAt, sa.UpdateAt, sa.UserId, sa.CoupleId, sa.SouvenirId, sa.AlbumId, sa.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sa.Id, _ = db.Result().LastInsertId()
	return sa, nil
}

// DelSouvenirAlbum
func DelSouvenirAlbum(sa *entity.SouvenirAlbum) error {
	sa.Status = entity.STATUS_DELETE
	sa.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_ALBUM).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sa.Status, sa.UpdateAt, sa.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirAlbumById
func GetSouvenirAlbumById(said int64) (*entity.SouvenirAlbum, error) {
	var sa entity.SouvenirAlbum
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,album_id,year").
		Form(TABLE_SOUVENIR_ALBUM).
		Where("id=?").
		Query(said).
		NextScan(&sa.Id, &sa.Status, &sa.CreateAt, &sa.UpdateAt, &sa.UserId, &sa.CoupleId, &sa.SouvenirId, &sa.AlbumId, &sa.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sa.Id <= 0 {
		return nil, nil
	} else if sa.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sa, nil
}

// GetSouvenirAlbumByCoupleSouvenirAlbum
func GetSouvenirAlbumByCoupleSouvenirAlbum(cid, sid, aid int64) (*entity.SouvenirAlbum, error) {
	var sa entity.SouvenirAlbum
	sa.CoupleId = cid
	sa.SouvenirId = sid
	sa.AlbumId = aid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_ALBUM).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND album_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, aid).
		NextScan(&sa.Id, &sa.CreateAt, &sa.UpdateAt, &sa.UserId, &sa.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sa.Id <= 0 {
		return nil, nil
	}
	return &sa, nil
}

// GetSouvenirAlbumListByCoupleSouvenir
func GetSouvenirAlbumListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirAlbum, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,album_id,year").
		Form(TABLE_SOUVENIR_ALBUM).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirAlbum, 0)
	for db.Next() {
		var sa entity.SouvenirAlbum
		sa.CoupleId = cid
		sa.SouvenirId = sid
		db.Scan(&sa.Id, &sa.CreateAt, &sa.UpdateAt, &sa.UserId, &sa.AlbumId, &sa.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sa)
	}
	return list, nil
}
