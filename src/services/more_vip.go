package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"time"
)

// AddVipByPay
func AddVipByPay(uid, cid, bid int64, v *entity.Vip) (*entity.Vip, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if v == nil || v.ExpireDays <= 0 {
		return nil, errors.New("nil_vip")
	} else if bid <= 0 {
		return nil, errors.New("nil_bill")
	}
	// bill
	b, err := mysql.GetBillById(bid)
	if err != nil {
		return nil, err
	} else if b.UserId != uid {
		return nil, errors.New("db_add_refuse")
	} else if b.CoupleId != cid {
		return nil, errors.New("db_add_refuse")
	}
	// 计算会员到期时间
	var bastTime time.Time
	vip, _ := mysql.GetVipLatest(cid)
	if vip != nil && vip.ExpireAt > time.Now().Unix() {
		bastTime = time.Unix(vip.ExpireAt, 0)
	} else {
		bastTime = time.Now()
	}
	v.ExpireAt = bastTime.Add(time.Hour * 24 * time.Duration(v.ExpireDays)).Unix()
	// mysql
	v.UserId = uid
	v.CoupleId = cid
	v.FromType = entity.VIP_FROM_TYPE_USER_BUY
	v, err = mysql.AddVip(v)
	return v, err
}

// AddVipByAdmin
func AddVipByAdmin(v *entity.Vip) (*entity.Vip, error) {
	if v == nil || v.ExpireDays <= 0 {
		return nil, errors.New("nil_vip")
	} else if v.UserId <= 0 {
		return nil, errors.New("nil_user")
	} else if v.CoupleId <= 0 {
		return nil, errors.New("nil_couple")
	}
	// 计算会员到期时间
	var bastTime time.Time
	vip, _ := mysql.GetVipLatest(v.CoupleId)
	if vip != nil && vip.ExpireAt > time.Now().Unix() {
		bastTime = time.Unix(vip.ExpireAt, 0)
	} else {
		bastTime = time.Now()
	}
	v.ExpireAt = bastTime.Add(time.Hour * 24 * time.Duration(v.ExpireDays)).Unix()
	// mysql
	v.FromType = entity.VIP_FROM_TYPE_SYS_SEND
	//v.BillId = 0
	v, err := mysql.AddVip(v)
	return v, err
}

// GetVipLatest
func GetVipLatest(cid int64) (*entity.Vip, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	v, err := mysql.GetVipLatest(cid)
	return v, err
}

// GetVipList
func GetVipList(uid, cid, bid int64, fromType int, page int) ([]*entity.Vip, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Vip
	offset := page * limit
	list, err := mysql.GetVipList(uid, cid, bid, fromType, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_vip")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetVipListByCouple
func GetVipListByCouple(cid int64, page int) ([]*entity.Vip, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Vip
	offset := page * limit
	list, err := mysql.GetVipListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_vip")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性和同步
	return list, err
}

// GetVipExpireDaysListByCreate
func GetVipExpireDaysListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	// mysql
	list, err := mysql.GetVipExpireDaysListByCreate(start, end)
	return list, err
}

// GetVipTotalByCreateWithDel
func GetVipTotalByCreateWithDel(start, end int64) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetVipTotalByCreateWithDel(start, end)
	return total
}

// IsVip
func IsVip(cid int64) bool {
	if cid <= 0 {
		return false
	}
	vip, _ := mysql.GetVipLatest(cid)
	if vip != nil && vip.ExpireAt > time.Now().Unix() {
		return true
	}
	return false
}
