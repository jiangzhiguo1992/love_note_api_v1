package mysql

import (
	"database/sql"
	"errors"
	"models/entity"
	"time"
)

// AddCoin
func AddCoin(c *entity.Coin) (*entity.Coin, error) {
	c.Status = entity.STATUS_VISIBLE
	c.CreateAt = time.Now().Unix()
	c.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_COIN).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,kind=?,`change`=?,count=?").
		Exec(c.Status, c.CreateAt, c.UpdateAt, c.UserId, c.CoupleId, c.Kind, c.Change, c.Count)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	c.Id, _ = db.Result().LastInsertId()
	return c, nil
}

// GetCoinById
func GetCoinById(cid int64) (*entity.Coin, error) {
	var c entity.Coin
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,kind,`change`,count").
		Form(TABLE_COIN).
		Where("id=?").
		Limit(0, 1).
		Query(cid).
		NextScan(&c.Id, &c.Status, &c.CreateAt, &c.UpdateAt, &c.UserId, &c.CoupleId, &c.Kind, &c.Change, &c.Count)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	} else if c.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &c, nil
}

// GetCoinLatest
func GetCoinLatest(cid int64) (*entity.Coin, error) {
	var c entity.Coin
	c.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,kind,`change`,count").
		Form(TABLE_COIN).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&c.Id, &c.CreateAt, &c.UpdateAt, &c.UserId, &c.Kind, &c.Change, &c.Count)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	}
	return &c, nil
}

// GetCoinLatestByUserCoupleKind
func GetCoinLatestByUserCoupleKind(uid, cid int64, kind int) (*entity.Coin, error) {
	var c entity.Coin
	c.UserId = uid
	c.CoupleId = cid
	c.Kind = kind
	db := mysqlDB().
		Select("id,create_at,update_at,`change`,count").
		Form(TABLE_COIN).
		Where("status>=? AND user_id=? AND couple_id=? AND kind=?").
		OrderDown("create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, kind).
		NextScan(&c.Id, &c.CreateAt, &c.UpdateAt, &c.Change, &c.Count)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	}
	return &c, nil
}

// GetCoinListByCouple
func GetCoinListByCouple(cid int64, offset, limit int) ([]*entity.Coin, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,kind,`change`,count").
		Form(TABLE_COIN).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Coin, 0)
	for db.Next() {
		var c entity.Coin
		c.CoupleId = cid
		db.Scan(&c.Id, &c.CreateAt, &c.UpdateAt, &c.UserId, &c.Kind, &c.Change, &c.Count)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &c)
	}
	return list, nil
}

// GetCoinCountByCreateUserCoupleKind
func GetCoinCountByCreateUserCoupleKind(create, uid, cid int64, kind int) int {
	var total sql.NullInt64
	db := mysqlDB().
		Select("SUM(`change`) as total").
		Form(TABLE_COIN).
		Where("status>=? AND create_at>=? AND user_id=? AND couple_id=? AND kind=?").
		Query(entity.STATUS_VISIBLE, create, uid, cid, kind).
		NextScan(&total)
	defer db.Close()
	return int(total.Int64)
}

/****************************************** admin ***************************************/

// GetCoinList
func GetCoinList(uid, cid, bid int64, kind int, offset, limit int) ([]*entity.Coin, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasCouple := cid > 0
	hasKind := kind != 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasCouple {
		where = where + " AND couple_id=?"
	}
	if hasKind {
		where = where + " AND kind=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,`change`,count").
		Form(TABLE_COIN).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasCouple {
			if !hasKind {
				db.Query(entity.STATUS_VISIBLE)
			} else {
				db.Query(entity.STATUS_VISIBLE, kind)
			}
		} else {
			if !hasKind {
				db.Query(entity.STATUS_VISIBLE, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, cid, kind)
			}
		}
	} else {
		if !hasCouple {
			if !hasKind {
				db.Query(entity.STATUS_VISIBLE, uid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, kind)
			}
		} else {
			if !hasKind {
				db.Query(entity.STATUS_VISIBLE, uid, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, cid, kind)
			}
		}
	}
	defer db.Close()
	list := make([]*entity.Coin, 0)
	for db.Next() {
		var c entity.Coin
		db.Scan(&c.Id, &c.CreateAt, &c.UpdateAt, &c.UserId, &c.CoupleId, &c.Kind, &c.Change, &c.Count)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &c)
	}
	return list, nil
}

// GetCoinChangeListByCreateWithPay
func GetCoinChangeListByCreateWithPay(start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select("`change`,COUNT(`change`) AS nums").
		Form(TABLE_COIN).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Group("`change`").
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

// GetCoinTotalByCreateKindWithDel
func GetCoinTotalByCreateKindWithDel(start, end int64, kind int) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_COIN).
		Where("(create_at BETWEEN ? AND ?) AND kind=?").
		Query(start, end, kind).
		NextScan(&total)
	defer db.Close()
	return total
}

// GetCoinByBill
//func GetCoinByBill(cid, bid int64) (*entity.Coin, error) {
//	var c entity.Coin
//	c.CoupleId = cid
//	c.BillId = bid
//	db := mysqlDB().
//		Select("id,create_at,update_at,user_id,kind,`change`,count").
//		Form(TABLE_COIN).
//		Where("status>=? AND couple_id=? AND bill_id=?").
//		Limit(0, 1).
//		Query(entity.STATUS_VISIBLE, cid, bid).
//		NextScan(&c.Id, &c.CreateAt, &c.UpdateAt, &c.UserId, &c.Kind, &c.Change, &c.Count)
//	defer db.Close()
//	if db.Err() != nil {
//		return nil, errors.New("db_query_fail")
//	} else if c.Id <= 0 {
//		return nil, nil
//	}
//	return &c, nil
//}
