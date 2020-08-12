package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostCommentReport
func AddPostCommentReport(pcr *entity.PostCommentReport) (*entity.PostCommentReport, error) {
	pcr.Status = entity.STATUS_VISIBLE
	pcr.CreateAt = time.Now().Unix()
	pcr.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_COMMENT_REPORT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_comment_id=?").
		Exec(pcr.Status, pcr.CreateAt, pcr.UpdateAt, pcr.UserId, pcr.CoupleId, pcr.PostCommentId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pcr.Id, _ = db.Result().LastInsertId()
	return pcr, nil
}

// GetPostCommentReportByUserCouple
func GetPostCommentReportByUserCouple(uid, cid, pcid int64) (*entity.PostCommentReport, error) {
	var pr entity.PostCommentReport
	pr.UserId = uid
	pr.CoupleId = cid
	pr.PostCommentId = pcid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_POST_COMMENT_REPORT).
		Where("status>=? AND user_id=? AND couple_id=? AND post_comment_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, pcid).
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

// GetPostCommentReportList
func GetPostCommentReportList(offset, limit int) ([]*entity.PostCommentReport, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,post_comment_id").
		Form(TABLE_POST_COMMENT_REPORT).
		Where("status>=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.PostCommentReport, 0)
	for db.Next() {
		var pr entity.PostCommentReport
		db.Scan(&pr.Id, &pr.CreateAt, &pr.UpdateAt, &pr.UserId, &pr.CoupleId, &pr.PostCommentId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pr)
	}
	return list, nil
}
