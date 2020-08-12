package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddEntry
func AddEntry(e *entity.Entry) (*entity.Entry, error) {
	e.Status = entity.STATUS_VISIBLE
	e.CreateAt = time.Now().Unix()
	e.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_ENTRY).
		Set("status=?,create_at=?,update_at=?,user_id=?,platform=?,os_version=?,language=?,device_factory=?,device_name=?,device_id=?,market=?,app_from=?,app_version=?").
		Exec(e.Status, e.CreateAt, e.UpdateAt, e.UserId, e.Platform, e.OsVersion, e.Language, "", e.DeviceName, e.DeviceId, e.Market, entity.APP_FROM_1, e.AppVersion)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	e.Id, _ = db.Result().LastInsertId()
	return e, nil
}

// UpdateEntry
func UpdateEntry(e *entity.Entry) (*entity.Entry, error) {
	e.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_ENTRY).
		Set("update_at=?,device_id=?,device_name=?,market=?,language=?,platform=?,os_version=?,app_version=?").
		Where("id=?").
		Exec(e.UpdateAt, e.DeviceId, e.DeviceName, e.Market, e.Language, e.Platform, e.OsVersion, e.AppVersion, e.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return e, nil
}

// GetEntryLatestByUser
func GetEntryLatestByUser(uid int64) (*entity.Entry, error) {
	var e entity.Entry
	e.UserId = uid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,device_id,device_name,market,language,platform,os_version,app_version").
		Form(TABLE_ENTRY).
		Where("status>=? AND user_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid).
		NextScan(&e.Id, &e.CreateAt, &e.UpdateAt, &e.UserId, &e.DeviceId, &e.DeviceName, &e.Market, &e.Language, &e.Platform, &e.OsVersion, &e.AppVersion)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if e.Id <= 0 {
		return nil, nil
	} else if e.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &e, nil
}

/****************************************** admin ***************************************/

// GetEntryList
func GetEntryList(uid int64, offset, limit int) ([]*entity.Entry, error) {
	hasUser := uid > 0
	where := "status>=?"
	if hasUser {
		where = where + " AND user_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,device_id,device_name,market,language,platform,os_version,app_version").
		Form(TABLE_ENTRY).
		Where(where).
		OrderDown("update_at").
		Limit(offset, limit)
	if !hasUser {
		db.Query(entity.STATUS_VISIBLE)
	} else {
		db.Query(entity.STATUS_VISIBLE, uid)
	}
	defer db.Close()
	list := make([]*entity.Entry, 0)
	for db.Next() {
		var e entity.Entry
		db.Scan(&e.Id, &e.CreateAt, &e.UpdateAt, &e.UserId, &e.DeviceId, &e.DeviceName, &e.Market, &e.Language, &e.Platform, &e.OsVersion, &e.AppVersion)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &e)
	}
	return list, nil
}

// GetEntryGroupListByFiled
func GetEntryGroupListByFiled(filed, at string, start, end int64) ([]*entity.FiledInfo, error) {
	if len(at) <= 0 {
		at = "update_at"
	}
	db := mysqlDB().
		Select(filed + ",COUNT(" + filed + ") AS nums").
		Form(TABLE_ENTRY).
		Where("status>=? AND (" + at + " BETWEEN ? AND ?)").
		Group(filed).
		OrderDown("nums").
		Query(entity.STATUS_VISIBLE, start, end)
	defer db.Close()
	infoList := make([]*entity.FiledInfo, 0)
	for db.Next() {
		var info entity.FiledInfo
		db.Scan(&info.Name, &info.Count)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		infoList = append(infoList, &info)
	}
	return infoList, nil
}

// GetEntryTotalByCreateWithDel
func GetEntryTotalByCreateWithDel(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_ENTRY).
		Where("create_at BETWEEN ? AND ?").
		Query(start, end).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetEntryTotalByUpdateWithDel
func GetEntryTotalByUpdateWithDel(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_ENTRY).
		Where("update_at BETWEEN ? AND ?").
		Query(start, end).
		NextScan(&total)
	defer db.Close()
	return total
}
