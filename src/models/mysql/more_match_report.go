package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMatchReport
func AddMatchReport(mr *entity.MatchReport) (*entity.MatchReport, error) {
	mr.Status = entity.STATUS_VISIBLE
	mr.CreateAt = time.Now().Unix()
	mr.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MATCH_REPORT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,match_period_id=?,match_work_id=?,reason=?").
		Exec(mr.Status, mr.CreateAt, mr.UpdateAt, mr.UserId, mr.CoupleId, mr.MatchPeriodId, mr.MatchWorkId, mr.Reason)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	mr.Id, _ = db.Result().LastInsertId()
	return mr, nil
}

// GetMatchReportByUserCoupleWork
func GetMatchReportByUserCoupleWork(uid, cid, mwid int64) (*entity.MatchReport, error) {
	var mr entity.MatchReport
	mr.UserId = uid
	mr.CoupleId = cid
	mr.MatchWorkId = mwid
	db := mysqlDB().
		Select("id,create_at,update_at,match_period_id,reason").
		Form(TABLE_MATCH_REPORT).
		Where("status>=? AND user_id=? AND couple_id=? AND match_work_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, mwid).
		NextScan(&mr.Id, &mr.CreateAt, &mr.UpdateAt, &mr.MatchPeriodId, &mr.Reason)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mr.Id <= 0 {
		return nil, nil
	}
	return &mr, nil
}

/****************************************** admin ***************************************/

// GetMatchReportList
func GetMatchReportList(offset, limit int) ([]*entity.MatchReport, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,match_period_id,match_work_id,reason").
		Form(TABLE_MATCH_REPORT).
		Where("status>=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.MatchReport, 0)
	for db.Next() {
		var mr entity.MatchReport
		db.Scan(&mr.Id, &mr.CreateAt, &mr.UpdateAt, &mr.UserId, &mr.CoupleId, &mr.MatchPeriodId, &mr.MatchWorkId, &mr.Reason)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &mr)
	}
	return list, nil
}
