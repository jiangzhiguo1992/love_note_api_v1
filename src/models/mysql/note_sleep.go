package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSleep
func AddSleep(s *entity.Sleep) (*entity.Sleep, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SLEEP).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,year=?,month_of_year=?,day_of_month=?,is_sleep=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, s.UserId, s.CoupleId, s.Year, s.MonthOfYear, s.DayOfMonth, s.IsSleep)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// GetSleepById
func GetSleepById(sid int64) (*entity.Sleep, error) {
	var s entity.Sleep
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,year,month_of_year,day_of_month,is_sleep").
		Form(TABLE_SLEEP).
		Where("id=?").
		Query(sid).
		NextScan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.CoupleId, &s.Year, &s.MonthOfYear, &s.DayOfMonth, &s.IsSleep)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	} else if s.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &s, nil
}

// DelSleep
func DelSleep(s *entity.Sleep) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SLEEP).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(s.Status, s.UpdateAt, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSleepLatestByUserCouple
func GetSleepLatestByUserCouple(uid, cid int64) (*entity.Sleep, error) {
	var s entity.Sleep
	s.UserId = uid
	s.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,year,month_of_year,day_of_month,is_sleep").
		Form(TABLE_SLEEP).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.Year, &s.MonthOfYear, &s.DayOfMonth, &s.IsSleep)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	}
	return &s, nil
}

// GetSleepListByCoupleYearMonth
func GetSleepListByCoupleYearMonth(cid int64, year, month, offset, limit int) ([]*entity.Sleep, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,day_of_month,is_sleep").
		Form(TABLE_SLEEP).
		Where("status>=? AND couple_id=? AND year=? AND month_of_year=?").
		OrderUp("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, year, month)
	defer db.Close()
	list := make([]*entity.Sleep, 0)
	for db.Next() {
		var s entity.Sleep
		s.CoupleId = cid
		s.Year = year
		s.MonthOfYear = month
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.DayOfMonth, &s.IsSleep)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetSleepTotalByUserCoupleDate
func GetSleepTotalByUserCoupleDate(uid, cid int64, year, month, day int) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SLEEP).
		Where("status>=? AND user_id=? AND couple_id=? AND year=? AND month_of_year=? AND day_of_month=?").
		Query(entity.STATUS_VISIBLE, uid, cid, year, month, day).
		NextScan(&total)
	defer db.Close()
	return total
}
