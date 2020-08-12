package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddVersion
func AddVersion(v *entity.Version) (*entity.Version, error) {
	v.Status = entity.STATUS_VISIBLE
	v.CreateAt = time.Now().Unix()
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_VERSION).
		Set("status=?,create_at=?,update_at=?,app_from=?,platform=?,version_name=?,version_code=?,update_log=?,update_url=?").
		Exec(v.Status, v.CreateAt, v.UpdateAt, entity.APP_FROM_1, v.Platform, v.VersionName, v.VersionCode, v.UpdateLog, v.UpdateUrl)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	v.Id, _ = db.Result().LastInsertId()
	return v, nil
}

// DelVersion
func DelVersion(v *entity.Version) error {
	v.Status = entity.STATUS_DELETE
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_VERSION).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(v.Status, v.UpdateAt, v.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetVersionById
func GetVersionById(vid int64) (*entity.Version, error) {
	var v entity.Version
	db := mysqlDB().
		Select("id,status,create_at,update_at,platform,version_name,version_code,update_log,update_url").
		Form(TABLE_VERSION).
		Where("id=?").
		Query(vid).
		NextScan(&v.Id, &v.Status, &v.CreateAt, &v.UpdateAt, &v.Platform, &v.VersionName, &v.VersionCode, &v.UpdateLog, &v.UpdateUrl)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if v.Id <= 0 {
		return nil, nil
	} else if v.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &v, nil
}

// GetVersionList
func GetVersionList(offset, limit int) ([]*entity.Version, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,platform,version_name,version_code,update_log,update_url").
		Form(TABLE_VERSION).
		Where("status>=?").
		OrderDown("version_code").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.Version, 0)
	for db.Next() {
		var v entity.Version
		db.Scan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.Platform, &v.VersionName, &v.VersionCode, &v.UpdateLog, &v.UpdateUrl)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &v)
	}
	return list, nil
}

// GetVersionListByCode
func GetVersionListByCode(platform string, code int) ([]*entity.Version, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,version_name,version_code,update_log,update_url").
		Form(TABLE_VERSION).
		Where("status>=? AND platform=? AND version_code>?").
		OrderDown("version_code").
		Query(entity.STATUS_VISIBLE, platform, code)
	defer db.Close()
	list := make([]*entity.Version, 0)
	for db.Next() {
		var v entity.Version
		v.Platform = platform
		db.Scan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.VersionName, &v.VersionCode, &v.UpdateLog, &v.UpdateUrl)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &v)
	}
	return list, nil
}
