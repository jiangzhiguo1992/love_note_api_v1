package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMatchPoint
func AddMatchPoint(mp *entity.MatchPoint) (*entity.MatchPoint, error) {
	mp.Status = entity.STATUS_VISIBLE
	mp.CreateAt = time.Now().Unix()
	mp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MATCH_POINT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,match_period_id=?,match_work_id=?").
		Exec(mp.Status, mp.CreateAt, mp.UpdateAt, mp.UserId, mp.CoupleId, mp.MatchPeriodId, mp.MatchWorkId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	mp.Id, _ = db.Result().LastInsertId()
	return mp, nil
}

// UpdateMatchPoint
func UpdateMatchPoint(mp *entity.MatchPoint) (*entity.MatchPoint, error) {
	mp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MATCH_POINT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(mp.Status, mp.UpdateAt, mp.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return mp, nil
}

// GetMatchPointByUserCoupleWork
func GetMatchPointByUserCoupleWork(uid, cid, mwid int64) (*entity.MatchPoint, error) {
	var mp entity.MatchPoint
	mp.UserId = uid
	mp.CoupleId = cid
	mp.MatchWorkId = mwid
	db := mysqlDB().
		Select("id,status,create_at,update_at,match_period_id").
		Form(TABLE_MATCH_POINT).
		Where("user_id=? AND couple_id=? AND match_work_id=?").
		Query(uid, cid, mwid).
		NextScan(&mp.Id, &mp.Status, &mp.CreateAt, &mp.UpdateAt, &mp.MatchPeriodId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mp.Id <= 0 {
		return nil, nil
	}
	return &mp, nil
}
