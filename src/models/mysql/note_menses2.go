package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMenses2
func AddMenses2(m *entity.Menses2) (*entity.Menses2, error) {
	m.Status = entity.STATUS_VISIBLE
	m.CreateAt = time.Now().Unix()
	m.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MENSES2).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,start_at=?,end_at=?,start_year=?,start_month_of_year=?,start_day_of_month=?,end_year=?,end_month_of_year=?,end_day_of_month=?").
		Exec(m.Status, m.CreateAt, m.UpdateAt, m.UserId, m.CoupleId, m.StartAt, m.EndAt, m.StartYear, m.StartMonthOfYear, m.StartDayOfMonth, m.EndYear, m.EndMonthOfYear, m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	m.Id, _ = db.Result().LastInsertId()
	return m, nil
}

// DelMenses2
func DelMenses2(m *entity.Menses2) error {
	m.Status = entity.STATUS_DELETE
	m.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES2).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(m.Status, m.UpdateAt, m.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateMenses2
func UpdateMenses2(m *entity.Menses2) (*entity.Menses2, error) {
	m.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES2).
		Set("status=?,update_at=?,start_at=?,end_at=?,start_year=?,start_month_of_year=?,start_day_of_month=?,end_year=?,end_month_of_year=?,end_day_of_month=?").
		Where("id=?").
		Exec(m.Status, m.UpdateAt, m.StartAt, m.EndAt, m.StartYear, m.StartMonthOfYear, m.StartDayOfMonth, m.EndYear, m.EndMonthOfYear, m.EndDayOfMonth, m.Id)
	defer db.Close()
	if db.Err() != nil {
		return m, errors.New("db_update_fail")
	}
	return m, nil
}

// GetMenses2AllByUserCoupleDateStart
func GetMenses2AllByUserCoupleDateStart(uid, cid int64, year, month, day int) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	m.StartYear = year
	m.StartMonthOfYear = month
	m.StartDayOfMonth = day
	db := mysqlDB().
		Select("id,status,create_at,update_at,start_at,end_at,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("user_id=? AND couple_id=? AND start_year=? AND start_month_of_year=? AND start_day_of_month=?").
		Limit(0, 1).
		Query(uid, cid, year, month, day).
		NextScan(&m.Id, &m.Status, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	// 不检查status
	m.IsReal = true
	return &m, nil
}

// GetMenses2ByUserCoupleDateStart
func GetMenses2ByUserCoupleDateStart(uid, cid int64, year, month, day int) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	m.StartYear = year
	m.StartMonthOfYear = month
	m.StartDayOfMonth = day
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND start_year=? AND start_month_of_year=? AND start_day_of_month=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, year, month, day).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	m.IsReal = true
	return &m, nil
}

// GetMenses2ByByUserCoupleStartDown
func GetMenses2ByUserCoupleStartDown(uid, cid, startAt int64) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,start_year,start_month_of_year,start_day_of_month,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND start_at<=?").
		OrderDown("start_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, startAt).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.StartYear, &m.StartMonthOfYear, &m.StartDayOfMonth, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	m.IsReal = true
	return &m, nil
}

// GetMenses2ByUserCoupleEndNear
func GetMenses2ByUserCoupleEndNear(uid, cid, nearAt int64) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,start_year,start_month_of_year,start_day_of_month,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND end_at<=?").
		OrderDown("end_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, nearAt).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.StartYear, &m.StartMonthOfYear, &m.StartDayOfMonth, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	m.IsReal = true
	return &m, nil
}

// GetMenses2ByUserCoupleStartNear
func GetMenses2ByUserCoupleStartNear(uid, cid, nearAt int64) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,start_year,start_month_of_year,start_day_of_month,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND start_at>=?").
		OrderUp("start_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, nearAt).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.StartYear, &m.StartMonthOfYear, &m.StartDayOfMonth, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	m.IsReal = true
	return &m, nil
}

// GetMenses2ByByUserCoupleAtBetween
func GetMenses2ByUserCoupleAtBetween(uid, cid, betweenAt int64) (*entity.Menses2, error) {
	var m entity.Menses2
	m.UserId = uid
	m.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,start_year,start_month_of_year,start_day_of_month,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND start_at<=? AND end_at>=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, betweenAt, betweenAt).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.StartYear, &m.StartMonthOfYear, &m.StartDayOfMonth, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	m.IsReal = true
	return &m, nil
}

// GetMensesListByUserCoupleYearMonth
func GetMenses2ListByUserCoupleYearMonth(uid, cid int64, year, month, offset, limit int) ([]*entity.Menses2, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,start_at,end_at,start_year,start_month_of_year,start_day_of_month,end_year,end_month_of_year,end_day_of_month").
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND ((start_year=? AND start_month_of_year=?) OR (end_year=? AND end_month_of_year=?))").
		OrderUp("start_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid, year, month, year, month)
	defer db.Close()
	list := make([]*entity.Menses2, 0)
	for db.Next() {
		var m entity.Menses2
		m.UserId = uid
		m.CoupleId = cid
		db.Scan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.StartAt, &m.EndAt, &m.StartYear, &m.StartMonthOfYear, &m.StartDayOfMonth, &m.EndYear, &m.EndMonthOfYear, &m.EndDayOfMonth)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		m.IsReal = true
		list = append(list, &m)
	}
	return list, nil
}

// GetMenses2TotalByUserCoupleDateStart
func GetMenses2TotalByUserCoupleDateStart(uid, cid int64, year, month int) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_MENSES2).
		Where("status>=? AND user_id=? AND couple_id=? AND start_year=? AND start_month_of_year=?").
		Query(entity.STATUS_VISIBLE, uid, cid, year, month).
		NextScan(&total)
	defer db.Close()
	return total
}
