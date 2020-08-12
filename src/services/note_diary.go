package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddDiary
func AddDiary(uid, cid int64, d *entity.Diary) (*entity.Diary, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if d == nil {
		return nil, errors.New("nil_diary")
	} else if d.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len([]rune(d.ContentText)) > GetLimit().DiaryContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// limit
	if len(d.ContentImageList) > 0 {
		imgLimit := GetVipLimitByCouple(cid).DiaryImageCount
		if imgLimit <= 0 {
			return nil, errors.New("limit_content_image_refuse")
		} else if len(d.ContentImageList) > imgLimit {
			return nil, errors.New("limit_content_image_over")
		}
	}
	// mysql
	d.UserId = uid
	d.CoupleId = cid
	d, err := mysql.AddDiary(d)
	if d == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_DIARY, d.Id)
		_, _ = AddTrends(trends)
		// push
		AddPushInCouple(uid, d.Id, "push_title_note_update", d.ContentText, entity.PUSH_TYPE_NOTE_DIARY)
	}()
	return d, err
}

// DelDiary
func DelDiary(uid, cid, did int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if did <= 0 {
		return errors.New("nil_diary")
	}
	// 旧数据检查
	d, err := mysql.GetDiaryById(did)
	if err != nil {
		return err
	} else if d == nil {
		return errors.New("nil_diary")
	} else if d.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelDiary(d)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_DIARY, did)
		_, _ = AddTrends(trends)
	}()
	return err
}

// UpdateDiary
func UpdateDiary(uid, cid int64, d *entity.Diary) (*entity.Diary, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if d == nil || d.Id <= 0 {
		return nil, errors.New("nil_diary")
	} else if d.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len([]rune(d.ContentText)) > GetLimit().DiaryContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 旧数据检查
	old, err := mysql.GetDiaryById(d.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_diary")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// limit
	limit := GetVipLimitByCouple(cid).DiaryImageCount
	if (len(d.ContentImageList) > limit) && (len(d.ContentImageList) > len(old.ContentImageList)) {
		// 修改的图数大于限制图数，如果是以前vip传上去的，则通过
		return old, errors.New("limit_content_image_over")
	}
	// mysql
	old.HappenAt = d.HappenAt
	old.ContentText = d.ContentText
	old.ContentImageList = d.ContentImageList
	d, err = mysql.UpdateDiary(old)
	if d == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_DIARY, d.Id)
		AddTrends(trends)
	}()
	return d, err
}

// GetDiaryById
func GetDiaryById(uid, cid, did int64) (*entity.Diary, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if did <= 0 {
		return nil, errors.New("nil_diary")
	}
	// mysql
	d, err := mysql.GetDiaryById(did)
	if err != nil {
		return nil, err
	} else if d == nil {
		return nil, errors.New("nil_diary")
	} else if d.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	d.ReadCount += 1
	// 同步
	go func() {
		// readCount
		mysql.UpdateDiaryReadCount(d)
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_DIARY, did)
		AddTrends(trends)
	}()
	return d, err
}

// GetDiaryListByUserCouple
func GetDiaryListByUserCouple(mid, suid, cid int64, page int) ([]*entity.Diary, error) {
	if mid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Diary
	offset := page * limit
	var list []*entity.Diary
	var err error
	if suid > 0 {
		list, err = mysql.GetDiaryListByUserCouple(suid, cid, offset, limit)
	} else {
		list, err = mysql.GetDiaryListByCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_diary")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(mid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_DIARY)
		AddTrends(trends)
	}()
	return list, err
}
