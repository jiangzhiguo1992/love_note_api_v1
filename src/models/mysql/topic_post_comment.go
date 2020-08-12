package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPostComment
func AddPostComment(pc *entity.PostComment) (*entity.PostComment, error) {
	pc.Status = entity.STATUS_VISIBLE
	pc.CreateAt = time.Now().Unix()
	pc.UpdateAt = time.Now().Unix()
	pc.SubCommentCount = 0
	pc.ReportCount = 0
	pc.PointCount = 0
	db := mysqlDB().
		Insert(TABLE_POST_COMMENT).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,post_id=?,to_comment_id=?,floor=?,kind=?,content_text=?,official=?,sub_comment_count=?,report_count=?,point_count=?").
		Exec(pc.Status, pc.CreateAt, pc.UpdateAt, pc.UserId, pc.CoupleId, pc.PostId, pc.ToCommentId, pc.Floor, pc.Kind, pc.ContentText, pc.Official, pc.SubCommentCount, pc.ReportCount, pc.PointCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	pc.Id, _ = db.Result().LastInsertId()
	return pc, nil
}

// DelPostComment
func DelPostComment(pc *entity.PostComment) error {
	pc.Status = entity.STATUS_DELETE
	pc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST_COMMENT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(pc.Status, pc.UpdateAt, pc.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdatePostCount
func UpdatePostCommentCount(pc *entity.PostComment, update bool) (*entity.PostComment, error) {
	if pc.SubCommentCount < 0 {
		pc.SubCommentCount = 0
	}
	if pc.Official || pc.ReportCount < 0 {
		pc.ReportCount = 0
	}
	if pc.PointCount < 0 {
		pc.PointCount = 0
	}
	if update {
		pc.UpdateAt = time.Now().Unix()
	}
	db := mysqlDB().
		Update(TABLE_POST_COMMENT).
		Set("update_at=?,sub_comment_count=?,report_count=?,point_count=?").
		Where("id=?").
		Exec(pc.UpdateAt, pc.SubCommentCount, pc.ReportCount, pc.PointCount, pc.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return pc, nil
}

// GetPostCommentById
func GetPostCommentById(pcid int64) (*entity.PostComment, error) {
	var pc entity.PostComment
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,post_id,to_comment_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("id=?").
		Query(pcid).
		NextScan(&pc.Id, &pc.Status, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.PostId, &pc.ToCommentId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pc.Id <= 0 {
		return nil, nil
	} else if pc.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &pc, nil
}

// GetPostCommentByUserCouplePostCommentKind
func GetPostCommentByUserCouplePostCommentKind(uid, cid, pid, tcid int64, kind int) (*entity.PostComment, error) {
	var pc entity.PostComment
	pc.UserId = uid
	pc.CoupleId = cid
	pc.PostId = pid
	pc.ToCommentId = tcid
	pc.Kind = kind
	db := mysqlDB().
		Select("id,create_at,update_at,floor,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND user_id=? AND couple_id=? AND post_id=? AND to_comment_id=? AND kind=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, pid, tcid, kind).
		NextScan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.Floor, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pc.Id <= 0 {
		return nil, nil
	}
	return &pc, nil
}

// GetPostCommentLatest
func GetPostCommentLatest(pid int64) (*entity.PostComment, error) {
	var pc entity.PostComment
	pc.PostId = pid
	pc.ToCommentId = 0
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND post_id=? AND to_comment_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, pid, 0).
		NextScan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pc.Id <= 0 {
		return nil, nil
	}
	return &pc, nil
}

// GetPostToCommentLatest
func GetPostToCommentLatest(pid, tcid int64) (*entity.PostComment, error) {
	var pc entity.PostComment
	pc.PostId = pid
	pc.ToCommentId = tcid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND post_id=? AND to_comment_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, pid, tcid).
		NextScan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if pc.Id <= 0 {
		return nil, nil
	}
	return &pc, nil
}

// GetPostCommentListByPost
func GetPostCommentListByPost(pid int64, limitReportCount int, order string, offset, limit int) ([]*entity.PostComment, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND post_id=? AND to_comment_id<=? AND report_count<?").
		Order("official DESC," + order).
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, pid, 0, limitReportCount)
	defer db.Close()
	list := make([]*entity.PostComment, 0)
	for db.Next() {
		var pc entity.PostComment
		pc.PostId = pid
		pc.ToCommentId = 0
		pc.Screen = false
		db.Scan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pc)
	}
	return list, nil
}

// GetPostCommentListByUserPost
func GetPostCommentListByUserPost(suid, pid int64, limitReportCount int, order string, offset, limit int) ([]*entity.PostComment, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,couple_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND user_id=? AND post_id=? AND to_comment_id<=? AND report_count<?").
		Order("official DESC," + order).
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, suid, pid, 0, limitReportCount)
	defer db.Close()
	list := make([]*entity.PostComment, 0)
	for db.Next() {
		var pc entity.PostComment
		pc.UserId = suid
		pc.PostId = pid
		pc.ToCommentId = 0
		pc.Screen = false
		db.Scan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.CoupleId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pc)
	}
	return list, nil
}

