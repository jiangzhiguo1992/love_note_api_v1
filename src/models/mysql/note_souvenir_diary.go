package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirDiary
func AddSouvenirDiary(sd *entity.SouvenirDiary) (*entity.SouvenirDiary, error) {
	sd.Status = entity.STATUS_VISIBLE
	sd.CreateAt = time.Now().Unix()
	sd.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_DIARY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,diary_id=?,year=?").
		Exec(sd.Status, sd.CreateAt, sd.UpdateAt, sd.UserId, sd.CoupleId, sd.SouvenirId, sd.DiaryId, sd.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sd.Id, _ = db.Result().LastInsertId()
	return sd, nil
}

// DelSouvenirDiary
func DelSouvenirDiary(sd *entity.SouvenirDiary) error {
	sd.Status = entity.STATUS_DELETE
	sd.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_DIARY).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sd.Status, sd.UpdateAt, sd.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirDiaryById
func GetSouvenirDiaryById(sdid int64) (*entity.SouvenirDiary, error) {
	var sd entity.SouvenirDiary
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,diary_id,year").
		Form(TABLE_SOUVENIR_DIARY).
		Where("id=?").
		Query(sdid).
		NextScan(&sd.Id, &sd.Status, &sd.CreateAt, &sd.UpdateAt, &sd.UserId, &sd.CoupleId, &sd.SouvenirId, &sd.DiaryId, &sd.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sd.Id <= 0 {
		return nil, nil
	} else if sd.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sd, nil
}

// GetSouvenirDiaryByCoupleSouvenirDiary
func GetSouvenirDiaryByCoupleSouvenirDiary(cid, sid, did int64) (*entity.SouvenirDiary, error) {
	var sd entity.SouvenirDiary
	sd.CoupleId = cid
	sd.SouvenirId = sid
	sd.DiaryId = did
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_DIARY).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND diary_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, did).
		NextScan(&sd.Id, &sd.CreateAt, &sd.UpdateAt, &sd.UserId, &sd.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sd.Id <= 0 {
		return nil, nil
	}
	return &sd, nil
}

// GetSouvenirDiaryListByCoupleSouvenir
func GetSouvenirDiaryListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirDiary, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,diary_id,year").
		Form(TABLE_SOUVENIR_DIARY).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirDiary, 0)
	for db.Next() {
		var sd entity.SouvenirDiary
		sd.CoupleId = cid
		sd.SouvenirId = sid
		db.Scan(&sd.Id, &sd.CreateAt, &sd.UpdateAt, &sd.UserId, &sd.DiaryId, &sd.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sd)
	}
	return list, nil
}
