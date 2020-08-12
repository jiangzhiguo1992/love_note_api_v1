package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddNotice
func AddNotice(n *entity.Notice) (*entity.Notice, error) {
	if n == nil {
		return nil, errors.New("nil_notice")
	} else if len(strings.TrimSpace(n.Title)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if n.ContentType != entity.NOTICE_TYPE_TEXT &&
		n.ContentType != entity.NOTICE_TYPE_URL &&
		n.ContentType != entity.NOTICE_TYPE_IMAGE {
		return nil, errors.New("limit_kind_nil")
	} else if len(strings.TrimSpace(n.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	}
	// mysql
	n, err := mysql.AddNotice(n)
	return n, err
}

// AddNoticeRead
func AddNoticeRead(uid, nid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if nid <= 0 {
		return errors.New("nil_notice")
	}
	nr, err := mysql.GetNoticeReadByUserNotice(uid, nid)
	if nr != nil || err != nil {
		// 过滤重复，和异常操作
		return nil
	}
	nr = &entity.NoticeRead{
		UserId:   uid,
		NoticeId: nid,
	}
	// mysql
	err = mysql.AddNoticeRead(nr)
	return err
}

// DelNotice
func DelNotice(nid int64) error {
	if nid <= 0 {
		return errors.New("nil_notice")
	}
	// 旧数据检查
	n, err := mysql.GetNoticeById(nid)
	if err != nil {
		return err
	} else if n == nil {
		return errors.New("nil_notice")
	}
	// mysql
	err = mysql.DelNotice(n)
	return err
}

// GetNoticeList
func GetNoticeList(page int) ([]*entity.Notice, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Notice
	offset := page * limit
	list, err := mysql.GetNoticeList(offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_notice")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetNoticeListByUser
func GetNoticeListByUser(uid int64, page int) ([]*entity.Notice, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	}
	list, err := GetNoticeList(page)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		return list, err
	}
	// 额外属性
	for _, n := range list {
		read, _ := mysql.GetNoticeReadByUserNotice(uid, n.Id)
		n.Read = read != nil
	}
	return list, err
}

// GetNoticeCountByNoRead
func GetNoticeCountByNoRead(uid int64) int {
	// 可看到的notice数量
	noticeCount := mysql.GetNoticeCount()
	// 阅读过的notice数量
	readCount := mysql.GetNoticeCountByRead(uid)
	// 没阅读的数量
	noReadCount := noticeCount - readCount
	if noReadCount < 0 {
		noReadCount = 0
	}
	return noReadCount
}
