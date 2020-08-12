package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddWallPaper
func AddWallPaper(wp *entity.WallPaper) (*entity.WallPaper, error) {
	wp.Status = entity.STATUS_VISIBLE
	wp.CreateAt = time.Now().Unix()
	wp.UpdateAt = time.Now().Unix()
	wp.ContentImages = entity.JoinStrByColon(wp.ContentImageList)
	db := mysqlDB().
		Insert(TABLE_WALL_PAPER).
		Set("status=?,create_at=?,update_at=?,couple_id=?,content_images=?").
		Exec(wp.Status, wp.CreateAt, wp.UpdateAt, wp.CoupleId, wp.ContentImages)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	wp.Id, _ = db.Result().LastInsertId()
	wp.ContentImages = ""
	return wp, nil
}

// UpdateWallPaper
func UpdateWallPaper(wp *entity.WallPaper) (*entity.WallPaper, error) {
	wp.UpdateAt = time.Now().Unix()
	wp.ContentImages = entity.JoinStrByColon(wp.ContentImageList)
	db := mysqlDB().
		Update(TABLE_WALL_PAPER).
		Set("update_at=?,content_images=?").
		Where("id=?").
		Exec(wp.UpdateAt, wp.ContentImages, wp.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	wp.ContentImages = ""
	return wp, nil
}

// GetWallPaperByCouple
func GetWallPaperByCouple(cid int64) (*entity.WallPaper, error) {
	var wp entity.WallPaper
	wp.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,content_images").
		Form(TABLE_WALL_PAPER).
		Where("status>=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&wp.Id, &wp.CreateAt, &wp.UpdateAt, &wp.ContentImages)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if wp.Id <= 0 {
		return nil, nil
	}
	wp.ContentImageList = entity.SplitStrByColon(wp.ContentImages)
	wp.ContentImages = ""
	return &wp, nil
}
