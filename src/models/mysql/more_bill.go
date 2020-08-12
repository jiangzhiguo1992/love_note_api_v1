package mysql

import (
	"database/sql"
	"errors"
	"models/entity"
	"strings"
	"time"
)

// AddBill
func AddBill(b *entity.Bill) (*entity.Bill, error) {
	b.Status = entity.STATUS_VISIBLE
	b.CreateAt = time.Now().Unix()
	b.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_BILL).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,app_from=?,platform_os=?,platform_pay=?,pay_type=?,pay_amount=?,trade_no=?,trade_receipt=?,goods_type=?,goods_id=?").
		Exec(b.Status, b.CreateAt, b.UpdateAt, b.UserId, b.CoupleId, entity.APP_FROM_1, b.PlatformOs, b.PlatformPay, b.PayType, b.PayAmount, b.TradeNo, b.TradeReceipt, b.GoodsType, b.GoodsId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	b.Id, _ = db.Result().LastInsertId()
	return b, nil
}

// UpdateBill
func UpdateBill(b *entity.Bill) (*entity.Bill, error) {
	b.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_BILL).
		Set("update_at=?,trade_no=?,trade_receipt=?,goods_id=?").
		Where("id=?").
		Exec(b.UpdateAt, b.TradeNo, b.TradeReceipt, b.GoodsId, b.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return b, nil
}

// GetBillById
func GetBillById(bid int64) (*entity.Bill, error) {
	var b entity.Bill
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,platform_os,platform_pay,pay_type,pay_amount,trade_no,trade_receipt,goods_type,goods_id").
		Form(TABLE_BILL).
		Where("id=?").
		Query(bid).
		NextScan(&b.Id, &b.Status, &b.CreateAt, &b.UpdateAt, &b.UserId, &b.CoupleId, &b.PlatformOs, &b.PlatformPay, &b.PayType, &b.PayAmount, &b.TradeNo, &b.TradeReceipt, &b.GoodsType, &b.GoodsId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if b.Id <= 0 {
		return nil, nil
	} else if b.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &b, nil
}

// GetBillByTradeNo
func GetBillByTradeNo(tradeNo, receipt string) (*entity.Bill, error) {
	var b entity.Bill
	b.TradeNo = tradeNo
	b.TradeReceipt = receipt
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,platform_os,platform_pay,pay_type,pay_amount,goods_type,goods_id").
		Form(TABLE_BILL).
		Where("trade_no=? AND trade_receipt=?").
		Limit(0, 1).
		Query(tradeNo, receipt).
		NextScan(&b.Id, &b.Status, &b.CreateAt, &b.UpdateAt, &b.UserId, &b.CoupleId, &b.PlatformOs, &b.PlatformPay, &b.PayType, &b.PayAmount, &b.GoodsType, &b.GoodsId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if b.Id <= 0 {
		return nil, nil
	} else if b.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &b, nil
}

// GetBillListByUserCoupleWithNoSync
func GetBillListByUserCoupleWithNoSync(uid, cid int64, offset, limit int) ([]*entity.Bill, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,platform_os,platform_pay,pay_type,pay_amount,trade_no,trade_receipt,goods_type,goods_id").
		Form(TABLE_BILL).
		Where("status>=? AND user_id=? AND couple_id=? AND goods_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid, 0)
	defer db.Close()
	list := make([]*entity.Bill, 0)
	for db.Next() {
		var b entity.Bill
		b.UserId = uid
		b.CoupleId = cid
		b.GoodsId = 0
		db.Scan(&b.Id, &b.CreateAt, &b.UpdateAt, &b.PlatformOs, &b.PlatformPay, &b.PayType, &b.PayAmount, &b.TradeNo, &b.TradeReceipt, &b.GoodsType, &b.GoodsId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &b)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetBillListWithNoSync
func GetBillListWithNoSync(offset, limit int) ([]*entity.Bill, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,platform_os,platform_pay,pay_type,pay_amount,trade_no,trade_receipt,goods_type,goods_id").
		Form(TABLE_BILL).
		Where("status>=? AND goods_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, 0)
	defer db.Close()
	list := make([]*entity.Bill, 0)
	for db.Next() {
		var b entity.Bill
		db.Scan(&b.Id, &b.CreateAt, &b.UpdateAt, &b.UserId, &b.CoupleId, &b.PlatformOs, &b.PlatformPay, &b.PayType, &b.PayAmount, &b.TradeNo, &b.TradeReceipt, &b.GoodsType, &b.GoodsId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &b)
	}
	return list, nil
}

// GetBillList
func GetBillList(uid, cid int64, tradeNo string, offset, limit int) ([]*entity.Bill, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasCouple := cid > 0
	hasTradeNo := len(strings.TrimSpace(tradeNo)) > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasCouple {
		where = where + " AND couple_id=?"
	}
	if hasTradeNo {
		where = where + " AND trade_no=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,platform_os,platform_pay,pay_type,pay_amount,trade_no,trade_receipt,goods_type,goods_id").
		Form(TABLE_BILL).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasCouple {
			if !hasTradeNo {
				db.Query(entity.STATUS_VISIBLE)
			} else {
				db.Query(entity.STATUS_VISIBLE, tradeNo)
			}
		} else {
			if !hasTradeNo {
				db.Query(entity.STATUS_VISIBLE, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, cid, tradeNo)
			}
		}
	} else {
		if !hasCouple {
			if !hasTradeNo {
				db.Query(entity.STATUS_VISIBLE, uid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, tradeNo)
			}
		} else {
			if !hasTradeNo {
				db.Query(entity.STATUS_VISIBLE, uid, cid)
			} else {
				db.Query(entity.STATUS_VISIBLE, uid, cid, tradeNo)
			}
		}
	}
	defer db.Close()
	list := make([]*entity.Bill, 0)
	for db.Next() {
		var b entity.Bill
		db.Scan(&b.Id, &b.CreateAt, &b.UpdateAt, &b.UserId, &b.CoupleId, &b.PlatformOs, &b.PlatformPay, &b.PayType, &b.PayAmount, &b.TradeNo, &b.TradeReceipt, &b.GoodsType, &b.GoodsId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &b)
	}
	return list, nil
}

// GetBillAmountTotalByCreateWithPay
func GetBillAmountTotalByCreateWithPay(start, end int64, platformOs string, platformPay, payType, goodsType int) float64 {
	where := "status>=? AND (create_at BETWEEN ? AND ?) AND goods_id<>?"
	hasPlatformOs := len(strings.TrimSpace(platformOs)) > 0
	hasPlatformPay := platformPay != 0
	hasPayType := payType != 0
	hasGoodsType := goodsType != 0
	if hasPlatformOs {
		where = where + " AND platform_os=?"
	}
	if hasPlatformPay {
		where = where + " AND platform_pay=?"
	}
	if hasPayType {
		where = where + " AND pay_type=?"
	}
	if hasGoodsType {
		where = where + " AND goods_type=?"
	}
	var total sql.NullFloat64
	db := mysqlDB().
		Select("SUM(pay_amount) as total").
		Form(TABLE_BILL).
		Where(where)
	if !hasPlatformOs {
		if !hasPlatformPay {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, payType)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, payType, goodsType)
				}
			}
		} else {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformPay)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformPay, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformPay, payType)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformPay, payType, goodsType)
				}
			}
		}
	} else {
		if !hasPlatformPay {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, payType)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, payType, goodsType)
				}
			}
		} else {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, platformPay)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, platformPay, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, platformPay, payType)
				} else {
					db.Query(entity.STATUS_VISIBLE, start, end, 0, platformOs, platformPay, payType, goodsType)
				}
			}
		}
	}
	db.NextScan(&total)
	defer db.Close()
	return total.Float64
}

