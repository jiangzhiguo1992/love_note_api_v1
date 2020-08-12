package mysql

import (
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddAudio
func AddAudio(a *entity.Audio) (*entity.Audio, error) {
	a.Status = entity.STATUS_VISIBLE
	a.CreateAt = time.Now().Unix()
	a.UpdateAt = time.Now().Unix()
	a.ContentAudio = strings.TrimSpace(a.ContentAudio)
	db := mysqlDB().
		Insert(TABLE_AUDIO).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?,content_audio=?,duration=?").
		Exec(a.Status, a.CreateAt, a.UpdateAt, a.UserId, a.CoupleId, a.HappenAt, a.Title, a.ContentAudio, a.Duration)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	a.Id, _ = db.Result().LastInsertId()
	return a, nil
}

// DelAudio
func DelAudio(a *entity.Audio) error {
	a.Status = entity.STATUS_DELETE
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_AUDIO).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(a.Status, a.UpdateAt, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetAudioById
func GetAudioById(aid int64) (*entity.Audio, error) {
	var a entity.Audio
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title,content_audio,duration").
		Form(TABLE_AUDIO).
		Where("id=?").
		Query(aid).
		NextScan(&a.Id, &a.Status, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.CoupleId, &a.HappenAt, &a.Title, &a.ContentAudio, &a.Duration)
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

// GetAudioListByCouple
func GetAudioListByCouple(cid int64, offset, limit int) ([]*entity.Audio, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,content_audio,duration").
		Form(TABLE_AUDIO).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Audio, 0)
	for db.Next() {
		var a entity.Audio
		a.CoupleId = cid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.HappenAt, &a.Title, &a.ContentAudio, &a.Duration)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAudioTotalByCouple
func GetAudioTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_AUDIO).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
