package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTrendsBrowse
func AddTrendsBrowse(tb *entity.TrendsBrowse) (*entity.TrendsBrowse, error) {
	tb.Status = entity.STATUS_VISIBLE
	tb.CreateAt = time.Now().Unix()
	tb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRENDS_BROWSE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?").
		Exec(tb.Status, tb.CreateAt, tb.UpdateAt, tb.UserId, tb.CoupleId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	tb.Id, _ = db.Result().LastInsertId()
	return tb, nil
}

// UpdateTrendsBrowse
func UpdateTrendsBrowse(tb *entity.TrendsBrowse) (*entity.TrendsBrowse, error) {
	tb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRENDS_BROWSE).
		Set("update_at=?").
		Where("id=?").
		Exec(tb.UpdateAt, tb.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return tb, nil
}

// GetTrendsBrowseByUserCouple
func GetTrendsBrowseByUserCouple(uid, cid int64) (*entity.TrendsBrowse, error) {
	var tb entity.TrendsBrowse
	tb.UserId = uid
	tb.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_TRENDS_BROWSE).
		Where("status>=? AND user_id=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&tb.Id, &tb.CreateAt, &tb.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tb.Id <= 0 {
		return nil, nil
	}
	return &tb, nil
}
