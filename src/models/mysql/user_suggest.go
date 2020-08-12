package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddSuggest
func AddSuggest(s *entity.Suggest) (*entity.Suggest, error) {
	s.Status = entity.STATUS_VISIBLE
	s.CreateAt = time.Now().Unix()
	s.UpdateAt = time.Now().Unix()
	//s.Top = false
	s.FollowCount = 0
	s.CommentCount = 0
	db := mysqlDB().
		Insert(TABLE_SUGGEST).
		Set("status=?,create_at=?,update_at=?,user_id=?,app_from=?,official=?,kind=?,title=?,content_text=?,content_image=?,follow_count=?,comment_count=?").
		Exec(s.Status, s.CreateAt, s.UpdateAt, s.UserId, entity.APP_FROM_1, s.Official, s.Kind, s.Title, s.ContentText, s.ContentImage, s.FollowCount, s.CommentCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	s.Id, _ = db.Result().LastInsertId()
	return s, nil
}

// DelSuggest
func DelSuggest(s *entity.Suggest) error {
	s.Status = entity.STATUS_DELETE
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SUGGEST).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(s.Status, s.UpdateAt, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateSuggest
func UpdateSuggest(s *entity.Suggest) (*entity.Suggest, error) {
	s.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_SUGGEST).
		Set("update_at=?,status=?,official=?").
		Where("id=?").
		Exec(s.UpdateAt, s.Status, s.Official, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return s, nil
}

// UpdateSuggestCount
func UpdateSuggestCount(s *entity.Suggest, update bool) (*entity.Suggest, error) {
	if s.FollowCount < 0 {
		s.FollowCount = 0
	}
	if s.CommentCount < 0 {
		s.CommentCount = 0
	}
	if update {
		s.UpdateAt = time.Now().Unix()
	}
	db := mysqlDB().
		Update(TABLE_SUGGEST).
		Set("update_at=?,follow_count=?,comment_count=?").
		Where("id=?").
		Exec(s.UpdateAt, s.FollowCount, s.CommentCount, s.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return s, nil
}

// GetSuggestById
func GetSuggestById(sid int64) (*entity.Suggest, error) {
	var s entity.Suggest
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,kind,title,content_text,content_image,official,follow_count,comment_count").
		Form(TABLE_SUGGEST).
		Where("id=?").
		Query(sid).
		NextScan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.Kind, &s.Title, &s.ContentText, &s.ContentImage, &s.Official, &s.FollowCount, &s.CommentCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if s.Id <= 0 {
		return nil, nil
	} else if s.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &s, nil
}

// GetSuggestListByStatusKind
func GetSuggestListByStatusKind(status, kind, offset, limit int) ([]*entity.Suggest, error) {
	where := ""
	if status > entity.STATUS_VISIBLE {
		where += "status=? AND app_from=?"
	} else {
		status = entity.STATUS_VISIBLE
		where += "status>=? AND app_from=?"
	}
	//hasKind := kind != entity.SUGGEST_KIND_ALL
	//if hasKind {
	//	where += " AND kind=?"
	//}
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,kind,title,content_text,content_image,official,follow_count,comment_count").
		Form(TABLE_SUGGEST).
		Where(where).
		Order("official DESC,update_at DESC").
		Limit(offset, limit).
		Query(status, entity.APP_FROM_1)
	//if hasKind {
	//	db.Query(status, kind)
	//} else {
	//	db.Query(status)
	//}
	defer db.Close()
	list := make([]*entity.Suggest, 0)
	for db.Next() {
		var s entity.Suggest
		db.Scan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.UserId, &s.Kind, &s.Title, &s.ContentText, &s.ContentImage, &s.Official, &s.FollowCount, &s.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

// GetSuggestListByUser
func GetSuggestListByUser(uid int64, offset, limit int) ([]*entity.Suggest, error) {
	db := mysqlDB().
		Select("id,status,create_at,update_at,kind,title,content_text,content_image,official,follow_count,comment_count").
		Form(TABLE_SUGGEST).
		Where("status>=? AND user_id=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid)
	defer db.Close()
	list := make([]*entity.Suggest, 0)
	for db.Next() {
		var s entity.Suggest
		s.UserId = uid
		db.Scan(&s.Id, &s.Status, &s.CreateAt, &s.UpdateAt, &s.Kind, &s.Title, &s.ContentText, &s.ContentImage, &s.Official, &s.FollowCount, &s.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &s)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetSuggestTotalByCreateWithDel
func GetSuggestTotalByCreateWithDel(create int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_SUGGEST).
		Where("create_at>?").
		Query(create).
		NextScan(&total)
	defer db.Close()
	return total
}
