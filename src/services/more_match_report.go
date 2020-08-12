package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddMatchReport
func AddMatchReport(uid, cid int64, mr *entity.MatchReport) (*entity.MatchReport, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if mr == nil || mr.MatchWorkId <= 0 {
		return nil, errors.New("nil_work")
	}
	// work检查
	mw, err := GetMatchWorkById(mr.MatchWorkId)
	if err != nil {
		return nil, err
	} else if mw == nil {
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
	// old
	old, err := mysql.GetMatchReportByUserCoupleWork(uid, cid, mr.MatchWorkId)
	if err != nil {
		return nil, err
	} else if old != nil {
		return nil, errors.New("report_repeat")
	}
	// mysql
	mr.UserId = uid
	mr.CoupleId = cid
	mr.MatchPeriodId = mp.Id
	mr.Reason = ""
	mr, err = mysql.AddMatchReport(mr)
	if mr == nil || err != nil {
		return mr, err
	}
	// 同步
	go func() {
		// work
		mw.ReportCount = mw.ReportCount + 1
		UpdateMatchWorkCount(mw)
		// period
		mp.ReportCount = mp.ReportCount + 1
		UpdateMatchPeriodCount(mp)
	}()
	return mr, err
}

// IsMatchWorkReportByUserCouple
func IsMatchWorkReportByUserCouple(uid, cid, mwid int64) bool {
	if uid <= 0 || cid <= 0 || mwid <= 0 {
		return false
	}
	report, _ := mysql.GetMatchReportByUserCoupleWork(uid, cid, mwid)
	if report == nil || report.Id <= 0 {
		return false
	}
	return true
}
