package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMatchCoin
func AddMatchCoin(mc *entity.MatchCoin) (*entity.MatchCoin, error) {
	mc.Status = entity.STATUS_VISIBLE
	mc.CreateAt = time.Now().Unix()
	mc.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MATCH_COIN).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,match_period_id=?,match_work_id=?,coin_id=?,coin_count=?").
		Exec(mc.Status, mc.CreateAt, mc.UpdateAt, mc.UserId, mc.CoupleId, mc.MatchPeriodId, mc.MatchWorkId, mc.CoinId, mc.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	mc.Id, _ = db.Result().LastInsertId()
	return mc, nil
}

// GetMatchCoinByUserCoupleWork
func GetMatchCoinByUserCoupleWork(uid, cid, mwid int64) (*entity.MatchCoin, error) {
	var mc entity.MatchCoin
	mc.UserId = uid
	mc.CoupleId = cid
	mc.MatchWorkId = mwid
	db := mysqlDB().
		Select("id,create_at,update_at,match_period_id,coin_id,coin_count").
		Form(TABLE_MATCH_COIN).
		Where("status>=? AND user_id=? AND couple_id=? AND match_work_id=?").
		Query(entity.STATUS_VISIBLE, uid, cid, mwid).
		NextScan(&mc.Id, &mc.CreateAt, &mc.UpdateAt, &mc.MatchPeriodId, &mc.CoinId, &mc.CoinCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if mc.Id <= 0 {
		return nil, nil
	}
	return &mc, nil
}
