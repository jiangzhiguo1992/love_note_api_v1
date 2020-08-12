package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddFood
func AddFood(f *entity.Food) (*entity.Food, error) {
	f.Status = entity.STATUS_VISIBLE
	f.CreateAt = time.Now().Unix()
	f.UpdateAt = time.Now().Unix()
	f.ContentImages = entity.JoinStrByColon(f.ContentImageList)
	db := mysqlDB().
		Insert(TABLE_FOOD).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_at=?,title=?,content_images=?,content_text=?,longitude=?,latitude=?,address=?,city_id=?").
		Exec(f.Status, f.CreateAt, f.UpdateAt, f.UserId, f.CoupleId, f.HappenAt, f.Title, f.ContentImages, f.ContentText, f.Longitude, f.Latitude, f.Address, f.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	f.Id, _ = db.Result().LastInsertId()
	f.ContentImages = ""
	return f, nil
}

// DelFood
func DelFood(f *entity.Food) error {
	f.Status = entity.STATUS_DELETE
	f.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_FOOD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(f.Status, f.UpdateAt, f.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateFood
func UpdateFood(f *entity.Food) (*entity.Food, error) {
	f.UpdateAt = time.Now().Unix()
	f.ContentImages = entity.JoinStrByColon(f.ContentImageList)
	db := mysqlDB().
		Update(TABLE_FOOD).
		Set("update_at=?,happen_at=?,title=?,content_images=?,content_text=?,longitude=?,latitude=?,address=?,city_id=?").
		Where("id=?").
		Exec(f.UpdateAt, f.HappenAt, f.Title, f.ContentImages, f.ContentText, f.Longitude, f.Latitude, f.Address, f.CityId, f.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	f.ContentImages = ""
	return f, nil
}

// GetFoodById
func GetFoodById(fid int64) (*entity.Food, error) {
	var f entity.Food
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_at,title,content_images,content_text,longitude,latitude,address,city_id").
		Form(TABLE_FOOD).
		Where("id=?").
		Query(fid).
		NextScan(&f.Id, &f.Status, &f.CreateAt, &f.UpdateAt, &f.UserId, &f.CoupleId, &f.HappenAt, &f.Title, &f.ContentImages, &f.ContentText, &f.Longitude, &f.Latitude, &f.Address, &f.CityId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if f.Id <= 0 {
		return nil, nil
	} else if f.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	f.ContentImageList = entity.SplitStrByColon(f.ContentImages)
	f.ContentImages = ""
	return &f, nil
}

// GetFoodListByCouple
func GetFoodListByCouple(cid int64, offset, limit int) ([]*entity.Food, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,content_images,content_text,longitude,latitude,address,city_id").
		Form(TABLE_FOOD).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Food, 0)
	for db.Next() {
		var f entity.Food
		f.CoupleId = cid
		db.Scan(&f.Id, &f.CreateAt, &f.UpdateAt, &f.UserId, &f.HappenAt, &f.Title, &f.ContentImages, &f.ContentText, &f.Longitude, &f.Latitude, &f.Address, &f.CityId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		f.ContentImageList = entity.SplitStrByColon(f.ContentImages)
		f.ContentImages = ""
		list = append(list, &f)
	}
	return list, nil
}

// GetFoodTotalByCouple
func GetFoodTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_FOOD).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
