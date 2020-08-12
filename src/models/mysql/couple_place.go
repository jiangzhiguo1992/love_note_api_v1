package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPlace
func AddPlace(p *entity.Place) (*entity.Place, error) {
	p.Status = entity.STATUS_VISIBLE
	p.CreateAt = time.Now().Unix()
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_PLACE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,longitude=?,latitude=?,address=?,country=?,province=?,city=?,district=?,street=?,city_id=?").
		Exec(p.Status, p.CreateAt, p.UpdateAt, p.UserId, p.CoupleId, p.Longitude, p.Latitude, p.Address, p.Country, p.Province, p.City, p.District, p.Street, p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	p.Id, _ = db.Result().LastInsertId()
	return p, nil
}

// GetPlaceLatestByUserCouple
func GetPlaceLatestByUserCouple(uid, cid int64) (*entity.Place, error) {
	var p entity.Place
	p.UserId = uid
	p.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,longitude,latitude,address,country,province,city,district,street,city_id").
		Form(TABLE_PLACE).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.Longitude, &p.Latitude, &p.Address, &p.Country, &p.Province, &p.City, &p.District, &p.Street, &p.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	}
	return &p, nil
}

// GetPlaceListByCouple
func GetPlaceListByCouple(cid int64, offset, limit int) ([]*entity.Place, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,longitude,latitude,address,country,province,city,district,street,city_id").
		Form(TABLE_PLACE).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Place, 0)
	for db.Next() {
		var p entity.Place
		p.CoupleId = cid
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.Longitude, &p.Latitude, &p.Address, &p.Country, &p.Province, &p.City, &p.District, &p.Street, &p.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &p)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetPlaceList
func GetPlaceList(uid int64, offset, limit int) ([]*entity.Place, error) {
	where := "status>=?"
	hasUser := uid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,longitude,latitude,address,country,province,city,district,street,city_id").
		Form(TABLE_PLACE).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		db.Query(entity.STATUS_VISIBLE)
	} else {
		db.Query(entity.STATUS_VISIBLE, uid)
	}
	defer db.Close()
	list := make([]*entity.Place, 0)
	for db.Next() {
		var p entity.Place
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.Longitude, &p.Latitude, &p.Address, &p.Country, &p.Province, &p.City, &p.District, &p.Street, &p.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &p)
	}
	return list, nil
}

// GetPlaceFilerListByCreate
func GetPlaceFilerListByCreate(filed string, start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select(filed + ",COUNT(" + filed + ") AS nums").
		Form(TABLE_PLACE).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
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
