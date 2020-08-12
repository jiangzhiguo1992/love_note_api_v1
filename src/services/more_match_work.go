package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"models/redis"
	"strings"
)

// AddMatchWork
func AddMatchWork(uid, cid int64, mw *entity.MatchWork) (*entity.MatchWork, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if mw == nil {
		return nil, errors.New("nil_work")
	} else if mw.MatchPeriodId <= 0 {
		return nil, errors.New("nil_period")
	}
	// period检查
	mp, err := mysql.GetMatchPeriodById(mw.MatchPeriodId)
	if err != nil {
		return nil, err
	} else if mp == nil {
		return nil, errors.New("nil_period")
	} else if !IsMatchPeriodPlay(mp) {
		return nil, errors.New("period_at_err")
	}
	// 内容检查
	if mp.Kind == entity.MATCH_KIND_WIFE_PICTURE {
		if len(strings.TrimSpace(mw.ContentImage)) <= 0 {
			return nil, errors.New("limit_content_image_nil")
		} else {
			mw.Title = ""
			mw.ContentText = ""
		}
	} else if mp.Kind == entity.MATCH_KIND_LETTER_SHOW {
		if len(strings.TrimSpace(mw.Title)) <= 0 {
			return nil, errors.New("limit_content_text_nil")
		} else if len([]rune(mw.Title)) > GetLimit().MatchWorkTitleLength {
			return nil, errors.New("limit_content_text_over")
		} else {
			mw.ContentImage = ""
			mw.Title = strings.TrimSpace(mw.Title)
			mw.ContentText = ""
		}
	} else if mp.Kind == entity.MATCH_KIND_DISCUSS_MEET {
		if len(strings.TrimSpace(mw.ContentText)) <= 0 {
			return nil, errors.New("limit_content_text_nil")
		} else if len([]rune(mw.ContentText)) > GetLimit().MatchWorkContentLength {
			return nil, errors.New("limit_content_text_over")
		} else {
			mw.ContentImage = ""
			mw.Title = ""
			mw.ContentText = strings.TrimSpace(mw.ContentText)
		}
	}
	// mysql
	mw.UserId = uid
	mw.CoupleId = cid
	mw.Kind = mp.Kind
	mw, err = mysql.AddMatchWork(mw)
	if mw == nil || err != nil {
		return mw, err
	}
	// 同步
	go func() {
		// coin 获取本期作品，包括删除的，防止刷金币
		total := mysql.GetMatchWorkTotalByUserCouplePeriod(uid, cid, mw.MatchPeriodId)
		if total <= 1 {
			coin := &entity.Coin{
				Kind:   entity.COIN_KIND_ADD_BY_MATCH_POST,
				Change: mp.CoinChange,
			}
			AddCoinByFree(uid, cid, coin)
		}
		// period
		mp.WorksCount = mp.WorksCount + 1
		UpdateMatchPeriodCount(mp)
	}()
	// redis-set
	redis.SetMatchWork(mw)
	// TODO redis
	//redis.AddMatchWorkInCreateList(mp.Period, mw)
	return mw, err
}

// DelMatchWork
func DelMatchWork(uid, mwid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if mwid <= 0 {
		return errors.New("nil_work")
	}
	// work检查
	mw, err := GetMatchWorkById(mwid)
	if err != nil {
		return err
	} else if mw == nil {
		return errors.New("nil_work")
	}
	// period检查
	mp, err := mysql.GetMatchPeriodById(mw.MatchPeriodId)
	if err != nil {
		return err
	} else if mp == nil {
		return errors.New("nil_period")
	}
	// admin
	u, _ := GetUserById(uid)
	if !IsAdminister(u) && mw.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelMatchWork(mw)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		mp.WorksCount = mp.WorksCount - 1
		mp.PointCount = mp.PointCount - mw.PointCount
		mp.CoinCount = mp.CoinCount - mw.CoinCount
		mp.ReportCount = mp.ReportCount - mw.ReportCount
		UpdateMatchPeriodCount(mp)
	}()
	// redis-del
	redis.DelMatchWork(mw)
	// TODO redis
	//redis.DelMatchWorkInCoinList(mp.Period, mw)
	//redis.DelMatchWorkInPointList(mp.Period, mw)
	//redis.DelMatchWorkInCreateList(mp.Period, mw)
	return err
}

// UpdateMatchWorkCount
func UpdateMatchWorkCount(mw *entity.MatchWork) (*entity.MatchWork, error) {
	if mw == nil || mw.Id <= 0 {
		return nil, errors.New("nil_work")
	}
	// period检查
	mp, err := mysql.GetMatchPeriodById(mw.MatchPeriodId)
	if err != nil {
		return nil, err
	} else if mp == nil {
		return nil, errors.New("nil_period")
	} else if !IsMatchPeriodPlay(mp) {
		return nil, errors.New("period_at_err")
	}
	// redis-del
	redis.DelMatchWork(mw)
	// mysql
	mw, err = mysql.UpdateMatchWorkCount(mw)
	if mw == nil || err != nil {
		return mw, err
	}
	// redis-set
	redis.SetMatchWork(mw)
	// TODO redis
	//redis.UpdateMatchWorkInCoinList(mp.Period, mw)
	//redis.UpdateMatchWorkInPointList(mp.Period, mw)
	//redis.UpdateMatchWorkInCreateList(mp.Period, mw)
	return mw, err
}

