package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostReport
func AddPostReport(pr *entity.PostReport) (*entity.PostReport, error) {
	pr.Status = entity.STATUS_VISIBLE
	pr.CreateAt = time.Now().Unix()
	pr.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_REPORT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_id=?").
		Exec(pr.Status, pr.CreateAt, pr.UpdateAt, pr.UserId, pr.CoupleId, pr.PostId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pr.Id, _ = db.Result().LastInsertId()
	return pr, nil
}

// GetPostReportByUserCouple
func GetPostReportByUserCouple(uid, cid, pid int64) (*entity.PostReport, error) {
	var pr entity.PostReport
	pr.UserId = uid
	pr.CoupleId = cid
	pr.PostId = pid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_POST_REPORT).
		Where("status>=? AND user_id=? AND couple_id=? AND post_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, pid).
		NextScan(&pr.Id, &pr.CreateAt, &pr.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pr.Id <= 0 {
		return nil, nil
	}
	return &pr, nil
}

/****************************************** admin ***************************************/

// GetPostReportList
func GetPostReportList(offset, limit int) ([]*entity.PostReport, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,post_id").
		Form(TABLE_POST_REPORT).
		Where("status>=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.PostReport, 0)
	for db.Next() {
		var pr entity.PostReport
		db.Scan(&pr.Id, &pr.CreateAt, &pr.UpdateAt, &pr.UserId, &pr.CoupleId, &pr.PostId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pr)
	}
	return list, nil
}
