package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddGift
func AddGift(uid, cid int64, g *entity.Gift) (*entity.Gift, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if g == nil {
		return nil, errors.New("nil_gift")
	} else if g.ReceiveId == 0 {
		return nil, errors.New("gift_receive_nil")
	} else if g.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(g.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(g.Title)) > GetLimit().GiftTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// limit
	if len(g.ContentImageList) > 0 {
		imgLimit := GetVipLimitByCouple(cid).GiftImageCount
		if imgLimit <= 0 {
			return nil, errors.New("limit_content_image_refuse")
		} else if len(g.ContentImageList) > imgLimit {
			return nil, errors.New("limit_content_image_over")
		}
	}
	// mysql
	g.UserId = uid
	g.CoupleId = cid
	g, err := mysql.AddGift(g)
	if g == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_GIFT, g.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, g.Id, "push_title_note_update", g.Title, entity.PUSH_TYPE_NOTE_GIFT)
	}()
	return g, err
}

// DelGift
func DelGift(uid, cid, gid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if gid <= 0 {
		return errors.New("nil_gift")
	}
	// 旧数据检查
	g, err := mysql.GetGiftById(gid)
	if err != nil {
		return err
	} else if g == nil {
		return errors.New("nil_gift")
	} else if g.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelGift(g)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_GIFT, gid)
		AddTrends(trends)
	}()
	return err
}

// UpdateGift
func UpdateGift(uid, cid int64, g *entity.Gift) (*entity.Gift, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if g == nil || g.Id <= 0 {
		return nil, errors.New("nil_gift")
	} else if g.ReceiveId == 0 {
		return nil, errors.New("gift_receive_nil")
	} else if g.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(g.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(g.Title)) > GetLimit().GiftTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// 旧数据检查
	old, err := mysql.GetGiftById(g.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_gift")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// 图片检查
	limit := GetVipLimitByCouple(cid).GiftImageCount
	if (len(g.ContentImageList) > limit) && (len(g.ContentImageList) > len(old.ContentImageList)) {
		// 修改的图数大于限制图数，如果是以前vip传上去的，则通过
		return old, errors.New("limit_content_image_over")
	}
	// mysql
	old.ReceiveId = g.ReceiveId
	old.HappenAt = g.HappenAt
	old.Title = g.Title
	old.ContentImageList = g.ContentImageList
	g, err = mysql.UpdateGift(old)
	if g == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_GIFT, g.Id)
		AddTrends(trends)
	}()
	return g, err
}

// GetGiftListByUserCouple
func GetGiftListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Gift, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Gift
	offset := page * limit
	var list []*entity.Gift
	var err error
	if suid > 0 {
		list, err = mysql.GetGiftListByCoupleReceiver(cid, suid, offset, limit)
	} else {
		list, err = mysql.GetGiftListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_gift")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_GIFT)
		AddTrends(trends)
	}()
	return list, err
}
