package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostCommentPoint
func AddPostCommentPoint(pcp *entity.PostCommentPoint) (*entity.PostCommentPoint, error) {
	pcp.Status = entity.STATUS_VISIBLE
	pcp.CreateAt = time.Now().Unix()
	pcp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_POST_COMMENT_POINT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_comment_id=?").
		Exec(pcp.Status, pcp.CreateAt, pcp.UpdateAt, pcp.UserId, pcp.CoupleId, pcp.PostCommentId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	pcp.Id, _ = db.Result().LastInsertId()
	return pcp, nil
}

// UpdatePostCommentPoint
func UpdatePostCommentPoint(pcp *entity.PostCommentPoint) (*entity.PostCommentPoint, error) {
	pcp.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST_COMMENT_POINT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(pcp.Status, pcp.UpdateAt, pcp.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return pcp, nil
}

// GetPostCommentPointByUserCouple
func GetPostCommentPointByUserCouple(uid, cid, pcid int64) (*entity.PostCommentPoint, error) {
	var pp entity.PostCommentPoint
	pp.UserId = uid
	pp.CoupleId = cid
	pp.PostCommentId = pcid
	db := mysqlDB().
		Select("id,status,create_at,update_at").
		Form(TABLE_POST_COMMENT_POINT).
		Where("user_id=? AND couple_id=? AND post_comment_id=?").
		Limit(0, 1).
		Query(uid, cid, pcid).
		NextScan(&pp.Id, &pp.Status, &pp.CreateAt, &pp.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pp.Id <= 0 {
		return nil, nil
	}
	return &pp, nil
}
