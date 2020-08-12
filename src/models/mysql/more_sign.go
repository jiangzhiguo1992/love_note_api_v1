package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSign
func AddSign(s *entity.Sign) (*entity.Sign, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	if s.ContinueDay < 0 {
		s.ContinueDay = 0
	}
	db := mysqlDB().
		Insert(TABLE_SIGN).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,year=?,month_of_year=?,day_of_month=?,continue_day=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, s.UserId, s.CoupleId, s.Year, s.MonthOfYear, s.DayOfMonth, s.ContinueDay)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// GetSignByCoupleYearMonthDay
func GetSignByCoupleYearMonthDay(cid int64, year, month, day int) (*entity.Sign, error) {
	var s entity.Sign
	s.Year = year
	s.MonthOfYear = month
	s.DayOfMonth = day
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,continue_day").
		Form(TABLE_SIGN).
		Where("status>=? AND couple_id=? AND year=? AND month_of_year=? AND day_of_month=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, year, month, day).
		NextScan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.ContinueDay)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	}
	return &s, nil
}

// GetSignListByCoupleYearMonth
func GetSignListByCoupleYearMonth(cid int64, year, month, offset, limit int) ([]*entity.Sign, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,day_of_month,continue_day").
		Form(TABLE_SIGN).
		Where("status>=? AND couple_id=? AND year=? AND month_of_year=?").
		OrderUp("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, year, month)
	defer db.Close()
	list := make([]*entity.Sign, 0)
	for db.Next() {
		var s entity.Sign
		s.CoupleId = cid
		s.Year = year
		s.MonthOfYear = month
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.DayOfMonth, &s.ContinueDay)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetSignList
func GetSignList(uid, cid int64, offset, limit int) ([]*entity.Sign, error) {
	where := "status>=?"
	hasUser := uid > 0
	HasCouple := cid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if HasCouple {
		where = where + " AND couple_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,year,month_of_year,day_of_month,continue_day").
		Form(TABLE_SIGN).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !HasCouple {
			db.Query(entity.STATUS_VISIBLE)
		} else {
			db.Query(entity.STATUS_VISIBLE, cid)
		}
	} else {
		if !HasCouple {
			db.Query(entity.STATUS_VISIBLE, uid)
		} else {
			db.Query(entity.STATUS_VISIBLE, uid, cid)
		}
	}
	defer db.Close()
	list := make([]*entity.Sign, 0)
	for db.Next() {
		var s entity.Sign
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.CoupleId, &s.Year, &s.MonthOfYear, &s.DayOfMonth, &s.ContinueDay)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetSignTotalWithDel
func GetSignTotalWithDel(year, month, day int) int64 {
	var total int64 = 0
	where := ""
	hasYear := year > 0
	hasMonth := month > 0
	hasDay := day > 0
	if hasYear {
		if len(where) <= 0 {
			where = "year=?"
		} else {
			where = where + " AND year=?"
		}
	}
	if hasMonth {
		if len(where) <= 0 {
			where = "month_of_year=?"
		} else {
			where = where + " AND month_of_year=?"
		}
	}
	if hasDay {
		if len(where) <= 0 {
			where = "day_of_month=?"
		} else {
			where = where + " AND day_of_month=?"
		}
	}
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SIGN).
		Where(where)
	if !hasYear {
		if !hasMonth {
			if !hasDay {
				db.Query()
			} else {
				db.Query(day)
			}
		} else {
			if !hasDay {
				db.Query(month)
			} else {
				db.Query(month, day)
			}
		}
	} else {
		if !hasMonth {
			if !hasDay {
				db.Query(year)
			} else {
				db.Query(year, day)
			}
		} else {
			if !hasDay {
				db.Query(year, month)
			} else {
				db.Query(year, month, day)
			}
		}
	}
	db.NextScan(&total)
	defer db.Close()
	return total
}