// GetBillTotalByCreateWithDelApy
func GetBillTotalByCreateWithDelApy(start, end int64, platformOs string, platformPay, payType, goodsType int) int64 {
	where := "(create_at BETWEEN ? AND ?) AND goods_id<>?"
	hasPlatformOs := len(strings.TrimSpace(platformOs)) > 0
	hasPlatformPay := platformPay != 0
	hasPayType := payType != 0
	hasGoodsType := goodsType != 0
	if hasPlatformOs {
		where = where + " AND platform_os=?"
	}
	if hasPlatformPay {
		where = where + " AND platform_pay=?"
	}
	if hasPayType {
		where = where + " AND pay_type=?"
	}
	if hasGoodsType {
		where = where + " AND goods_type=?"
	}
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_BILL).
		Where(where)
	if !hasPlatformOs {
		if !hasPlatformPay {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(start, end, 0)
				} else {
					db.Query(start, end, 0, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(start, end, 0, payType)
				} else {
					db.Query(start, end, 0, payType, goodsType)
				}
			}
		} else {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(start, end, 0, platformPay)
				} else {
					db.Query(start, end, 0, platformPay, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(start, end, 0, platformPay, payType)
				} else {
					db.Query(start, end, 0, platformPay, payType, goodsType)
				}
			}
		}
	} else {
		if !hasPlatformPay {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(start, end, 0, platformOs)
				} else {
					db.Query(start, end, 0, platformOs, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(start, end, 0, platformOs, payType)
				} else {
					db.Query(start, end, 0, platformOs, payType, goodsType)
				}
			}
		} else {
			if !hasPayType {
				if !hasGoodsType {
					db.Query(start, end, 0, platformOs, platformPay)
				} else {
					db.Query(start, end, 0, platformOs, platformPay, goodsType)
				}
			} else {
				if !hasGoodsType {
					db.Query(start, end, 0, platformOs, platformPay, payType)
				} else {
					db.Query(start, end, 0, platformOs, platformPay, payType, goodsType)
				}
			}
		}
	}
	db.NextScan(&total)
	defer db.Close()
	return total
}
