package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddAward
func AddAward(a *entity.Award) (*entity.Award, error) {
	a.Status = entity.STATUS_VISIBLE
	a.CreateAt = time.Now().Unix()
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_AWARD).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,happen_id=?,award_rule_id=?,happen_at=?,content_text=?,score_change=?").
		Exec(a.Status, a.CreateAt, a.UpdateAt, a.UserId, a.CoupleId, a.HappenId, a.AwardRuleId, a.HappenAt, a.ContentText, a.ScoreChange)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	a.Id, _ = db.Result().LastInsertId()
	return a, nil
}

// DelAward
func DelAward(a *entity.Award) error {
	a.Status = entity.STATUS_DELETE
	a.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_AWARD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(a.Status, a.UpdateAt, a.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// GetAwardById
func GetAwardById(aid int64) (*entity.Award, error) {
	var a entity.Award
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,happen_id,award_rule_id,happen_at,content_text,score_change").
		Form(TABLE_AWARD).
		Where("id=?").
		Query(aid).
		NextScan(&a.Id, &a.Status, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.CoupleId, &a.HappenId, &a.AwardRuleId, &a.HappenAt, &a.ContentText, &a.ScoreChange)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if a.Id <= 0 {
		return nil, nil
	} else if a.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &a, nil
}

// GetAwardListByCouple
func GetAwardListByCouple(cid int64, offset, limit int) ([]*entity.Award, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,happen_id,award_rule_id,happen_at,content_text,score_change").
		Form(TABLE_AWARD).
		Where("status>=? AND couple_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Award, 0)
	for db.Next() {
		var a entity.Award
		a.CoupleId = cid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.HappenId, &a.AwardRuleId, &a.HappenAt, &a.ContentText, &a.ScoreChange)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAwardListByCoupleHappenUser
func GetAwardListByCoupleHappenUser(cid, hid int64, offset, limit int) ([]*entity.Award, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,award_rule_id,happen_at,content_text,score_change").
		Form(TABLE_AWARD).
		Where("status>=? AND couple_id=? AND happen_id=?").
		OrderDown("happen_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, hid)
	defer db.Close()
	list := make([]*entity.Award, 0)
	for db.Next() {
		var a entity.Award
		a.CoupleId = cid
		a.HappenId = hid
		db.Scan(&a.Id, &a.CreateAt, &a.UpdateAt, &a.UserId, &a.AwardRuleId, &a.HappenAt, &a.ContentText, &a.ScoreChange)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &a)
	}
	return list, nil
}

// GetAwardTotalByCouple
func GetAwardTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_AWARD).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
