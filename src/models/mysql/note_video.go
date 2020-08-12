package mysql

import (
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddVideo
func AddVideo(v *entity.Video) (*entity.Video, error) {
	v.Status = entity.STATUS_VISIBLE
	v.CreateAt = time.Now().Unix()
	v.UpdateAt = time.Now().Unix()
	v.ContentThumb = strings.TrimSpace(v.ContentThumb)
	v.ContentVideo = strings.TrimSpace(v.ContentVideo)
	db := mysqlDB().
		Insert(TABLE_VIDEO).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?,content_thumb=?,content_video=?,duration=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(v.Status, v.CreateAt, v.UpdateAt, v.UserId, v.CoupleId, v.HappenAt, v.Title, v.ContentThumb, v.ContentVideo, v.Duration, v.Longitude, v.Latitude, v.Address, v.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	v.Id, _ = db.Result().LastInsertId()
	return v, nil
}

// DelVideo
func DelVideo(v *entity.Video) error {
	v.Status = entity.STATUS_DELETE
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_VIDEO).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(v.Status, v.UpdateAt, v.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateVideo
func UpdateVideo(v *entity.Video) (*entity.Video, error) {
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_VIDEO).
		Set("update_at=?,happen_at=?,title=?,longitude=?,latitude=?,address=?,city_id=?").
		Where("id=?").
		Exec(v.UpdateAt, v.HappenAt, v.Title, v.Longitude, v.Latitude, v.Address, v.CityId, v.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return v, nil
}

// GetVideoById
func GetVideoById(vid int64) (*entity.Video, error) {
	var v entity.Video
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title,content_thumb,content_video,duration,longitude,latitude,address,city_id").
		Form(TABLE_VIDEO).
		Where("id=?").
		Query(vid).
		NextScan(&v.Id, &v.Status, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.CoupleId, &v.HappenAt, &v.Title, &v.ContentThumb, &v.ContentVideo, &v.Duration, &v.Longitude, &v.Latitude, &v.Address, &v.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if v.Id <= 0 {
		return nil, nil
	} else if v.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &v, nil
}

// GetVideoListByCouple
func GetVideoListByCouple(cid int64, offset, limit int) ([]*entity.Video, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,content_thumb,content_video,duration,longitude,latitude,address,city_id").
		Form(TABLE_VIDEO).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Video, 0)
	for db.Next() {
		var v entity.Video
		v.CoupleId = cid
		db.Scan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.HappenAt, &v.Title, &v.ContentThumb, &v.ContentVideo, &v.Duration, &v.Longitude, &v.Latitude, &v.Address, &v.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &v)
	}
	return list, nil
}

// GetVideoTotalByCouple
func GetVideoTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_VIDEO).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
