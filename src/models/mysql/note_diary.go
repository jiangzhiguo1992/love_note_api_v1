package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddDiary
func AddDiary(d *entity.Diary) (*entity.Diary, error) {
	d.Status = entity.STATUS_VISIBLE
	d.CreateAt = time.Now().Unix()
	d.UpdateAt = time.Now().Unix()
	d.ReadCount = 0
	d.ContentImages = entity.JoinStrByColon(d.ContentImageList)
	db := mysqlDB().
		Insert(TABLE_DIARY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,content_text=?,content_images=?,read_count=?").
		Exec(d.Status, d.CreateAt, d.UpdateAt, d.UserId, d.CoupleId, d.HappenAt, d.ContentText, d.ContentImages, d.ReadCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	d.Id, _ = db.Result().LastInsertId()
	d.ContentImages = ""
	return d, nil
}

// DelDiary
func DelDiary(d *entity.Diary) error {
	d.Status = entity.STATUS_DELETE
	d.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_DIARY).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(d.Status, d.UpdateAt, d.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateDiary
func UpdateDiary(d *entity.Diary) (*entity.Diary, error) {
	d.UpdateAt = time.Now().Unix()
	d.ContentImages = entity.JoinStrByColon(d.ContentImageList)
	db := mysqlDB().
		Update(TABLE_DIARY).
		Set("update_at=?,happen_at=?,content_text=?,content_images=?").
		Where("id=?").
		Exec(d.UpdateAt, d.HappenAt, d.ContentText, d.ContentImages, d.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	d.ContentImages = ""
	return d, nil
}

// UpdateDiaryReadCount
func UpdateDiaryReadCount(d *entity.Diary) (*entity.Diary, error) {
	db := mysqlDB().
		Update(TABLE_DIARY).
		Set("read_count=?").
		Where("id=?").
		Exec(d.ReadCount, d.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return d, nil
}

// GetDiaryById
func GetDiaryById(did int64) (*entity.Diary, error) {
	var d entity.Diary
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,content_text,content_images,read_count").
		Form(TABLE_DIARY).
		Where("id=?").
		Query(did).
		NextScan(&d.Id, &d.Status, &d.CreateAt, &d.UpdateAt, &d.UserId, &d.CoupleId, &d.HappenAt, &d.ContentText, &d.ContentImages, &d.ReadCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if d.Id <= 0 {
		return nil, nil
	} else if d.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	d.ContentImageList = entity.SplitStrByColon(d.ContentImages)
	d.ContentImages = ""
	return &d, nil
}

// GetDiaryListByCouple
func GetDiaryListByCouple(cid int64, offset, limit int) ([]*entity.Diary, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,content_text,content_images,read_count").
		Form(TABLE_DIARY).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Diary, 0)
	for db.Next() {
		var d entity.Diary
		d.CoupleId = cid
		db.Scan(&d.Id, &d.CreateAt, &d.UpdateAt, &d.UserId, &d.HappenAt, &d.ContentText, &d.ContentImages, &d.ReadCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		d.ContentImageList = entity.SplitStrByColon(d.ContentImages)
		d.ContentImages = ""
		list = append(list, &d)
	}
	return list, nil
}

// GetDiaryListByUserCouple
func GetDiaryListByUserCouple(uid, cid int64, offset, limit int) ([]*entity.Diary, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,happen_at,content_text,content_images,read_count").
		Form(TABLE_DIARY).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid)
	defer db.Close()
	list := make([]*entity.Diary, 0)
	for db.Next() {
		var d entity.Diary
		d.UserId = uid
		d.CoupleId = cid
		db.Scan(&d.Id, &d.CreateAt, &d.UpdateAt, &d.HappenAt, &d.ContentText, &d.ContentImages, &d.ReadCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		d.ContentImageList = entity.SplitStrByColon(d.ContentImages)
		d.ContentImages = ""
		list = append(list, &d)
	}
	return list, nil
}

// GetDiaryTotalByCouple
func GetDiaryTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_DIARY).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
