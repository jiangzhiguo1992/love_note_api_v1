package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSuggestComment
func AddSuggestComment(sc *entity.SuggestComment) (*entity.SuggestComment, error) {
	sc.Status = entity.STATUS_VISIBLE
	sc.CreateAt = time.Now().Unix()
	sc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_SUGGEST_COMMENT).
		Set("status=?,create_at=?,update_at=?,user_id=?,suggest_id=?,content_text=?,official=?").
		Exec(sc.Status, sc.CreateAt, sc.UpdateAt, sc.UserId, sc.SuggestId, sc.ContentText, sc.Official)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	sc.Id, _ = db.Result().LastInsertId()
	return sc, nil
}

// DelSuggestComment
func DelSuggestComment(sc *entity.SuggestComment) error {
	sc.Status = entity.STATUS_DELETE
	sc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SUGGEST_COMMENT).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(sc.Status, sc.UpdateAt, sc.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetSuggestCommentById
func GetSuggestCommentById(scid int64) (*entity.SuggestComment, error) {
	var sc entity.SuggestComment
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,suggest_id,content_text,official").
		Form(TABLE_SUGGEST_COMMENT).
		Where("id=?").
		Query(scid).
		NextScan(&sc.Id, &sc.Status, &sc.CreateAt, &sc.UpdateAt, &sc.UserId, &sc.SuggestId, &sc.ContentText, &sc.Official)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if sc.Id <= 0 {
		return nil, nil
	} else if sc.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &sc, nil
}

// GetSuggestCommentList
func GetSuggestCommentList(uid, sid int64, offset, limit int) ([]*entity.SuggestComment, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasSuggest := sid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasSuggest {
		where = where + " AND suggest_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,suggest_id,content_text,official").
		Form(TABLE_SUGGEST_COMMENT).
		Where(where)
	if hasSuggest {
		db.Order("official DESC,create_at ASC")
	} else {
		db.OrderDown("create_at")
	}
	db.Limit(offset, limit)
	if !hasUser {
		if !hasSuggest {
			db.Query(entity.STATUS_VISIBLE)
		} else {
			db.Query(entity.STATUS_VISIBLE, sid)
		}
	} else {
		if !hasSuggest {
			db.Query(entity.STATUS_VISIBLE, uid)
		} else {
			db.Query(entity.STATUS_VISIBLE, uid, sid)
		}
	}
	defer db.Close()
	list := make([]*entity.SuggestComment, 0)
	for db.Next() {
		var sc entity.SuggestComment
		db.Scan(&sc.Id, &sc.CreateAt, &sc.UpdateAt, &sc.UserId, &sc.SuggestId, &sc.ContentText, &sc.Official)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &sc)
	}
	return list, nil
}

// GetSuggestCommentTotalByUser
func GetSuggestCommentTotalByUser(uid, sid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SUGGEST_COMMENT).
		Where("status>=? AND user_id=? AND suggest_id=?").
		Query(entity.STATUS_VISIBLE, uid, sid).
		NextScan(&total)
	defer db.Close()
	return total
}

/****************************************** admin ***************************************/

// GetSuggestCommentTotalByCreateWithDel
func GetSuggestCommentTotalByCreateWithDel(create int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SUGGEST_COMMENT).
		Where("create_at>?").
		Query(create).
		NextScan(&total)
	defer db.Close()
	return total
}
