package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddVideo 添加视频
func AddVideo(uid, cid int64, v *entity.Video) (*entity.Video, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if v == nil {
		return nil, errors.New("nil_video")
	} else if v.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(v.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(v.Title)) > GetLimit().VideoTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len(v.ContentVideo) <= 0 {
		return nil, errors.New("limit_content_video_nil")
	}
	// limit
	//totalLimit := GetVipLimitByCouple(cid).VideoTotalCount
	//if totalLimit <= 0 {
	//	return nil, errors.New("db_add_refuse")
	//} else if mysql.GetVideoTotalByCouple(cid) >= int64(totalLimit) {
	//	return nil, errors.New("limit_total_over")
	//}
	// mysql
	v.UserId = uid
	v.CoupleId = cid
	v, err := mysql.AddVideo(v)
	if v == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_VIDEO, v.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, v.Id, "push_title_note_update", v.Title, entity.PUSH_TYPE_NOTE_VIDEO)
	}()
	return v, err
}

// DelVideo
func DelVideo(uid, cid, vid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if vid <= 0 {
		return errors.New("nil_video")
	}
	// 旧数据检查
	v, err := mysql.GetVideoById(vid)
	if err != nil {
		return err
	} else if v == nil {
		return errors.New("nil_video")
	} else if v.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelVideo(v)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_VIDEO, vid)
		AddTrends(trends)
	}()
	return err
}

// UpdateVideo
func UpdateVideo(uid, cid int64, v *entity.Video) (*entity.Video, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if v == nil || v.Id <= 0 {
		return nil, errors.New("nil_video")
	} else if v.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(v.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(v.Title)) > GetLimit().VideoTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// 旧数据检查
	old, err := mysql.GetVideoById(v.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return nil, errors.New("nil_video")
	} else if old.UserId != uid {
		return nil, errors.New("db_update_refuse")
	}
	// mysql
	old.HappenAt = v.HappenAt
	old.Title = v.Title
	old.Longitude = v.Longitude
	old.Latitude = v.Latitude
	old.Address = v.Address
	old.CityId = v.CityId
	v, err = mysql.UpdateVideo(old)
	if v == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_VIDEO, v.Id)
		AddTrends(trends)
	}()
	return v, err
}

// GetVideoListByCouple
func GetVideoListByCouple(uid, cid int64, page int) ([]*entity.Video, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Video
	offset := page * limit
	list, err := mysql.GetVideoListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_video")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_VIDEO)
		AddTrends(trends)
	}()
	return list, err
}
