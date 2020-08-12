package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMatchPeriod
func AddMatchPeriod(mp *entity.MatchPeriod) (*entity.MatchPeriod, error) {
	mp.Status = entity.STATUS_VISIBLE
	mp.CreateAt = time.Now().Unix()
	mp.UpdateAt = time.Now().Unix()
	mp.WorksCount = 0
	mp.ReportCount = 0
	mp.PointCount = 0
	mp.CoinCount = 0
	db := mysqlDB().
		Insert(TABLE_MATCH_PERIOD).
		Set("status=?,create_at=?,update_at=?,start_at=?,end_at=?,period=?,kind=?,title=?,coin_change=?,works_count=?,report_count=?,point_count=?,coin_count=?").
		Exec(mp.Status, mp.CreateAt, mp.UpdateAt, mp.StartAt, mp.EndAt, mp.Period, mp.Kind, mp.Title, mp.CoinChange, mp.WorksCount, mp.ReportCount, mp.PointCount, mp.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	mp.Id, _ = db.Result().LastInsertId()
	return mp, nil
}

// DelMatchPeriod
func DelMatchPeriod(mw *entity.MatchPeriod) error {
	mw.Status = entity.STATUS_DELETE
	mw.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MATCH_PERIOD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(mw.Status, mw.UpdateAt, mw.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateMatchPeriodCount
func UpdateMatchPeriodCount(mp *entity.MatchPeriod) (*entity.MatchPeriod, error) {
	if mp.WorksCount < 0 {
		mp.WorksCount = 0
	}
	if mp.ReportCount < 0 {
		mp.ReportCount = 0
	}
	if mp.PointCount < 0 {
		mp.PointCount = 0
	}
	if mp.CoinCount < 0 {
		mp.CoinCount = 0
	}
	mp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MATCH_PERIOD).
		Set("update_at=?,works_count=?,report_count=?,point_count=?,coin_count=?").
		Where("id=?").
		Exec(mp.UpdateAt, mp.WorksCount, mp.ReportCount, mp.PointCount, mp.CoinCount, mp.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return mp, nil
}

// GetMatchPeriodById
func GetMatchPeriodById(mpid int64) (*entity.MatchPeriod, error) {
	var mp entity.MatchPeriod
	db := mysqlDB().
		Select("id,status,create_at,update_at,start_at,end_at,period,kind,title,coin_change,works_count,report_count,point_count,coin_count").
		Form(TABLE_MATCH_PERIOD).
		Where("id=?").
		Query(mpid).
		NextScan(&mp.Id, &mp.Status, &mp.CreateAt, &mp.UpdateAt, &mp.StartAt, &mp.EndAt, &mp.Period, &mp.Kind, &mp.Title, &mp.CoinChange, &mp.WorksCount, &mp.ReportCount, &mp.PointCount, &mp.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mp.Id <= 0 {
		return nil, nil
	} else if mp.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &mp, nil
}

// GetMatchPeriodNow
func GetMatchPeriodNow(kind int) (*entity.MatchPeriod, error) {
	now := time.Now().Unix()
	var mp entity.MatchPeriod
	mp.Kind = kind
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,period,title,coin_change,works_count,report_count,point_count,coin_count").
		Form(TABLE_MATCH_PERIOD).
		Where("status>=? AND kind=? AND start_at<? AND end_at>?").
		Query(entity.STATUS_VISIBLE, kind, now, now).
		NextScan(&mp.Id, &mp.CreateAt, &mp.UpdateAt, &mp.StartAt, &mp.EndAt, &mp.Period, &mp.Title, &mp.CoinChange, &mp.WorksCount, &mp.ReportCount, &mp.PointCount, &mp.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mp.Id <= 0 {
		return nil, nil
	}
	return &mp, nil
}

// GetMatchPeriodLatest
func GetMatchPeriodLatest(kind int) (*entity.MatchPeriod, error) {
	var mp entity.MatchPeriod
	mp.Kind = kind
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,period,title,coin_change,works_count,report_count,point_count,coin_count").
		Form(TABLE_MATCH_PERIOD).
		Where("status>=? AND kind=?").
		OrderDown("period").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, kind).
		NextScan(&mp.Id, &mp.CreateAt, &mp.UpdateAt, &mp.StartAt, &mp.EndAt, &mp.Period, &mp.Title, &mp.CoinChange, &mp.WorksCount, &mp.ReportCount, &mp.PointCount, &mp.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mp.Id <= 0 {
		return nil, nil
	}
	return &mp, nil
}

// GetMatchPeriodList
func GetMatchPeriodList(kind, offset, limit int) ([]*entity.MatchPeriod, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,period,title,coin_change,works_count,report_count,point_count,coin_count").
		Form(TABLE_MATCH_PERIOD).
		Where("status>=? AND kind=?").
		OrderDown("period").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, kind)
	defer db.Close()
	list := make([]*entity.MatchPeriod, 0)
	for db.Next() {
		var mp entity.MatchPeriod
		mp.Kind = kind
		db.Scan(&mp.Id, &mp.CreateAt, &mp.UpdateAt, &mp.StartAt, &mp.EndAt, &mp.Period, &mp.Title, &mp.CoinChange, &mp.WorksCount, &mp.ReportCount, &mp.PointCount, &mp.CoinCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &mp)
	}
	return list, nil
}
