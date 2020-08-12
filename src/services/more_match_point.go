package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// ToggleMatchPoint
func ToggleMatchPoint(uid, cid int64, mp *entity.MatchPoint) (*entity.MatchPoint, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if mp == nil || mp.MatchWorkId <= 0 {
		return nil, errors.New("nil_work")
	}
	// work检查
	mw, err := GetMatchWorkById(mp.MatchWorkId)
	if err != nil {
		return nil, err
	} else if mw == nil {
		return nil, errors.New("nil_work")
	}
	// period检查
	period, err := mysql.GetMatchPeriodById(mw.MatchPeriodId)
	if err != nil {
		return nil, err
	} else if period == nil {
		return nil, errors.New("nil_period")
	} else if !IsMatchPeriodPlay(period) {
		return nil, errors.New("period_at_err")
	}
	// mysql
	old, err := mysql.GetMatchPointByUserCoupleWork(uid, cid, mp.MatchWorkId)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		// 没点赞
		mp.UserId = uid
		mp.CoupleId = cid
		mp.MatchPeriodId = period.Id
		mp, err = mysql.AddMatchPoint(mp)
	} else {
		// 已点赞
		if old.Status >= entity.STATUS_VISIBLE {
			old.Status = entity.STATUS_DELETE
		} else {
			old.Status = entity.STATUS_VISIBLE
		}
		mp, err = mysql.UpdateMatchPoint(old)
	}
	if mp == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		if mp.Status >= entity.STATUS_VISIBLE {
			mw.PointCount = mw.PointCount + 1
			period.PointCount = period.PointCount + 1
		} else {
			mw.PointCount = mw.PointCount - 1
			period.PointCount = period.PointCount - 1
		}
		// work
		UpdateMatchWorkCount(mw)
		// period
		UpdateMatchPeriodCount(period)
	}()
	return mp, err
}

// IsMatchWorkPointByUserCouple
func IsMatchWorkPointByUserCouple(uid, cid, mwid int64) bool {
	if uid <= 0 || cid <= 0 || mwid <= 0 {
		return false
	}
	point, _ := mysql.GetMatchPointByUserCoupleWork(uid, cid, mwid)
	if point == nil || point.Id <= 0 {
		return false
	} else if point.Status < entity.STATUS_VISIBLE {
		return false
	}
	return true
}
