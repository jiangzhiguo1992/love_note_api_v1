package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSouvenirGift
func AddSouvenirGift(sg *entity.SouvenirGift) (*entity.SouvenirGift, error) {
	sg.Status = entity.STATUS_VISIBLE
	sg.CreateAt = time.Now().Unix()
	sg.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SOUVENIR_GIFT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,souvenir_id=?,gift_id=?,year=?").
		Exec(sg.Status, sg.CreateAt, sg.UpdateAt, sg.UserId, sg.CoupleId, sg.SouvenirId, sg.GiftId, sg.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sg.Id, _ = db.Result().LastInsertId()
	return sg, nil
}

// DelSouvenirGift
func DelSouvenirGift(sg *entity.SouvenirGift) error {
	sg.Status = entity.STATUS_DELETE
	sg.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SOUVENIR_GIFT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sg.Status, sg.UpdateAt, sg.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSouvenirGiftById
func GetSouvenirGiftById(sgid int64) (*entity.SouvenirGift, error) {
	var sg entity.SouvenirGift
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,souvenir_id,gift_id,year").
		Form(TABLE_SOUVENIR_GIFT).
		Where("id=?").
		Query(sgid).
		NextScan(&sg.Id, &sg.Status, &sg.CreateAt, &sg.UpdateAt, &sg.UserId, &sg.CoupleId, &sg.SouvenirId, &sg.GiftId, &sg.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sg.Id <= 0 {
		return nil, nil
	} else if sg.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sg, nil
}

// GetSouvenirGiftByCoupleSouvenirGift
func GetSouvenirGiftByCoupleSouvenirGift(cid, sid, gid int64) (*entity.SouvenirGift, error) {
	var sg entity.SouvenirGift
	sg.CoupleId = cid
	sg.SouvenirId = sid
	sg.GiftId = gid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,year").
		Form(TABLE_SOUVENIR_GIFT).
		Where("status>=? AND couple_id=? AND souvenir_id=? AND gift_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, sid, gid).
		NextScan(&sg.Id, &sg.CreateAt, &sg.UpdateAt, &sg.UserId, &sg.Year)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sg.Id <= 0 {
		return nil, nil
	}
	return &sg, nil
}

// GetSouvenirGiftListByCoupleSouvenir
func GetSouvenirGiftListByCoupleSouvenir(cid, sid int64) ([]*entity.SouvenirGift, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,gift_id,year").
		Form(TABLE_SOUVENIR_GIFT).
		Where("status>=? AND couple_id=? AND souvenir_id=?").
		OrderUp("update_at").
		Query(entity.STATUS_VISIBLE, cid, sid)
	defer db.Close()
	list := make([]*entity.SouvenirGift, 0)
	for db.Next() {
		var sg entity.SouvenirGift
		sg.CoupleId = cid
		sg.SouvenirId = sid
		db.Scan(&sg.Id, &sg.CreateAt, &sg.UpdateAt, &sg.UserId, &sg.GiftId, &sg.Year)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sg)
	}
	return list, nil
}
