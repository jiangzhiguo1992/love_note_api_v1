package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddApi
func AddApi(a *entity.Api) (*entity.Api, error) {
	a.Status = entity.STATUS_VISIBLE
	a.CreateAt = time.Now().Unix()
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_API).
		Set("status=?,create_at=?,update_at=?,user_id=?,platform=?,language=?,uri=?,method=?,params=?,body=?,result=?,duration=?").
		Exec(a.Status, a.CreateAt, a.UpdateAt, a.UserId, a.Platform, a.Language, a.URI, a.Method, a.Params, a.Body, a.Result, a.Duration)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	a.Id, _ = db.Result().LastInsertId()
	return a, nil
}

// GetApiList
func GetApiList(start, end, uid int64, offset, limit int) ([]*entity.Api, error) {
	where := "status>=?"
	hasCreate := start < end
	hasUser := uid > 0
	if hasCreate {
		where = where + " AND (create_at BETWEEN ? AND ?)"
	}
	if hasUser {
		where = where + " AND user_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,platform,language,uri,method,params,body,result,duration").
		Form(TABLE_API).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasCreate {
		if !hasUser {
			db.Query(entity.STATUS_VISIBLE)
		} else {
			db.Query(entity.STATUS_VISIBLE, uid)
		}
	} else {
		if !hasUser {
			db.Query(entity.STATUS_VISIBLE, start, end)
		} else {
			db.Query(entity.STATUS_VISIBLE, start, end, uid)
		}
	}
	defer db.Close()
	list := make([]*entity.Api, 0)
	for db.Next() {
		var a entity.Api
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.Platform, &a.Language, &a.URI, &a.Method, &a.Params, &a.Body, &a.Result, &a.Duration)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetApiUriListByCreate
func GetApiUriListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select("uri,COUNT(uri) AS nums").
		Form(TABLE_API).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Group("uri").
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

// GetApiTotalByCreateWithDel
func GetApiTotalByCreateWithDel(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_API).
		Where("create_at BETWEEN ? AND ?").
		Query(start, end).
		NextScan(&total)
	defer db.Close()
	return total
}
