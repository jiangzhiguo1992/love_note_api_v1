package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddAwardRule
func AddAwardRule(ar *entity.AwardRule) (*entity.AwardRule, error) {
	ar.Status = entity.STATUS_VISIBLE
	ar.CreateAt = time.Now().Unix()
	ar.UpdateAt = time.Now().Unix()
	ar.UseCount = 0
	db := mysqlDB().
		Insert(TABLE_AWARD_RULE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,title=?,score=?,use_count=?").
		Exec(ar.Status, ar.CreateAt, ar.UpdateAt, ar.UserId, ar.CoupleId, ar.Title, ar.Score, ar.UseCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	ar.Id, _ = db.Result().LastInsertId()
	return ar, nil
}

// DelAwardRule
func DelAwardRule(ar *entity.AwardRule) error {
	ar.Status = entity.STATUS_DELETE
	ar.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_AWARD_RULE).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(ar.Status, ar.UpdateAt, ar.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdateAwardRule
func UpdateAwardRule(ar *entity.AwardRule) (*entity.AwardRule, error) {
	if ar.UseCount < 0 {
		ar.UseCount = 0
	}
	ar.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_AWARD_RULE).
		Set("update_at=?,title=?,score=?,use_count=?").
		Where("id=?").
		Exec(ar.UpdateAt, ar.Title, ar.Score, ar.UseCount, ar.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return ar, nil
}

// GetAwardRuleById
func GetAwardRuleById(arid int64) (*entity.AwardRule, error) {
	var ar entity.AwardRule
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,title,score,use_count").
		Form(TABLE_AWARD_RULE).
		Where("id=?").
		Query(arid).
		NextScan(&ar.Id, &ar.Status, &ar.CreateAt, &ar.UpdateAt, &ar.UserId, &ar.CoupleId, &ar.Title, &ar.Score, &ar.UseCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if ar.Id <= 0 {
		return nil, nil
	} else if ar.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &ar, nil
}

// GetAwardRuleListByCouple
func GetAwardRuleListByCouple(cid int64, offset, limit int) ([]*entity.AwardRule, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,title,score,use_count").
		Form(TABLE_AWARD_RULE).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.AwardRule, 0)
	for db.Next() {
		var ar entity.AwardRule
		ar.CoupleId = cid
		db.Scan(&ar.Id, &ar.CreateAt, &ar.UpdateAt, &ar.UserId, &ar.Title, &ar.Score, &ar.UseCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &ar)
	}
	return list, nil
}
