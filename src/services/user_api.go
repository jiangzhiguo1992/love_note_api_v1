package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddApi
func AddApi(a *entity.Api) (*entity.Api, error) {
	if a == nil || a.UserId == 0 {
		//utils.LogErr("AddApi", "a == nil")
		return nil, nil
	}
	// 不要管理员的api
	administer, err := mysql.GetAdminByUser(a.UserId)
	if administer != nil {
		return nil, nil
	}
	// mysql
	a, err = mysql.AddApi(a)
	return a, err
}

// GetApiList
func GetApiList(start, end, uid int64, page int) ([]*entity.Api, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Api
	offset := page * limit
	list, err := mysql.GetApiList(start, end, uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_common")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetApiUriListByCreate
func GetApiUriListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	// mysql
	list, err := mysql.GetApiUriListByCreate(start, end)
	return list, err
}

// GetApiTotalByCreateWithDel
func GetApiTotalByCreateWithDel(start, end int64) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetApiTotalByCreateWithDel(start, end)
	return total
}
