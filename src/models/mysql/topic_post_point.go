package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostPoint
func AddPostPoint(pp *entity.PostPoint) (*entity.PostPoint, error) {
	pp.Status = entity.STATUS_VISIBLE
	pp.CreateAt = time.Now().Unix()
	pp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_POINT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_id=?").
		Exec(pp.Status, pp.CreateAt, pp.UpdateAt, pp.UserId, pp.CoupleId, pp.PostId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pp.Id, _ = db.Result().LastInsertId()
	return pp, nil
}

// UpdatePostPoint
func UpdatePostPoint(pp *entity.PostPoint) (*entity.PostPoint, error) {
	pp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST_POINT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(pp.Status, pp.UpdateAt, pp.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return pp, nil
}

// GetPostPointByUserCouple
func GetPostPointByUserCouple(uid, cid, pid int64) (*entity.PostPoint, error) {
	var pp entity.PostPoint
	pp.UserId = uid
	pp.CoupleId = cid
	pp.PostId = pid
	db := mysqlDB().
		Select("id,status,create_at,update_at").
		Form(TABLE_POST_POINT).
		Where("user_id=? AND couple_id=? AND post_id=?").
		Limit(0, 1).
		Query(uid, cid, pid).
		NextScan(&pp.Id, &pp.Status, &pp.CreateAt, &pp.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pp.Id <= 0 {
		return nil, nil
	}
	return &pp, nil
}
