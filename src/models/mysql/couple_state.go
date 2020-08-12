package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddCoupleState
func AddCoupleState(cs *entity.CoupleState) (*entity.CoupleState, error) {
	cs.Status = entity.STATUS_VISIBLE
	cs.CreateAt = time.Now().Unix()
	cs.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_COUPLE_STATE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,state=?").
		Exec(cs.Status, cs.CreateAt, cs.UpdateAt, cs.UserId, cs.CoupleId, cs.State)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	cs.Id, _ = db.Result().LastInsertId()
	return cs, nil
}

// GetCoupleStateLatestByState
func GetCoupleStateLatestByState(cid int64, state int) (*entity.CoupleState, error) {
	var cs entity.CoupleState
	cs.CoupleId = cid
	cs.State = state
	db := mysqlDB().
		Select("id,create_at,update_at,user_id").
		Form(TABLE_COUPLE_STATE).
		Where("status>=? AND couple_id=? AND state=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid, state).
		NextScan(&cs.Id, &cs.CreateAt, &cs.UpdateAt, &cs.UserId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if cs.Id <= 0 {
		return nil, nil
	}
	return &cs, nil
}

// GetCoupleStateListByCouple
func GetCoupleStateListByCouple(cid int64, offset, limit int) ([]*entity.CoupleState, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,state").
		Form(TABLE_COUPLE_STATE).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.CoupleState, 0)
	for db.Next() {
		var cs entity.CoupleState
		cs.CoupleId = cid
		db.Scan(&cs.Id, &cs.CreateAt, &cs.UpdateAt, &cs.UserId, &cs.State)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &cs)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetCoupleStateStateListByCreate
func GetCoupleStateStateListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select("state,COUNT(state) AS nums").
		Form(TABLE_COUPLE_STATE).
		Where("create_at BETWEEN ? AND ?").
		Group("state").
		OrderDown("nums").
		Query(start, end)
	defer db.Close()
	infoList := make([]*entity.FiledInfo, 0)
	for db.Next() {
		var info entity.FiledInfo
		db.Scan(&info.Name, &info.Count)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		infoList = append(infoList, &info)
	}
	return infoList, nil
}
