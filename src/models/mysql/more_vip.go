package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddVip
func AddVip(v *entity.Vip) (*entity.Vip, error) {
	v.Status = entity.STATUS_VISIBLE
	v.CreateAt = time.Now().Unix()
	v.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_VIP).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,from_type=?,expire_days=?,expire_at=?").
		Exec(v.Status, v.CreateAt, v.UpdateAt, v.UserId, v.CoupleId, v.FromType, v.ExpireDays, v.ExpireAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	v.Id, _ = db.Result().LastInsertId()
	return v, nil
}

// GetVipById
func GetVipById(vid int64) (*entity.Vip, error) {
	var v entity.Vip
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,from_type,expire_days,expire_at").
		Form(TABLE_VIP).
		Where("id=?").
		Limit(0, 1).
		Query(vid).
		NextScan(&v.Id, &v.Status, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.CoupleId, &v.FromType, &v.ExpireDays, &v.ExpireAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if v.Id <= 0 {
		return nil, nil
	} else if v.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &v, nil
}

// GetVipLatest
func GetVipLatest(cid int64) (*entity.Vip, error) {
	var v entity.Vip
	v.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,from_type,expire_days,expire_at").
		Form(TABLE_VIP).
		Where("status>=? AND couple_id=?").
		OrderDown("expire_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.FromType, &v.ExpireDays, &v.ExpireAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if v.Id <= 0 {
		return nil, nil
	}
	return &v, nil
}

// GetVipListByCouple
func GetVipListByCouple(cid int64, offset, limit int) ([]*entity.Vip, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,from_type,expire_days,expire_at").
		Form(TABLE_VIP).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Vip, 0)
	for db.Next() {
		var v entity.Vip
		v.CoupleId = cid
		db.Scan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.FromType, &v.ExpireDays, &v.ExpireAt)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &v)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetVipList
func GetVipList(uid, cid, bid int64, fromType int, offset, limit int) ([]*entity.Vip, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasCouple := cid > 0
	hasFromType := fromType > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasCouple {
		where = where + " AND couple_id=?"
	}
	if hasFromType {
		where = where + " AND from_type=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,from_type,expire_days,expire_at").
		Form(TABLE_VIP).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasCouple {
			if !hasFromType {
				db.Query(entity.STATUS_VISIBLE)
			} else {
				db.Query(entity.STATUS_VISIBLE, fromType)
			}
		} else {
			if !hasFromType {
				db.Query(entity.STATUS_VISIBLE, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, cid, fromType)
			}
		}
	} else {
		if !hasCouple {
			if !hasFromType {
				db.Query(entity.STATUS_VISIBLE, uid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, fromType)
			}
		} else {
			if !hasFromType {
				db.Query(entity.STATUS_VISIBLE, uid, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, cid, fromType)
			}
		}
	}
	defer db.Close()
	list := make([]*entity.Vip, 0)
	for db.Next() {
		var v entity.Vip
		db.Scan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.CoupleId, &v.FromType, &v.ExpireDays, &v.ExpireAt)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &v)
	}
	return list, nil
}

// GetVipExpireDaysListByCreate
func GetVipExpireDaysListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select("expire_days,COUNT(expire_days) AS nums").
		Form(TABLE_VIP).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Group("expire_days").
		OrderDown("nums").
		Query(entity.STATUS_VISIBLE, start, end)
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

// GetVipTotalByCreateWithDel
func GetVipTotalByCreateWithDel(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_VIP).
		Where("create_at BETWEEN ? AND ?").
		Query(start, end).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetVipByBill
//func GetVipByBill(cid, bid int64) (*entity.Vip, error) {
//	var v entity.Vip
//	v.CoupleId = cid
//	v.BillId = bid
//	db := mysqlDB().
//		Select("id,create_at,update_at,user_id,from_type,expire_days,expire_at").
//		Form(TABLE_VIP).
//		Where("status>=? AND couple_id=? AND bill_id=?").
//		Limit(0, 1).
//		Query(entity.STATUS_VISIBLE, cid, bid).
//		NextScan(&v.Id, &v.CreateAt, &v.UpdateAt, &v.UserId, &v.FromType, &v.ExpireDays, &v.ExpireAt)
//	defer db.Close()
//	if db.Err() != nil {
//		return nil, errors.New("db_query_fail")
//	} else if v.Id <= 0 {
//		return nil, nil
//	}
//	return &v, nil
//}
