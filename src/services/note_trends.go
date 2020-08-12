package services

import (
	"errors"
	"fmt"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// CreateTrends
func CreateTrends(uid int64, cid int64, atp int, ctp int, conId int64) *entity.Trends {
	trends := &entity.Trends{}
	trends.UserId = uid
	trends.CoupleId = cid
	trends.ActionType = atp
	trends.ContentType = ctp
	trends.ContentId = conId
	return trends
}

// CreateTrendsByList
func CreateTrendsByList(uid int64, cid int64, atp int, ctp int) *entity.Trends {
	trends := &entity.Trends{}
	trends.UserId = uid
	trends.CoupleId = cid
	trends.ActionType = atp
	trends.ContentType = ctp
	trends.ContentId = entity.TRENDS_CON_ID_LIST
	return trends
}

// AddTrends
func AddTrends(t *entity.Trends) (*entity.Trends, error) {
	if t == nil {
		utils.LogErr("AddTrends", "缺失 Trends")
		return nil, errors.New("nil_trends")
	}
	if t.UserId <= 0 {
		utils.LogErr("AddTrends", "缺失 UserId "+fmt.Sprintf("%+v", t))
		return nil, errors.New("nil_user")
	} else if t.CoupleId <= 0 {
		utils.LogErr("AddTrends", "缺失 CoupleId "+fmt.Sprintf("%+v", t))
		return nil, errors.New("nil_couple")
	} else if t.ActionType <= 0 {
		utils.LogErr("AddTrends", "缺失 ActionType "+fmt.Sprintf("%+v", t))
		return nil, errors.New("data_err")
	} else if t.ContentType <= 0 {
		utils.LogErr("AddTrends", "缺失 ContentType "+fmt.Sprintf("%+v", t))
		return nil, errors.New("data_err")
	} else if t.ContentId < entity.TRENDS_CON_ID_LIST {
		utils.LogErr("AddTrends", "缺失 ContentId "+fmt.Sprintf("%+v", t))
		return nil, errors.New("data_err")
	}
	// 暂时不要删、查
	if t.ActionType == entity.TRENDS_ACT_TYPE_DELETE ||
		t.ActionType == entity.TRENDS_ACT_TYPE_QUERY {
		return nil, nil
	}
	// old
	old, err := mysql.GetTrendsByUserCoupleTypeId(t.UserId, t.CoupleId, t.ActionType, t.ContentType, t.ContentId)
	if err != nil {
		return nil, err
	} else if old == nil || old.Id <= 0 {
		t, err = mysql.AddTrends(t)
	} else {
		t, err = mysql.UpdateTrends(old)
	}
	if t == nil || err != nil {
		return old, err
	}
	return t, err
}

// GetTrendsListByCreateUserCouple
func GetTrendsListByCreateUserCouple(create, uid, cid int64, page int) ([]*entity.Trends, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	limit := GetPageSizeLimit().Trends
	offset := page * limit
	// mysql
	list, err := mysql.GetTrendsListByCreateCouple(create, cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_trends")
		} else {
			return nil, nil
		}
	}
	// 同步
	go func() {
		AddTrendsBrowse(uid, cid)
	}()
	return list, err
}

// GetTrendsList
func GetTrendsList(uid, cid int64, actType, conType, page int) ([]*entity.Trends, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	limit := GetPageSizeLimit().Trends
	offset := page * limit
	// mysql
	list, err := mysql.GetTrendsList(uid, cid, actType, conType, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_trends")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetTrendsContentTypeList
func GetTrendsContentTypeList(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	// mysql
	list, err := mysql.GetTrendsContentTypeList(start, end)
	return list, err
}

// GetTrendsNoteTotal
func GetTrendsNoteTotal(cid int64) *NoteTotal {
	noteTotal := &NoteTotal{
		TotalSouvenir: mysql.GetSouvenirTotalByCouple(cid),
		TotalWord:     mysql.GetWordTotalByCouple(cid),
		TotalDiary:    mysql.GetDiaryTotalByCouple(cid),
		TotalAward:    mysql.GetAwardTotalByCouple(cid),
		TotalDream:    mysql.GetDreamTotalByCouple(cid),
		TotalGift:     mysql.GetGiftTotalByCouple(cid),
		TotalFood:     mysql.GetFoodTotalByCouple(cid),
		TotalTravel:   mysql.GetTravelTotalByCouple(cid),
		TotalAngry:    mysql.GetAngryTotalByCouple(cid),
		TotalPromise:  mysql.GetPromiseTotalByCouple(cid),
		TotalAudio:    mysql.GetAudioTotalByCouple(cid),
		TotalVideo:    mysql.GetVideoTotalByCouple(cid),
		TotalAlbum:    mysql.GetAlbumTotalByCouple(cid),
		TotalPicture:  mysql.GetPictureTotalByCouple(cid),
		TotalMovie:    mysql.GetMovieTotalByCouple(cid),
	}
	return noteTotal
}
