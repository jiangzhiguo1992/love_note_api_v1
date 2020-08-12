package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSms
func AddSms(s *entity.Sms) (*entity.Sms, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SMS).
		Set("status=?,create_at=?,update_at=?,phone_area=?,phone_number=?,send_type=?,content=?,operator=?,done=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, entity.PHONE_AREA_CHINA, s.Phone, s.SendType, s.Content, "", false)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// GetSmsByPhoneType
func GetSmsByPhoneType(phone string, sendType int) (*entity.Sms, error) {
	var s entity.Sms
	db := mysqlDB().
		Select("id,create_at,update_at,content").
		Form(TABLE_SMS).
		Where("status>=? AND phone_area=? AND phone_number=? AND send_type=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, entity.PHONE_AREA_CHINA, phone, sendType).
		NextScan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.Content)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	}
	return &s, nil
}

// GetSmsList
func GetSmsList(phone string, sendType, offset, limit int) ([]*entity.Sms, error) {
	where := "status>=?"
	hasPhone := len(phone) > 0
	hasType := sendType != 0
	if hasPhone {
		where = where + " AND phone_number=?"
	}
	if hasType {
		where = where + " AND send_type=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,phone_number,send_type,content").
		Form(TABLE_SMS).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasPhone {
		if !hasType {
			db.Query(entity.STATUS_VISIBLE)
		} else {
			db.Query(entity.STATUS_VISIBLE, sendType)
		}
	} else {
		if !hasType {
			db.Query(entity.STATUS_VISIBLE, phone)
		} else {
			db.Query(entity.STATUS_VISIBLE, phone, sendType)
		}
	}
	defer db.Close()
	list := make([]*entity.Sms, 0)
	for db.Next() {
		var s entity.Sms
		db.Scan(&s.Id, &s.CreateAt, &s.UpdateAt, &s.Phone, &s.SendType, &s.Content)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetSmsTotalByCreateWithDel
func GetSmsTotalByCreateWithDel(start, end int64, phone string, sendType int) int64 {
	where := "(create_at BETWEEN ? AND ?)"
	hasPhone := len(phone) > 0
	hasType := sendType != 0
	if hasPhone {
		where = where + " AND phone_number=?"
	}
	if hasType {
		where = where + " AND send_type=?"
	}
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SMS).
		Where(where)
	if !hasPhone {
		if !hasType {
			db.Query(start, end)
		} else {
			db.Query(start, end, sendType)
		}
	} else {
		if !hasType {
			db.Query(start, end, phone)
		} else {
			db.Query(start, end, phone, sendType)
		}
	}
	db.NextScan(&total)
	defer db.Close()
	return total
}
