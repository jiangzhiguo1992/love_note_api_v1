package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirVideo
func AddSouvenirVideo(sv *entity.SouvenirVideo) (*entity.SouvenirVideo, error) {
	sv.Status = entity.STATUS_VISIBLE
	sv.CreateAt = time.Now().Unix()
	sv.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_VIDEO).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,video_id=?,year=?").
		Exec(sv.Status, sv.CreateAt, sv.UpdateAt, sv.UserId, sv.CoupleId, sv.SouvenirId, sv.VideoId, sv.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sv.Id, _ = db.Result().LastInsertId()
	return sv, nil
}

// DelSouvenirVideo
func DelSouvenirVideo(sv *entity.SouvenirVideo) error {
	sv.Status = entity.STATUS_DELETE
	sv.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_VIDEO).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sv.Status, sv.UpdateAt, sv.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirVideoById
func GetSouvenirVideoById(svid int64) (*entity.SouvenirVideo, error) {
	var sv entity.SouvenirVideo
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,video_id,year").
		Form(TABLE_SOUVENIR_VIDEO).
		Where("id=?").
		Query(svid).
		NextScan(&sv.Id, &sv.Status, &sv.CreateAt, &sv.UpdateAt, &sv.UserId, &sv.CoupleId, &sv.SouvenirId, &sv.VideoId, &sv.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sv.Id <= 0 {
		return nil, nil
	} else if sv.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sv, nil
}

// GetSouvenirVideoByCoupleSouvenirVideo
func GetSouvenirVideoByCoupleSouvenirVideo(cid, sid, vid int64) (*entity.SouvenirVideo, error) {
	var sv entity.SouvenirVideo
	sv.CoupleId = cid
	sv.SouvenirId = sid
	sv.VideoId = vid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_VIDEO).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND video_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, vid).
		NextScan(&sv.Id, &sv.CreateAt, &sv.UpdateAt, &sv.UserId, &sv.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sv.Id <= 0 {
		return nil, nil
	}
	return &sv, nil
}

// GetSouvenirVideoListByCoupleSouvenir
func GetSouvenirVideoListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirVideo, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,video_id,year").
		Form(TABLE_SOUVENIR_VIDEO).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirVideo, 0)
	for db.Next() {
		var sv entity.SouvenirVideo
		sv.CoupleId = cid
		sv.SouvenirId = sid
		db.Scan(&sv.Id, &sv.CreateAt, &sv.UpdateAt, &sv.UserId, &sv.VideoId, &sv.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sv)
	}
	return list, nil
}
