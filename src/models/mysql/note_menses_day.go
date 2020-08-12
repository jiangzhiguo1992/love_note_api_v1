package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMensesDay
func AddMensesDay(md *entity.MensesDay) (*entity.MensesDay, error) {
	md.Status = entity.STATUS_VISIBLE
	md.CreateAt = time.Now().Unix()
	md.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MENSES_DAY).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,menses2_id=?,year=?,month_of_year=?,day_of_month=?,blood=?,pain=?,mood=?").
		Exec(md.Status, md.CreateAt, md.UpdateAt, md.UserId, md.CoupleId, md.Menses2Id, md.Year, md.MonthOfYear, md.DayOfMonth, md.Blood, md.Pain, md.Mood)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	md.Id, _ = db.Result().LastInsertId()
	return md, nil
}

// UpdateMensesDay
func UpdateMensesDay(md *entity.MensesDay) (*entity.MensesDay, error) {
	md.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES_DAY).
		Set("update_at=?,blood=?,pain=?,mood=?").
		Where("id=?").
		Exec(md.UpdateAt, md.Blood, md.Pain, md.Mood, md.Id)
	defer db.Close()
	if db.Err() != nil {
		return md, errors.New("db_update_fail")
	}
	return md, nil
}

// GetMensesDayByMensesDate
func GetMensesDayByMensesDate(mid int64, year, month, day int) (*entity.MensesDay, error) {
	var md entity.MensesDay
	md.Menses2Id = mid
	md.Year = year
	md.MonthOfYear = month
	md.DayOfMonth = day
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,blood,pain,mood").
		Form(TABLE_MENSES_DAY).
		Where("status>=? AND menses2_id=? AND year=? AND month_of_year=? AND day_of_month=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, mid, year, month, day).
		NextScan(&md.Id, &md.CreateAt, &md.UpdateAt, &md.UserId, &md.CoupleId, &md.Blood, &md.Pain, &md.Mood)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if md.Id <= 0 {
		return nil, nil
	}
	return &md, nil
}

// GetMensesDayListByMenses
func GetMensesDayListByMenses(mid int64, offset, limit int) ([]*entity.MensesDay, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,year,month_of_year,day_of_month,blood,pain,mood").
		Form(TABLE_MENSES_DAY).
		Where("status>=? AND menses2_id=?").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, mid)
	defer db.Close()
	list := make([]*entity.MensesDay, 0)
	for db.Next() {
		var md entity.MensesDay
		md.Menses2Id = mid
		db.Scan(&md.Id, &md.CreateAt, &md.UpdateAt, &md.UserId, &md.CoupleId, &md.Year, &md.MonthOfYear, &md.DayOfMonth, &md.Blood, &md.Pain, &md.Mood)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &md)
	}
	return list, nil
}
