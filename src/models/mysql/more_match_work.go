package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMatchWork
func AddMatchWork(mw *entity.MatchWork) (*entity.MatchWork, error) {
	mw.Status = entity.STATUS_VISIBLE
	mw.CreateAt = time.Now().Unix()
	mw.UpdateAt = time.Now().Unix()
	mw.ReportCount = 0
	mw.PointCount = 0
	mw.CoinCount = 0
	db := mysqlDB().
		Insert(TABLE_MATCH_WORK).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,match_period_id=?,kind=?,title=?,content_text=?,content_image=?,report_count=?,point_count=?,coin_count=?").
		Exec(mw.Status, mw.CreateAt, mw.UpdateAt, mw.UserId, mw.CoupleId, mw.MatchPeriodId, mw.Kind, mw.Title, mw.ContentText, mw.ContentImage, mw.ReportCount, mw.PointCount, mw.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	mw.Id, _ = db.Result().LastInsertId()
	return mw, nil
}

// DelMatchWork
func DelMatchWork(mw *entity.MatchWork) error {
	mw.Status = entity.STATUS_DELETE
	mw.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MATCH_WORK).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(mw.Status, mw.UpdateAt, mw.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateMatchWorkCount
func UpdateMatchWorkCount(mw *entity.MatchWork) (*entity.MatchWork, error) {
	if mw.ReportCount < 0 {
		mw.ReportCount = 0
	}
	if mw.PointCount < 0 {
		mw.PointCount = 0
	}
	if mw.CoinCount < 0 {
		mw.CoinCount = 0
	}
	mw.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MATCH_WORK).
		Set("update_at=?,report_count=?,point_count=?,coin_count=?").
		Where("id=?").
		Exec(mw.UpdateAt, mw.ReportCount, mw.PointCount, mw.CoinCount, mw.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return mw, nil
}

// GetMatchWorkById
func GetMatchWorkById(mwid int64) (*entity.MatchWork, error) {
	var mw entity.MatchWork
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,match_period_id,kind,title,content_text,content_image,report_count,point_count,coin_count").
		Form(TABLE_MATCH_WORK).
		Where("id=?").
		Query(mwid).
		NextScan(&mw.Id, &mw.Status, &mw.CreateAt, &mw.UpdateAt, &mw.UserId, &mw.CoupleId, &mw.MatchPeriodId, &mw.Kind, &mw.Title, &mw.ContentText, &mw.ContentImage, &mw.ReportCount, &mw.PointCount, &mw.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mw.Id <= 0 {
		return nil, nil
	} else if mw.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &mw, nil
}

// GetMatchWorkListByPeriodOrder
func GetMatchWorkListByPeriodOrder(mpid int64, order string, limitReportCount, offset, limit int) ([]*entity.MatchWork, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,title,content_text,content_image,report_count,point_count,coin_count").
		Form(TABLE_MATCH_WORK).
		Where("status>=? AND match_period_id=? AND report_count<?").
		Order(order).
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, mpid, limitReportCount)
	defer db.Close()
	list := make([]*entity.MatchWork, 0)
	for db.Next() {
		var mw entity.MatchWork
		mw.MatchPeriodId = mpid
		mw.Screen = false
		db.Scan(&mw.Id, &mw.CreateAt, &mw.UpdateAt, &mw.UserId, &mw.CoupleId, &mw.Kind, &mw.Title, &mw.ContentText, &mw.ContentImage, &mw.ReportCount, &mw.PointCount, &mw.CoinCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &mw)
	}
	return list, nil
}

// GetMatchWorkListByCoupleKind
func GetMatchWorkListByCoupleKind(cid int64, kind, offset, limit int) ([]*entity.MatchWork, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,match_period_id,title,content_text,content_image,report_count,point_count,coin_count").
		Form(TABLE_MATCH_WORK).
		Where("status>=? AND couple_id=? AND kind=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, kind)
	defer db.Close()
	list := make([]*entity.MatchWork, 0)
	for db.Next() {
		var mw entity.MatchWork
		mw.CoupleId = cid
		mw.Kind = kind
		mw.Screen = false
		db.Scan(&mw.Id, &mw.CreateAt, &mw.UpdateAt, &mw.UserId, &mw.MatchPeriodId, &mw.Title, &mw.ContentText, &mw.ContentImage, &mw.ReportCount, &mw.PointCount, &mw.CoinCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &mw)
	}
	return list, nil
}

// GetMatchWorkTotalByUserCouplePeriod
func GetMatchWorkTotalByUserCouplePeriod(uid, cid, mpid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_MATCH_WORK).
		Where("user_id=? AND couple_id=? AND match_period_id=?").
		Query(uid, cid, mpid).
		NextScan(&total)
	defer db.Close()
	return total
}

/****************************************** admin ***************************************/

// GetMatchWorkList
func GetMatchWorkList(uid, mpid int64, offset, limit int) ([]*entity.MatchWork, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasPeriod := mpid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasPeriod {
		where = where + " AND match_period_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,match_period_id,kind,title,content_text,content_image,report_count,point_count,coin_count").
		Form(TABLE_MATCH_WORK).
		Where(where).
		OrderDown("update_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasPeriod {
			db.Query(entity.STATUS_VISIBLE)
		} else {
			db.Query(entity.STATUS_VISIBLE, mpid)
		}
	} else {
		if !hasPeriod {
			db.Query(entity.STATUS_VISIBLE, uid)
		} else {
			db.Query(entity.STATUS_VISIBLE, uid, mpid)
		}
	}
	defer db.Close()
	list := make([]*entity.MatchWork, 0)
	for db.Next() {
		var mw entity.MatchWork
		db.Scan(&mw.Id, &mw.CreateAt, &mw.UpdateAt, &mw.UserId, &mw.CoupleId, &mw.MatchPeriodId, &mw.Kind, &mw.Title, &mw.ContentText, &mw.ContentImage, &mw.ReportCount, &mw.PointCount, &mw.CoinCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &mw)
	}
	return list, nil
}

// GetMatchWorkTotalByKindWithDel
func GetMatchWorkTotalByKindWithDel(create int64, kind int) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_MATCH_WORK).
		Where("create_at>? AND kind=?").
		Query(create, kind).
		NextScan(&total)
	defer db.Close()
	return total
}