// GetPostToCommentList
func GetPostToCommentList(pid, tcId int64, limitReportCount int, order string, offset, limit int) ([]*entity.PostComment, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND post_id=? AND to_comment_id=? AND report_count<?").
		Order("official DESC," + order).
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, pid, tcId, limitReportCount)
	defer db.Close()
	list := make([]*entity.PostComment, 0)
	for db.Next() {
		var pc entity.PostComment
		pc.PostId = pid
		pc.ToCommentId = tcId
		pc.Screen = false
		db.Scan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pc)
	}
	return list, nil
}

// GetPostCommentTotalByUserCouple
func GetPostCommentTotalByUserCouple(uid, cid, pid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND user_id=? AND couple_id=? AND post_id=?").
		Query(entity.STATUS_VISIBLE, uid, cid, pid).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetPostToCommentTotalByUserCouple
func GetPostToCommentTotalByUserCouple(uid, cid, pid, tcid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_POST_COMMENT).
		Where("status>=? AND user_id=? AND couple_id=? AND post_id=? AND to_comment_id=?").
		Query(entity.STATUS_VISIBLE, uid, cid, pid, tcid).
		NextScan(&total)
	defer db.Close()
	return total
}

/****************************************** admin ***************************************/

// GetPostCommentList
func GetPostCommentList(uid, pid, tcid int64, offset, limit int) ([]*entity.PostComment, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasPost := pid > 0
	hasToComment := tcid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasPost {
		where = where + " AND post_id=?"
	}
	if hasToComment {
		where = where + " AND to_comment_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,post_id,to_comment_id,floor,kind,content_text,official,sub_comment_count,report_count,point_count").
		Form(TABLE_POST_COMMENT).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasPost {
			if !hasToComment {
				db.Query(entity.STATUS_VISIBLE)
			} else {
				db.Query(entity.STATUS_VISIBLE, tcid)
			}
		} else {
			if !hasToComment {
				db.Query(entity.STATUS_VISIBLE, pid)
			} else {
				db.Query(entity.STATUS_VISIBLE, pid, tcid)
			}
		}
	} else {
		if !hasPost {
			if !hasToComment {
				db.Query(entity.STATUS_VISIBLE, uid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, tcid)
			}
		} else {
			if !hasToComment {
				db.Query(entity.STATUS_VISIBLE, uid, pid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, pid, tcid)
			}
		}
	}
	defer db.Close()
	list := make([]*entity.PostComment, 0)
	for db.Next() {
		var pc entity.PostComment
		db.Scan(&pc.Id, &pc.CreateAt, &pc.UpdateAt, &pc.UserId, &pc.CoupleId, &pc.PostId, &pc.ToCommentId, &pc.Floor, &pc.Kind, &pc.ContentText, &pc.Official, &pc.SubCommentCount, &pc.ReportCount, &pc.PointCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &pc)
	}
	return list, nil
}

// GetPostCommentTotalByCreateWithDel
func GetPostCommentTotalByCreateWithDel(create int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_POST_COMMENT).
		Where("create_at>?").
		Query(create).
		NextScan(&total)
	defer db.Close()
	return total
}
