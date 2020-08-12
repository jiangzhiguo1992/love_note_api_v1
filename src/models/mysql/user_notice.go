package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddNotice
func AddNotice(n *entity.Notice) (*entity.Notice, error) {
	n.Status = entity.STATUS_VISIBLE
	n.CreateAt = time.Now().Unix()
	n.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_NOTICE).
		Set("status=?,create_at=?,update_at=?,app_from=?,title=?,content_type=?,content_text=?,read_count=?").
		Exec(n.Status, n.CreateAt, n.UpdateAt, entity.APP_FROM_1, n.Title, n.ContentType, n.ContentText, 0)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	n.Id, _ = db.Result().LastInsertId()
	return n, nil
}

// DelNotice
func DelNotice(n *entity.Notice) error {
	n.Status = entity.STATUS_DELETE
	n.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_NOTICE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(n.Status, n.UpdateAt, n.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetNoticeById
func GetNoticeById(nid int64) (*entity.Notice, error) {
	var n entity.Notice
	db := mysqlDB().
		Select("id,status,create_at,update_at,title,content_type,content_text").
		Form(TABLE_NOTICE).
		Where("id=?").
		Query(nid).
		NextScan(&n.Id, &n.Status, &n.CreateAt, &n.UpdateAt, &n.Title, &n.ContentType, &n.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if n.Id <= 0 {
		return nil, nil
	} else if n.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &n, nil
}

// GetNoticeList
func GetNoticeList(offset, limit int) ([]*entity.Notice, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,title,content_type,content_text").
		Form(TABLE_NOTICE).
		Where("status>=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE)
	defer db.Close()
	list := make([]*entity.Notice, 0)
	for db.Next() {
		var n entity.Notice
		db.Scan(&n.Id, &n.CreateAt, &n.UpdateAt, &n.Title, &n.ContentType, &n.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &n)
	}
	return list, nil
}

// AddNoticeRead
func AddNoticeRead(n *entity.NoticeRead) error {
	n.Status = entity.STATUS_VISIBLE
	n.CreateAt = time.Now().Unix()
	n.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_NOTICE_READ).
		Set("status=?,create_at=?,update_at=?,user_id=?,notice_id=?").
		Exec(n.Status, n.CreateAt, n.UpdateAt, n.UserId, n.NoticeId)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_update_fail")
	}
	return nil
}

// GetNoticeReadByUserNotice
func GetNoticeReadByUserNotice(uid, nid int64) (*entity.NoticeRead, error) {
	var nr entity.NoticeRead
	nr.UserId = uid
	nr.NoticeId = nid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_NOTICE_READ).
		Where("status>=? AND user_id=? AND notice_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, nid).
		NextScan(&nr.Id, &nr.CreateAt, &nr.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if nr.Id <= 0 {
		return nil, nil
	}
	return &nr, nil
}

// GetNoticeCount
func GetNoticeCount() int {
	var noticeCount int
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_NOTICE).
		Where("status>=?").
		Query(entity.STATUS_VISIBLE).
		NextScan(&noticeCount)
	defer db.Close()
	if db.Err() != nil {
		return 0
	}
	return noticeCount
}

// GetNoticeCountByRead
func GetNoticeCountByRead(uid int64) int {
	var readCount int
	db := mysqlDB().
		Select("notice_id").
		Form(TABLE_NOTICE_READ).
		Where("status>=? AND user_id=?").
		Group("notice_id").
		Query(entity.STATUS_VISIBLE, uid)
	defer db.Close()
	if db.Err() != nil {
		return 0
	}
	for db.Next() {
		readCount += 1
	}
	return readCount
}
