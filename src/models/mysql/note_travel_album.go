package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelAlbum
func AddTravelAlbum(ta *entity.TravelAlbum) (*entity.TravelAlbum, error) {
	ta.Status = entity.STATUS_VISIBLE
	ta.CreateAt = time.Now().Unix()
	ta.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_ALBUM).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,album_id=?").
		Exec(ta.Status, ta.CreateAt, ta.UpdateAt, ta.UserId, ta.CoupleId, ta.TravelId, ta.AlbumId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	ta.Id, _ = db.Result().LastInsertId()
	return ta, nil
}

// DelTravelAlbum
func DelTravelAlbum(ta *entity.TravelAlbum) error {
	ta.Status = entity.STATUS_DELETE
	ta.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_ALBUM).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(ta.Status, ta.UpdateAt, ta.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelAlbumById
func GetTravelAlbumById(taid int64) (*entity.TravelAlbum, error) {
	var ta entity.TravelAlbum
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,album_id").
		Form(TABLE_TRAVEL_ALBUM).
		Where("id=?").
		Query(taid).
		NextScan(&ta.Id, &ta.Status, &ta.CreateAt, &ta.UpdateAt, &ta.UserId, &ta.CoupleId, &ta.TravelId, &ta.AlbumId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if ta.Id <= 0 {
		return nil, nil
	} else if ta.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &ta, nil
}

// GetTravelAlbumByCoupleAlbum
func GetTravelAlbumByCoupleTravelAlbum(cid, tid, aid int64) (*entity.TravelAlbum, error) {
	var ta entity.TravelAlbum
	ta.CoupleId = cid
	ta.TravelId = tid
	ta.AlbumId = aid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_TRAVEL_ALBUM).
		Where("status>=? AND couple_id=? AND travel_id=? AND album_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, tid, aid).
		NextScan(&ta.Id, &ta.CreateAt, &ta.UpdateAt, &ta.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if ta.Id <= 0 {
		return nil, nil
	}
	return &ta, nil
}

// GetTravelAlbumListByCoupleTravel
func GetTravelAlbumListByCoupleTravel(cid, tid int64) ([]*entity.TravelAlbum, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,album_id").
		Form(TABLE_TRAVEL_ALBUM).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelAlbum, 0)
	for db.Next() {
		var ta entity.TravelAlbum
		ta.CoupleId = cid
		ta.TravelId = tid
		db.Scan(&ta.Id, &ta.CreateAt, &ta.UpdateAt, &ta.UserId, &ta.AlbumId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &ta)
	}
	return list, nil
}

// GetTravelAlbumTotalByCoupleTravel
func GetTravelAlbumTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_ALBUM).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
