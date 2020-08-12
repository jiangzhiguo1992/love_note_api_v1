package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddShy
func AddShy(s *entity.Shy) (*entity.Shy, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SHY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,year=?,month_of_year=?,day_of_month=?,happen_at=?,end_at=?,safe=?,`desc`=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, s.UserId, s.CoupleId, s.Year, s.MonthOfYear, s.DayOfMonth, s.HappenAt, s.EndAt, s.Safe, s.Desc)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// GetShyById
func GetShyById(sid int64) (*entity.Shy, error) {
	var s entity.Shy
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,year,month_of_year,day_of_month,happen_at,end_at,safe,`desc`").
		Form(TABLE_SHY).
		Where("id=?").
		Query(sid).
		NextScan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.CoupleId, &s.Year, &s.MonthOfYear, &s.DayOfMonth, &s.HappenAt, &s.EndAt, &s.Safe, &s.Desc)
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

// DelShy
func DelShy(s *entity.Shy) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SHY).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(s.Status, s.UpdateAt, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetShyListByCoupleYearMonth
func GetShyListByCoupleYearMonth(cid int64, year, month, offset, limit int) ([]*entity.Shy, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,day_of_month,happen_at,end_at,safe,`desc`").
		Form(TABLE_SHY).
		Where("status>=? AND couple_id=? AND year=? AND month_of_year=?").
		OrderUp("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, year, month)
	defer db.Close()
	list := make([]*entity.Shy, 0)
	for db.Next() {
		var s entity.Shy
		s.CoupleId = cid
		s.Year = year
		s.MonthOfYear = month
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.DayOfMonth, &s.HappenAt, &s.EndAt, &s.Safe, &s.Desc)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetShyTotalByUserCoupleDate
func GetShyTotalByUserCoupleDate(uid, cid int64, year, month, day int) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SHY).
		Where("status>=? AND user_id=? AND couple_id=? AND year=? AND month_of_year=? AND day_of_month=?").
		Query(entity.STATUS_VISIBLE, uid, cid, year, month, day).
		NextScan(&total)
	defer db.Close()
	return total
}
