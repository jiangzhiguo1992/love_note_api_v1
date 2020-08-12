package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddTrendsBrowse
func AddTrendsBrowse(uid, cid int64) (*entity.TrendsBrowse, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// old
	old, err := mysql.GetTrendsBrowseByUserCouple(uid, cid)
	if err != nil {
		return nil, err
	} else if old == nil || old.Id <= 0 {
		tb := &entity.TrendsBrowse{
			BaseCp: entity.BaseCp{
				UserId:   uid,
				CoupleId: cid,
			},
		}
		old, err = mysql.AddTrendsBrowse(tb)
	} else {
		old, err = mysql.UpdateTrendsBrowse(old)
	}
	if old == nil || err != nil {
		return old, err
	}
	return old, err
}

// GetTrendsCountByUserCouple
func GetTrendsCountByUserCouple(uid, taId, cid int64) int {
	if uid <= 0 {
		return 0
	} else if taId <= 0 {
		return 0
	} else if cid <= 0 {
		return 0
	}
	// trendsBrowseMe
	tbMe, err := mysql.GetTrendsBrowseByUserCouple(uid, cid)
	if err != nil {
		return 0
	} else if tbMe == nil {
		return 0
	}
	// mysql
	return int(mysql.GetTrendsTotalByUpdateUserCouple(tbMe.UpdateAt, taId, cid))
}
