package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddMatchCoin
func AddMatchCoin(uid, cid int64, mc *entity.MatchCoin) (*entity.MatchCoin, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if mc == nil || mc.MatchWorkId <= 0 {
		return nil, errors.New("nil_work")
	} else if mc.CoinCount <= 0 {
		// 防止负负得正
		return nil, errors.New("nil_coin")
	}
	// work检查
	mw, err := GetMatchWorkById(mc.MatchWorkId)
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
	// coin检查
	coin := &entity.Coin{
		Kind:   entity.COIN_KIND_SUB_BY_MATCH_UP,
		Change: -mc.CoinCount,
	}
	coin, err = AddCoinByFree(uid, cid, coin)
	if coin == nil || coin.Id <= 0 || err != nil {
		return nil, err
	}
	// mysql
	mc.UserId = uid
	mc.CoupleId = cid
	mc.MatchPeriodId = mp.Id
	mc.CoinId = coin.Id
	mc, err = mysql.AddMatchCoin(mc)
	if mc == nil || err != nil {
		return mc, err
	}
	// 同步
	go func() {
		// work
		mw.CoinCount = mw.CoinCount + mc.CoinCount
		UpdateMatchWorkCount(mw)
		// period
		mp.CoinCount = mp.CoinCount + mc.CoinCount
		UpdateMatchPeriodCount(mp)
	}()
	return mc, err
}

// IsMatchWorkCoinByUserCouple
func IsMatchWorkCoinByUserCouple(uid, cid, mwid int64) bool {
	if uid <= 0 || cid <= 0 || mwid <= 0 {
		return false
	}
	coin, _ := mysql.GetMatchCoinByUserCoupleWork(uid, cid, mwid)
	if coin == nil || coin.Id <= 0 {
		return false
	}
	return true
}
