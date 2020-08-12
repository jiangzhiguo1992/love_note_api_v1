package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddGift
func AddGift(g *entity.Gift) (*entity.Gift, error) {
	g.Status = entity.STATUS_VISIBLE
	g.CreateAt = time.Now().Unix()
	g.UpdateAt = time.Now().Unix()
	g.ContentImages = entity.JoinStrByColon(g.ContentImageList)
	db := mysqlDB().
		Insert(TABLE_GIFT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,receive_id=?,happen_at=?,title=?,content_images=?").
		Exec(g.Status, g.CreateAt, g.UpdateAt, g.UserId, g.CoupleId, g.ReceiveId, g.HappenAt, g.Title, g.ContentImages)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	g.Id, _ = db.Result().LastInsertId()
	g.ContentImages = ""
	return g, nil
}

// DelGift
func DelGift(g *entity.Gift) error {
	g.Status = entity.STATUS_DELETE
	g.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_GIFT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(g.Status, g.UpdateAt, g.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateGift
func UpdateGift(g *entity.Gift) (*entity.Gift, error) {
	g.UpdateAt = time.Now().Unix()
	g.ContentImages = entity.JoinStrByColon(g.ContentImageList)
	db := mysqlDB().
		Update(TABLE_GIFT).
		Set("update_at=?,receive_id=?,happen_at=?,title=?,content_images=?").
		Where("id=?").
		Exec(g.UpdateAt, g.ReceiveId, g.HappenAt, g.Title, g.ContentImages, g.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	g.ContentImages = ""
	return g, nil
}

// GetGiftById
func GetGiftById(gid int64) (*entity.Gift, error) {
	var g entity.Gift
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,receive_id,happen_at,title,content_images").
		Form(TABLE_GIFT).
		Where("id=?").
		Query(gid).
		NextScan(&g.Id, &g.Status, &g.CreateAt, &g.UpdateAt, &g.UserId, &g.CoupleId, &g.ReceiveId, &g.HappenAt, &g.Title, &g.ContentImages)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if g.Id <= 0 {
		return nil, nil
	} else if g.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	g.ContentImageList = entity.SplitStrByColon(g.ContentImages)
	g.ContentImages = ""
	return &g, nil
}

// GetGiftListByCouple
func GetGiftListByCouple(cid int64, offset, limit int) ([]*entity.Gift, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,receive_id,happen_at,title,content_images").
		Form(TABLE_GIFT).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Gift, 0)
	for db.Next() {
		var g entity.Gift
		g.CoupleId = cid
		db.Scan(&g.Id, &g.CreateAt, &g.UpdateAt, &g.UserId, &g.ReceiveId, &g.HappenAt, &g.Title, &g.ContentImages)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		g.ContentImageList = entity.SplitStrByColon(g.ContentImages)
		g.ContentImages = ""
		list = append(list, &g)
	}
	return list, nil
}

// GetGiftListByCoupleReceiver
func GetGiftListByCoupleReceiver(cid, rid int64, offset, limit int) ([]*entity.Gift, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_at,title,content_images").
		Form(TABLE_GIFT).
		Where("status>=? AND couple_id=? AND receive_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, rid)
	defer db.Close()
	list := make([]*entity.Gift, 0)
	for db.Next() {
		var g entity.Gift
		g.CoupleId = cid
		g.ReceiveId = rid
		db.Scan(&g.Id, &g.CreateAt, &g.UpdateAt, &g.UserId, &g.HappenAt, &g.Title, &g.ContentImages)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		g.ContentImageList = entity.SplitStrByColon(g.ContentImages)
		g.ContentImages = ""
		list = append(list, &g)
	}
	return list, nil
}

// GetGiftTotalByCouple
func GetGiftTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_GIFT).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
