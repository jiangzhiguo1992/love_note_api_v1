package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSuggestFollow
func AddSuggestFollow(sf *entity.SuggestFollow) (*entity.SuggestFollow, error) {
	sf.Status = entity.STATUS_VISIBLE
	sf.CreateAt = time.Now().Unix()
	sf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SUGGEST_FOLLOW).
		Set("status=?,create_at=?,update_at=?,user_id=?,suggest_id=?").
		Exec(sf.Status, sf.CreateAt, sf.UpdateAt, sf.UserId, sf.SuggestId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	sf.Id, _ = db.Result().LastInsertId()
	return sf, nil
}

// UpdateSuggestFollowStatus
func UpdateSuggestFollowStatus(sf *entity.SuggestFollow) (*entity.SuggestFollow, error) {
	sf.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SUGGEST_FOLLOW).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sf.Status, sf.UpdateAt, sf.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return sf, nil
}

// GetSuggestFollowByUser
func GetSuggestFollowByUser(uid, sid int64) (*entity.SuggestFollow, error) {
	var sf entity.SuggestFollow
	sf.UserId = uid
	sf.SuggestId = sid
	db := mysqlDB().
		Select("id,status,create_at,update_at").
		Form(TABLE_SUGGEST_FOLLOW).
		Where("user_id=? AND suggest_id=?").
		Limit(0, 1).
		Query(uid, sid).
		NextScan(&sf.Id, &sf.Status, &sf.CreateAt, &sf.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sf.Id <= 0 {
		return nil, nil
	}
	return &sf, nil
}

// GetSuggestFollowListByUser
func GetSuggestFollowListByUser(uid int64, offset, limit int) ([]*entity.SuggestFollow, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,suggest_id").
		Form(TABLE_SUGGEST_FOLLOW).
		Where("status>=? AND user_id=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid)
	defer db.Close()
	list := make([]*entity.SuggestFollow, 0)
	for db.Next() {
		var sf entity.SuggestFollow
		sf.UserId = uid
		db.Scan(&sf.Id, &sf.CreateAt, &sf.UpdateAt, &sf.SuggestId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sf)
	}
	return list, nil
}
