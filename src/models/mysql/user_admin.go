package mysql

import (
	"errors"
	"models/entity"
)

// GetAdminByUser
func GetAdminByUser(uid int64) (*entity.Administer, error) {
	var a entity.Administer
	a.UserId = uid
	db := mysqlDB().
		Select("id,status,create_at,update_at,permission").
		Form(TABLE_ADMIN).
		Where("status>=? AND user_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid).
		NextScan(&a.Id, &a.Status, &a.CreateAt, &a.UpdateAt, &a.Permission)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if a.Id <= 0 {
		return nil, nil
	}
	return &a, nil
}
