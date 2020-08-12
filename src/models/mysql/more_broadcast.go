package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddBroadcast
func AddBroadcast(v *entity.Broadcast) (*entity.Broadcast, error) {
	v.Status = entity.STATUS_VISIBLE
	v.CreateAt = time.Now().Unix()
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_BROADCAST).
		Set("status=?,create_at=?,update_at=?,title=?,cover=?,start_at=?,end_at=?,content_type=?,content_text=?,is_end=?").
		Exec(v.Status, v.CreateAt, v.UpdateAt, v.Title, v.Cover, v.StartAt, v.EndAt, v.ContentType, v.ContentText, v.IsEnd)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	v.Id, _ = db.Result().LastInsertId()
	return v, nil
}

// DelBroadcast
func DelBroadcast(v *entity.Broadcast) error {
	v.Status = entity.STATUS_DELETE
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_BROADCAST).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(v.Status, v.UpdateAt, v.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetBroadcastById
func GetBroadcastById(bid int64) (*entity.Broadcast, error) {
	var b entity.Broadcast
	db := mysqlDB().
		Select("id,status,create_at,update_at,title,cover,start_at,end_at,content_type,content_text,is_end").
		Form(TABLE_BROADCAST).
		Where("id=?").
		Query(bid).
		NextScan(&b.Id, &b.Status, &b.CreateAt, &b.UpdateAt, &b.Title, &b.Cover, &b.StartAt, &b.EndAt, &b.ContentType, &b.ContentText, &b.IsEnd)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if b.Id <= 0 {
		return nil, nil
	} else if b.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &b, nil
}

// GetBroadcastListNoEnd
func GetBroadcastListNoEnd(offset, limit int) ([]*entity.Broadcast, error) {
	nowAt := time.Now().Unix()
	db := mysqlDB().
		Select("id,create_at,update_at,title,cover,start_at,end_at,content_type,content_text").
		Form(TABLE_BROADCAST).
		Where("status>=? AND start_at<=? AND (end_at=? OR end_at>=?) AND is_end=?").
		OrderDown("start_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, nowAt, 0, nowAt, false)
	defer db.Close()
	list := make([]*entity.Broadcast, 0)
	for db.Next() {
		var b entity.Broadcast
		b.IsEnd = false
		db.Scan(&b.Id, &b.CreateAt, &b.UpdateAt, &b.Title, &b.Cover, &b.StartAt, &b.EndAt, &b.ContentType, &b.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &b)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetBroadcastList
func GetBroadcastList(offset, limit int) ([]*entity.Broadcast, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,title,cover,start_at,end_at,content_type,content_text,is_end").
		Form(TABLE_BROADCAST).
		Where("status>=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.Broadcast, 0)
	for db.Next() {
		var b entity.Broadcast
		db.Scan(&b.Id, &b.CreateAt, &b.UpdateAt, &b.Title, &b.Cover, &b.StartAt, &b.EndAt, &b.ContentType, &b.ContentText, &b.IsEnd)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &b)
	}
	return list, nil
}
