package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddBroadcast
func AddBroadcast(b *entity.Broadcast) (*entity.Broadcast, error) {
	if b == nil {
		return nil, errors.New("nil_broadcast")
	} else if len(strings.TrimSpace(b.Title)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len(strings.TrimSpace(b.Cover)) <= 0 {
		return nil, errors.New("limit_content_image_nil")
	} else if b.StartAt <= 0 || b.EndAt <= 0 || b.StartAt >= b.EndAt {
		return nil, errors.New("limit_happen_err")
	} else if b.ContentType != entity.BROADCAST_TYPE_TEXT &&
		b.ContentType != entity.BROADCAST_TYPE_URL &&
		b.ContentType != entity.BROADCAST_TYPE_IMAGE {
		return nil, errors.New("limit_kind_nil")
	} else if len(strings.TrimSpace(b.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	}
	// mysql
	b.IsEnd = false
	b, err := mysql.AddBroadcast(b)
	return b, err
}

// DelBroadcast
func DelBroadcast(bid int64) error {
	if bid <= 0 {
		return errors.New("nil_broadcast")
	}
	// 旧数据检查
	v, err := mysql.GetBroadcastById(bid)
	if err != nil {
		return err
	} else if v == nil {
		return errors.New("nil_broadcast")
	}
	// mysql
	err = mysql.DelBroadcast(v)
	return err
}

// GetBroadcastList
func GetBroadcastList(page int) ([]*entity.Broadcast, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Broadcast
	offset := page * limit
	list, err := mysql.GetBroadcastList(offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_common")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetBroadcastListNoEnd
func GetBroadcastListNoEnd() ([]*entity.Broadcast, error) {
	// mysql
	limit := GetPageSizeLimit().Broadcast
	list, err := mysql.GetBroadcastListNoEnd(0, limit)
	return list, err
}
