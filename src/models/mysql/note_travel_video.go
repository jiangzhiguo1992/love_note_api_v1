package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTravelVideo
func AddTravelVideo(tv *entity.TravelVideo) (*entity.TravelVideo, error) {
	tv.Status = entity.STATUS_VISIBLE
	tv.CreateAt = time.Now().Unix()
	tv.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRAVEL_VIDEO).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,travel_id=?,video_id=?").
		Exec(tv.Status, tv.CreateAt, tv.UpdateAt, tv.UserId, tv.CoupleId, tv.TravelId, tv.VideoId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	tv.Id, _ = db.Result().LastInsertId()
	return tv, nil
}

// DelTravelVideo
func DelTravelVideo(tv *entity.TravelVideo) error {
	tv.Status = entity.STATUS_DELETE
	tv.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRAVEL_VIDEO).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(tv.Status, tv.UpdateAt, tv.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetTravelVideoById
func GetTravelVideoById(tvid int64) (*entity.TravelVideo, error) {
	var tv entity.TravelVideo
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,travel_id,video_id").
		Form(TABLE_TRAVEL_VIDEO).
		Where("id=?").
		Query(tvid).
		NextScan(&tv.Id, &tv.Status, &tv.CreateAt, &tv.UpdateAt, &tv.UserId, &tv.CoupleId, &tv.TravelId, &tv.VideoId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tv.Id <= 0 {
		return nil, nil
	} else if tv.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &tv, nil
}

// GetTravelVideoByCoupleTravelVideo
func GetTravelVideoByCoupleTravelVideo(cid, tid, vid int64) (*entity.TravelVideo, error) {
	var tv entity.TravelVideo
	tv.CoupleId = cid
	tv.TravelId = tid
	tv.VideoId = vid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_TRAVEL_VIDEO).
		Where("status>=? AND couple_id=? AND travel_id=? AND video_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, tid, vid).
		NextScan(&tv.Id, &tv.CreateAt, &tv.UpdateAt, &tv.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tv.Id <= 0 {
		return nil, nil
	}
	return &tv, nil
}

// GetTravelVideoListByCoupleTravel
func GetTravelVideoListByCoupleTravel(cid, tid int64) ([]*entity.TravelVideo, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,video_id").
		Form(TABLE_TRAVEL_VIDEO).
		Where("status>=? AND couple_id=? AND travel_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, tid)
	defer db.Close()
	list := make([]*entity.TravelVideo, 0)
	for db.Next() {
		var tv entity.TravelVideo
		tv.CoupleId = cid
		tv.TravelId = tid
		db.Scan(&tv.Id, &tv.CreateAt, &tv.UpdateAt, &tv.UserId, &tv.VideoId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &tv)
	}
	return list, nil
}

// GetTravelVideoTotalByCoupleTravel
func GetTravelVideoTotalByCoupleTravel(cid, tid int64) int {
	total := 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRAVEL_VIDEO).
		Where("status>=? AND couple_id=? AND travel_id=?").
		Query(entity.STATUS_VISIBLE, cid, tid).
		NextScan(&total)
	defer db.Close()
	return total
}