// GetMatchWorkById
func GetMatchWorkById(mwid int64) (*entity.MatchWork, error) {
	if mwid <= 0 {
		return nil, errors.New("nil_work")
	}
	// redis-get
	mw, err := redis.GetMatchWorkById(mwid)
	if mw != nil && mw.Id > 0 && err == nil {
		return mw, err
	}
	// mysql
	mw, err = mysql.GetMatchWorkById(mwid)
	if mw == nil || mw.Id <= 0 || err != nil {
		return mw, err
	}
	// redis-set
	redis.SetMatchWork(mw)
	return mw, err
}

// GetMatchWorkListByPeriodOrder
func GetMatchWorkListByPeriodOrder(uid, cid, mpid int64, order, page int) ([]*entity.MatchWork, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if mpid <= 0 {
		return nil, errors.New("nil_period")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// period检查
	mp, err := mysql.GetMatchPeriodById(mpid)
	if err != nil {
		return nil, err
	} else if mp == nil {
		return nil, errors.New("nil_period")
	} else if !IsMatchPeriodPlay(mp) && order == MATCH_WORK_ORDER_CREATE {
		// 不能查看往期新人榜
		order = MATCH_WORK_ORDER_COIN
	}
	var list []*entity.MatchWork
	limit := GetPageSizeLimit().MatchWork
	offset := page * limit
	// TODO redis
	//if order == MATCH_WORK_ORDER_POINT {
	//	list, _ = redis.GetMatchWorkListPointByKind(mp.Period, mp.Kind, offset, limit)
	//} else if order == MATCH_WORK_ORDER_CREATE {
	//	list, _ = redis.GetMatchWorkListCreateByKind(mp.Period, mp.Kind, offset, limit)
	//} else {
	//	list, _ = redis.GetMatchWorkListCoinByKind(mp.Period, mp.Kind, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		orderBy := MatchWorkOrderList[0]
		if order > 0 && order < len(MatchWorkOrderList) {
			orderBy = MatchWorkOrderList[order]
		}
		maxReportCount := GetLimit().MatchWorkScreenReportCount
		list, err = mysql.GetMatchWorkListByPeriodOrder(mpid, orderBy, maxReportCount, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if order == MATCH_WORK_ORDER_POINT {
		//		redis.SetMatchWorkListPointByKind(mp.Period, mp.Kind, list, true)
		//	} else if order == MATCH_WORK_ORDER_CREATE {
		//		redis.SetMatchWorkListCreateByKind(mp.Period, mp.Kind, list, true)
		//	} else {
		//		redis.SetMatchWorkListCoinByKind(mp.Period, mp.Kind, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_match_work")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		LoadMatchWorkWithAll(uid, cid, v)
	}
	return list, err
}

// GetMatchWorkListByCoupleKind
func GetMatchWorkListByCoupleKind(uid, cid int64, kind, page int) ([]*entity.MatchWork, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if kind != entity.MATCH_KIND_WIFE_PICTURE &&
		kind != entity.MATCH_KIND_LETTER_SHOW &&
		kind != entity.MATCH_KIND_DISCUSS_MEET {
		return nil, errors.New("limit_kind_nil")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().MatchWork
	offset := page * limit
	list, err := mysql.GetMatchWorkListByCoupleKind(cid, kind, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_match_work")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		LoadMatchWorkWithAll(uid, cid, v)
	}
	return list, err
}

// GetMatchWorkList
func GetMatchWorkList(uid, mpid int64, page int) ([]*entity.MatchWork, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().MatchWork
	offset := page * limit
	list, err := mysql.GetMatchWorkList(uid, mpid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_match_work")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetMatchWorkReportList
func GetMatchWorkReportList(page int) ([]*entity.MatchWork, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// reportList
	limit := GetPageSizeLimit().MatchWork
	offset := page * limit
	reportList, err := mysql.GetMatchReportList(offset, limit)
	if err != nil {
		return nil, err
	} else if reportList == nil || len(reportList) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_match_work")
		} else {
			return nil, nil
		}
	}
	// workList
	list := make([]*entity.MatchWork, 0)
	for _, v := range reportList {
		if v == nil || v.MatchWorkId <= 0 {
			continue
		}
		mw, _ := GetMatchWorkById(v.MatchWorkId)
		if mw != nil && mw.Id > 0 {
			list = append(list, mw)
		}
	}
	return list, nil
}

// GetMatchWorkTotalByKindWithDel
func GetMatchWorkTotalByKindWithDel(create int64, kind int) int64 {
	if create == 0 || kind == 0 {
		return 0
	}
	// mysql
	total := mysql.GetMatchWorkTotalByKindWithDel(create, kind)
	return total
}

// LoadMatchWorkWithAll
func LoadMatchWorkWithAll(uid, cid int64, mw *entity.MatchWork) *entity.MatchWork {
	if mw == nil || mw.Id <= 0 {
		return nil
	}
	// 额外属性
	if mw.ReportCount < GetLimit().MatchWorkScreenReportCount {
		mw.Screen = false
	} else {
		mw.Screen = true
	}
	mw.Couple, _ = GetCoupleVisibleByUser(mw.UserId)
	if cid <= 0 {
		// 没配对
		mw.Mine = false
		mw.Our = false
		mw.Report = false
		mw.Point = false
		mw.Coin = false
	} else {
		mw.Mine = mw.UserId == uid
		mw.Our = mw.CoupleId == cid
		//mw.Report = IsMatchWorkReportByUserCouple(uid, cid, mw.Id)
		mw.Point = IsMatchWorkPointByUserCouple(uid, cid, mw.Id)
		//mw.Coin = IsMatchWorkCoinByUserCouple(uid, cid, mw.Id)
	}
	return mw
}
