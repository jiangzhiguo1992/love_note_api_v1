package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddAwardScore
func AddAwardScore(as *entity.AwardScore) (*entity.AwardScore, error) {
	as.Status = entity.STATUS_VISIBLE
	as.CreateAt = time.Now().Unix()
	as.UpdateAt = time.Now().Unix()
	as.ChangeCount = 0
	as.TotalScore = 0
	db := mysqlDB().
		Insert(TABLE_AWARD_SCORE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,change_count=?,total_score=?").
		Exec(as.Status, as.CreateAt, as.UpdateAt, as.UserId, as.CoupleId, as.ChangeCount, as.TotalScore)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	as.Id, _ = db.Result().LastInsertId()
	return as, nil
}

// UpdateAwardScore
func UpdateAwardScore(as *entity.AwardScore) (*entity.AwardScore, error) {
	if as.ChangeCount < 0 {
		as.ChangeCount = 0
	}
	as.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_AWARD_SCORE).
		Set("update_at=?,change_count=?,total_score=?").
		Where("id=?").
		Exec(as.UpdateAt, as.ChangeCount, as.TotalScore, as.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return as, nil
}

// GetAwardScoreByUserCouple
func GetAwardScoreByUserCouple(uid, cid int64) (*entity.AwardScore, error) {
	var as entity.AwardScore
	as.UserId = uid
	as.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,change_count,total_score").
		Form(TABLE_AWARD_SCORE).
		Where("status>=? AND user_id=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&as.Id, &as.CreateAt, &as.UpdateAt, &as.ChangeCount, &as.TotalScore)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if as.Id <= 0 {
		return nil, nil
	}
	return &as, nil
}
