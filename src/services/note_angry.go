package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddAngry
func AddAngry(uid, cid int64, a *entity.Angry) (*entity.Angry, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil {
		return nil, errors.New("nil_angry")
	} else if a.HappenId <= 0 {
		return nil, errors.New("angry_nil_happen_id")
	} else if a.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(a.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(a.ContentText)) > GetLimit().AngryContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// mysql
	a.UserId = uid
	a.CoupleId = cid
	a, err := mysql.AddAngry(a)
	if a == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_ANGRY, a.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, a.Id, "push_title_note_update", a.ContentText, entity.PUSH_TYPE_NOTE_ANGRY)
	}()
	return a, err
}

// DelAngry
func DelAngry(uid, cid, aid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if aid <= 0 {
		return errors.New("nil_angry")
	}
	// 旧数据检查
	a, err := mysql.GetAngryById(aid)
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("nil_angry")
	} else if a.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelAngry(a)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_ANGRY, aid)
		AddTrends(trends)
	}()
	return err
}

// UpdateAngry 主要修改礼物和承诺
func UpdateAngry(uid, cid int64, a *entity.Angry) (*entity.Angry, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil || a.Id <= 0 {
		return nil, errors.New("nil_angry")
	}
	// 旧数据检查
	old, err := mysql.GetAngryById(a.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_angry")
	} else if old.CoupleId != cid {
		// 这里检查的是coupleId
		return old, errors.New("db_update_refuse")
	}
	// gift
	if a.GiftId > 0 {
		gift, _ := mysql.GetGiftById(a.GiftId)
		if gift == nil || gift.CoupleId != cid {
			return old, errors.New("nil_gift")
		} else {
			old.Gift = gift
			old.GiftId = a.GiftId
		}
	} else {
		old.Gift = nil
		old.GiftId = 0
	}
	// promise
	if a.PromiseId > 0 {
		promise, _ := mysql.GetPromiseById(a.PromiseId)
		if promise == nil || promise.CoupleId != cid {
			return old, errors.New("nil_promise")
		} else {
			old.Promise = promise
			old.PromiseId = a.PromiseId
		}
	} else {
		old.Promise = nil
		old.PromiseId = 0
	}
	// mysql
	a, err = mysql.UpdateAngry(old)
	if a == nil || err != nil {
		return old, err
	}
	// 不用load了，关联在上面赋值进去了
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_ANGRY, a.Id)
		AddTrends(trends)
	}()
	return a, err
}

// GetAngryByIdWithGiftPromise
func GetAngryByIdWithGiftPromise(uid, cid, aid int64) (*entity.Angry, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if aid <= 0 {
		return nil, errors.New("nil_angry")
	}
	// mysql
	a, err := mysql.GetAngryById(aid)
	if err != nil {
		return nil, err
	} else if a == nil {
		return nil, errors.New("nil_angry")
	} else if a.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	LoadAngryWithGiftPromise(a)
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_ANGRY, aid)
		AddTrends(trends)
	}()
	return a, err
}

// GetAngryListByUserCouple
func GetAngryListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Angry, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Angry
	offset := page * limit
	var list []*entity.Angry
	var err error
	if suid > 0 {
		list, err = mysql.GetAngryListByCoupleHappenUser(cid, suid, offset, limit)
	} else {
		list, err = mysql.GetAngryListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_angry")
		} else {
			return nil, nil
		}
	}
	// 这里不要额外属性
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_ANGRY)
		AddTrends(trends)
	}()
	return list, err
}

// LoadAngryWithGiftPromise
func LoadAngryWithGiftPromise(a *entity.Angry) *entity.Angry {
	if a == nil {
		return a
	}
	if a.GiftId > 0 {
		a.Gift, _ = mysql.GetGiftById(a.GiftId)
	}
	if a.PromiseId > 0 {
		a.Promise, _ = mysql.GetPromiseById(a.PromiseId)
	}
	return a
}
