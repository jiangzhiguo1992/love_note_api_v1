package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddAudio
func AddAudio(uid, cid int64, a *entity.Audio) (*entity.Audio, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil {
		return nil, errors.New("nil_audio")
	} else if a.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(a.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(a.Title)) > GetLimit().AudioTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len(a.ContentAudio) <= 0 {
		return nil, errors.New("limit_content_audio_nil")
	}
	// limit
	//totalLimit := GetVipLimitByCouple(cid).AudioTotalCount
	//if totalLimit <= 0 {
	//	return nil, errors.New("db_add_refuse")
	//} else if mysql.GetAudioTotalByCouple(cid) >= int64(totalLimit) {
	//	return nil, errors.New("limit_total_over")
	//}
	// mysql
	a.UserId = uid
	a.CoupleId = cid
	a, err := mysql.AddAudio(a)
	if a == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_AUDIO, a.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, a.Id, "push_title_note_update", a.Title, entity.PUSH_TYPE_NOTE_AUDIO)
	}()
	return a, err
}

// DelAudio
func DelAudio(uid, cid, aid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if aid <= 0 {
		return errors.New("nil_audio")
	}
	// 旧数据检查
	a, err := mysql.GetAudioById(aid)
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("nil_audio")
	} else if a.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelAudio(a)
	if err != nil {
		return err
	}
	// 动态
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_AUDIO, aid)
		AddTrends(trends)
	}()
	return err
}

// GetAudioListByCouple
func GetAudioListByCouple(uid, cid int64, page int) ([]*entity.Audio, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Audio
	offset := page * limit
	list, err := mysql.GetAudioListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_audio")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_AUDIO)
		AddTrends(trends)
	}()
	return list, err
}
