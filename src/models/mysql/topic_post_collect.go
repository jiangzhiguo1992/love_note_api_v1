package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostCollect
func AddPostCollect(pc *entity.PostCollect) (*entity.PostCollect, error) {
	pc.Status = entity.STATUS_VISIBLE
	pc.CreateAt = time.Now().Unix()
	pc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_COLLECT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_id=?").
		Exec(pc.Status, pc.CreateAt, pc.UpdateAt, pc.UserId, pc.CoupleId, pc.PostId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pc.Id, _ = db.Result().LastInsertId()
	return pc, nil
}

// UpdatePostCollect
func UpdatePostCollect(pc *entity.PostCollect) (*entity.PostCollect, error) {
	pc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST_COLLECT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(pc.Status, pc.UpdateAt, pc.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return pc, nil
}

// GetPostCollectByUserCouple
func GetPostCollectByUserCouple(uid, cid, pid int64) (*entity.PostCollect, error) {
	var pc entity.PostCollect
	pc.UserId = uid
	pc.CoupleId = cid
	pc.PostId = pid
	db := mysqlDB().
		Select("id,status,create_at,update_at").
		Form(TABLE_POST_COLLECT).
		Where("user_id=? AND couple_id=? AND post_id=?").
		Limit(0, 1).
		Query(uid, cid, pid).
		NextScan(&pc.Id, &pc.Status, &pc.CreateAt, &pc.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pc.Id <= 0 {
		return nil, nil
	}
	return &pc, nil
}

// GetPostCollectListByUserCouple
func GetPostCollectListByUserCouple(uid, cid int64, offset, limit int) ([]*entity.PostCollect, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,post_id").
		Form(TABLE_POST_COLLECT).
		Where("status>=? AND user_id=? AND couple_id=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid)
	defer db.Close()
	list := make([]*entity.PostCollect, 0)
	for db.Next() {
		var pc entity.PostCollect
		pc.UserId = uid
		pc.CoupleId = cid
		db.Scan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.PostId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pc)
	}
	return list, nil
}
