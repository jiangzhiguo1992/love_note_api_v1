package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMenses
func AddMenses(m *entity.Menses) (*entity.Menses, error) {
	m.Status = entity.STATUS_VISIBLE
	m.CreateAt = time.Now().Unix()
	m.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MENSES).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,year=?,month_of_year=?,day_of_month=?,is_start=?").
		Exec(m.Status, m.CreateAt, m.UpdateAt, m.UserId, m.CoupleId, m.Year, m.MonthOfYear, m.DayOfMonth, m.IsStart)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	m.Id, _ = db.Result().LastInsertId()
	return m, nil
}

// GetMensesById
func GetMensesById(mid int64) (*entity.Menses, error) {
	var m entity.Menses
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,year,month_of_year,day_of_month,is_start").
		Form(TABLE_MENSES).
		Where("id=?").
		Query(mid).
		NextScan(&m.Id, &m.Status, &m.CreateAt, &m.UpdateAt, &m.UserId, &m.CoupleId, &m.Year, &m.MonthOfYear, &m.DayOfMonth, &m.IsStart)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	} else if m.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &m, nil
}

// DelMenses
func DelMenses(s *entity.Menses) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(s.Status, s.UpdateAt, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// DelMensesByUserCoupleDate
func DelMensesByUserCoupleDate(s *entity.Menses) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES).
		Set("status=?,update_at=?").
		Where("user_id=? AND couple_id=? AND year=? AND month_of_year=? AND day_of_month=?").
		Exec(s.Status, s.UpdateAt, s.UserId, s.CoupleId, s.Year, s.MonthOfYear, s.DayOfMonth)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetMensesLatestByUserCouple
func GetMensesLatestByUserCouple(uid, cid int64) (*entity.Menses, error) {
	var m entity.Menses
	m.UserId = uid
	m.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,year,month_of_year,day_of_month,is_start").
		Form(TABLE_MENSES).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.Year, &m.MonthOfYear, &m.DayOfMonth, &m.IsStart)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if m.Id <= 0 {
		return nil, nil
	}
	return &m, nil
}

// GetMensesListByUserCoupleYearMonth
func GetMensesListByUserCoupleYearMonth(uid, cid int64, year, month, offset, limit int) ([]*entity.Menses, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,day_of_month,is_start").
		Form(TABLE_MENSES).
		Where("status>=? AND user_id=? AND couple_id=? AND year=? AND month_of_year=?").
		OrderUp("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid, year, month)
	defer db.Close()
	list := make([]*entity.Menses, 0)
	for db.Next() {
		var m entity.Menses
		m.UserId = uid
		m.CoupleId = cid
		m.Year = year
		m.MonthOfYear = month
		db.Scan(&m.Id, &m.CreateAt, &m.UpdateAt, &m.DayOfMonth, &m.IsStart)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &m)
	}
	return list, nil
}
