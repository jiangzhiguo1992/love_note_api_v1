package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddWhisper
func AddWhisper(uid, cid int64, w *entity.Whisper) (*entity.Whisper, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if w == nil {
		return nil, errors.New("nil_whisper")
	} else if len(w.Channel) <= 0 {
		return nil, errors.New("whisper_channel_nil")
	} else if len([]rune(w.Channel)) > GetLimit().WhisperChannelLength {
		return nil, errors.New("whisper_channel_over")
	} else {
		if !w.IsImage {
			if len(w.Content) <= 0 {
				return nil, errors.New("limit_content_text_nil")
			} else if len([]rune(w.Content)) > GetLimit().WhisperContentLength {
				return nil, errors.New("limit_content_text_over")
			}
		} else {
			if !GetVipLimitByCouple(cid).WhisperImageEnable {
				return nil, errors.New("limit_content_image_refuse")
			}
			w.Content = strings.TrimSpace(w.Content)
		}
	}
	// mysql
	w.UserId = uid
	w.CoupleId = cid
	w, err := mysql.AddWhisper(w)
	if w == nil || err != nil {
		return nil, err
	}
	// 同步
	//go func() {
	//	trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_WHISPER, w.Id)
	//	AddTrends(trends)
	//	// 没有推送
	//}()
	return w, err
}

// GetWhisperListByCoupleChannel
func GetWhisperListByCoupleChannel(uid, cid int64, channel string, page int) ([]*entity.Whisper, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if len([]rune(channel)) > GetLimit().WhisperChannelLength {
		return nil, errors.New("whisper_channel_over")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Whisper
	offset := page * limit
	list, err := mysql.GetWhisperListByCoupleChannel(cid, channel, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_whisper")
		} else {
			return nil, nil
		}
	}
	// 同步
	//go func() {
	//	trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_WHISPER)
	//	AddTrends(trends)
	//}()
	return list, err
}
