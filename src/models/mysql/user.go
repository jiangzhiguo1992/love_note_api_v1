package mysql

import (
	"database/sql"
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddUser
func AddUser(u *entity.User) (*entity.User, error) {
	u.Status = entity.STATUS_VISIBLE
	u.CreateAt = time.Now().Unix()
	u.UpdateAt = time.Now().Unix()
	u.Phone = strings.TrimSpace(u.Phone)
	u.Sex = 0      // 性别默认没有，且只能修改一次
	u.Birthday = 0 // 生日默认没有，且只能修改一次
	db := mysqlDB().
		Insert(TABLE_USER).
		Set("status=?,create_at=?,update_at=?,phone_area=?,phone_number=?,password=?,sex=?,birth_type=?,birthday=?,user_token=?").
		Exec(u.Status, u.CreateAt, u.UpdateAt, entity.PHONE_AREA_CHINA, u.Phone, u.Password, u.Sex, entity.USER_BIRTH_LIGHT, u.Birthday, u.UserToken)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("user_register_fail")
	}
	u.Id, _ = db.Result().LastInsertId()
	return u, nil
}

// UpdateUser
func UpdateUser(u *entity.User) (*entity.User, error) {
	u.UpdateAt = time.Now().Unix()
	u.Phone = strings.TrimSpace(u.Phone)
	db := mysqlDB().
		Update(TABLE_USER).
		Set("update_at=?,phone_number=?,password=?,sex=?,birthday=?,user_token=?").
		Where("id=?").
		Exec(u.UpdateAt, u.Phone, u.Password, u.Sex, u.Birthday, u.UserToken, u.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return u, nil
}

// UpdateUserStatus
func UpdateUserStatus(u *entity.User) (*entity.User, error) {
	u.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_USER).
		Set("update_at=?,status=?").
		Where("id=?").
		Exec(u.UpdateAt, u.Status, u.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return u, nil
}

// GetUserByToken
func GetUserByToken(token string) (*entity.User, error) {
	var u entity.User
	u.UserToken = strings.TrimSpace(token)
	db := mysqlDB().
		Select("id,status,create_at,update_at,phone_number,password,sex,birthday").
		Form(TABLE_USER).
		Where("user_token=?").
		Limit(0, 1).
		Query(token).
		NextScan(&u.Id, &u.Status, &u.CreateAt, &u.UpdateAt, &u.Phone, &u.Password, &u.Sex, &u.Birthday)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if u.Id <= 0 {
		return nil, nil
	}
	return &u, nil
}

// GetUserByPhone
func GetUserByPhone(phone string) (*entity.User, error) {
	var u entity.User
	u.Phone = phone
	db := mysqlDB().
		Select("id,status,create_at,update_at,password,sex,birthday,user_token").
		Form(TABLE_USER).
		Where("phone_area=? AND phone_number=?").
		Limit(0, 1).
		Query(entity.PHONE_AREA_CHINA, phone).
		NextScan(&u.Id, &u.Status, &u.CreateAt, &u.UpdateAt, &u.Password, &u.Sex, &u.Birthday, &u.UserToken)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if u.Id <= 0 {
		return nil, nil
	}
	return &u, nil
}

// GetUserById
func GetUserById(uid int64) (*entity.User, error) {
	var u entity.User
	db := mysqlDB().
		Select("id,status,create_at,update_at,phone_number,password,sex,birthday,user_token").
		Form(TABLE_USER).
		Where("id=?").
		Limit(0, 1).
		Query(uid).
		NextScan(&u.Id, &u.Status, &u.CreateAt, &u.UpdateAt, &u.Phone, &u.Password, &u.Sex, &u.Birthday, &u.UserToken)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if u.Id <= 0 {
		return nil, nil
	}
	return &u, nil
}

/****************************************** admin ***************************************/

// GetUserList
func GetUserList(offset, limit int) ([]*entity.User, error) {
	db := mysqlDB().
		Select("id,status,create_at,update_at,phone_number,password,sex,birthday,user_token").
		Form(TABLE_USER).
		OrderDown("create_at").
		Limit(offset, limit).
		Query()
	defer db.Close()
	list := make([]*entity.User, 0)
	for db.Next() {
		var u entity.User
		db.Scan(&u.Id, &u.Status, &u.CreateAt, &u.UpdateAt, &u.Phone, &u.Password, &u.Sex, &u.Birthday, &u.UserToken)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		u.Password = ""
		u.UserToken = ""
		list = append(list, &u)
	}
	return list, nil
}

// GetUserListByBlack
func GetUserListByBlack(offset, limit int) ([]*entity.User, error) {
	db := mysqlDB().
		Select("id,status,create_at,update_at,phone_number,password,sex,birthday,user_token").
		Form(TABLE_USER).
		Where("status<>?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.User, 0)
	for db.Next() {
		var u entity.User
		db.Scan(&u.Id, &u.Status, &u.CreateAt, &u.UpdateAt, &u.Phone, &u.Password, &u.Sex, &u.Birthday, &u.UserToken)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		u.Password = ""
		u.UserToken = ""
		list = append(list, &u)
	}
	return list, nil
}

// GetUserTotalByCreate
func GetUserTotalByCreate(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_USER).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Query(entity.STATUS_VISIBLE, start, end).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetUserTotalByCreate
func GetUserBirthAvgByCreate(start, end int64) float64 {
	var avg sql.NullFloat64
	db := mysqlDB().
		Select("AVG(birthday) as birthday_avg").
		Form(TABLE_USER).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Query(entity.STATUS_VISIBLE, start, end).
		NextScan(&avg)
	defer db.Close()
	return avg.Float64
}
