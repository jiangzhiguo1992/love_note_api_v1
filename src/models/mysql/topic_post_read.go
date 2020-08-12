package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostRead
func AddPostRead(pr *entity.PostRead) (*entity.PostRead, error) {
	pr.Status = entity.STATUS_VISIBLE
	pr.CreateAt = time.Now().Unix()
	pr.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_READ).
		Set("status=?,create_at=?,update_at=?,user_id=?,post_id=?").
		Exec(pr.Status, pr.CreateAt, pr.UpdateAt, pr.UserId, pr.PostId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pr.Id, _ = db.Result().LastInsertId()
	return pr, nil
}

// UpdatePostRead
func UpdatePostRead(pr *entity.PostRead) (*entity.PostRead, error) {
	pr.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST_READ).
		Set("update_at=?").
		Where("id=?").
		Exec(pr.UpdateAt, pr.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return pr, nil
}

// GetPostReadByUser
func GetPostReadByUser(uid, pid int64) (*entity.PostRead, error) {
	var pr entity.PostRead
	pr.UserId = uid
	pr.PostId = pid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_POST_READ).
		Where("status>=? AND user_id=? AND post_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, pid).
		NextScan(&pr.Id, &pr.CreateAt, &pr.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pr.Id <= 0 {
		return nil, nil
	}
	return &pr, nil
}
