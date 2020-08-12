package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddCouple
func AddCouple(c *entity.Couple) (*entity.Couple, error) {
	now := time.Now().Unix()
	c.Status = entity.STATUS_VISIBLE
	c.CreateAt = now
	c.UpdateAt = now
	c.TogetherAt = now
	c.CreatorName = ""
	c.InviteeName = ""
	c.CreatorAvatar = ""
	c.InviteeAvatar = ""
	// 开始事务
	db := mysqlTX().
		Insert(TABLE_COUPLE).
		Set("status=?,create_at=?,update_at=?,together_at=?,creator_id=?,invitee_id=?,creator_name=?,invitee_name=?,creator_avatar=?,invitee_avatar=?").
		Exec(c.Status, c.CreateAt, c.UpdateAt, c.TogetherAt, c.CreatorId, c.InviteeId, c.CreatorName, c.InviteeName, c.CreatorAvatar, c.InviteeAvatar)
	defer db.Close()
	if db.Err() != nil {
		db.tx.Rollback() // 回滚
		return nil, errors.New("db_add_fail")
	}
	var err error
	c.Id, err = db.Result().LastInsertId()
	if c.Id <= 0 || err != nil {
		db.tx.Rollback() // 回滚
		return nil, errors.New("db_add_fail")
	}
	// state操作
	cs := &entity.CoupleState{
		BaseObj: entity.BaseObj{
			Status:   entity.STATUS_VISIBLE,
			CreateAt: now,
			UpdateAt: now,
		},
		BaseCp: entity.BaseCp{
			UserId:   c.CreatorId,
			CoupleId: c.Id,
		},
		State: entity.COUPLE_STATE_INVITE,
	}
	db2 := mysqlTX2(db.tx).
		Insert(TABLE_COUPLE_STATE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,state=?").
		Exec(cs.Status, cs.CreateAt, cs.UpdateAt, cs.UserId, cs.CoupleId, cs.State)
	defer db2.Close()
	if db2.Err() != nil {
		db2.tx.Rollback() // 回滚
		return nil, errors.New("db_add_fail")
	}
	c.Id, err = db2.Result().LastInsertId()
	if c.Id <= 0 || err != nil {
		db2.tx.Rollback() // 回滚
		return nil, errors.New("db_add_fail")
	}
	err = db2.Commit()
	if err != nil {
		return nil, errors.New("db_add_fail")
	}
	c.State = cs
	return c, nil
}

// UpdateCouple 修改配对信息
func UpdateCouple(c *entity.Couple) (*entity.Couple, error) {
	c.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_COUPLE).
		Set("update_at=?,together_at=?,creator_name=?,invitee_name=?,creator_avatar=?,invitee_avatar=?").
		Where("id=?").
		Exec(c.UpdateAt, c.TogetherAt, c.CreatorName, c.InviteeName, c.CreatorAvatar, c.InviteeAvatar, c.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return c, nil
}

// GetCoupleById
func GetCoupleById(cid int64) (*entity.Couple, error) {
	var c entity.Couple
	var cs entity.CoupleState
	db := mysqlDB().
		Select(TABLE_COUPLE + ".id," + TABLE_COUPLE + ".status," + TABLE_COUPLE + ".create_at," + TABLE_COUPLE + ".update_at," + TABLE_COUPLE + ".together_at," + TABLE_COUPLE + ".creator_id," + TABLE_COUPLE + ".invitee_id," + TABLE_COUPLE + ".creator_name," + TABLE_COUPLE + ".invitee_name," + TABLE_COUPLE + ".creator_avatar," + TABLE_COUPLE + ".invitee_avatar," + TABLE_COUPLE_STATE + ".id," + TABLE_COUPLE_STATE + ".status," + TABLE_COUPLE_STATE + ".create_at," + TABLE_COUPLE_STATE + ".user_id," + TABLE_COUPLE_STATE + ".state").
		Form(TABLE_COUPLE).
		LeftJoin("(SELECT * FROM " + TABLE_COUPLE_STATE + " WHERE status>=? AND couple_id=? ORDER BY create_at DESC LIMIT 0,1) AS " + TABLE_COUPLE_STATE).
		On(TABLE_COUPLE + ".id=" + TABLE_COUPLE_STATE + ".couple_id").
		Where(TABLE_COUPLE + ".status>=? AND " + TABLE_COUPLE + ".id=?").
		Query(entity.STATUS_VISIBLE, cid, entity.STATUS_VISIBLE, cid).
		NextScan(&c.Id, &c.Status, &c.CreateAt, &c.UpdateAt, &c.TogetherAt, &c.CreatorId, &c.InviteeId, &c.CreatorName, &c.InviteeName, &c.CreatorAvatar, &c.InviteeAvatar, &cs.Id, &cs.Status, &cs.CreateAt, &cs.UserId, &cs.State)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	} else if cs.Id <= 0 {
		return nil, errors.New("couple_state_error")
	}
	c.State = &cs
	return &c, nil
}

// GetCoupleByUser
func GetCoupleByUser(uid int64) (*entity.Couple, error) {
	var c entity.Couple
	var cs entity.CoupleState
	db := mysqlDB().
		Select(TABLE_COUPLE + ".id," + TABLE_COUPLE + ".status," + TABLE_COUPLE + ".create_at," + TABLE_COUPLE + ".update_at," + TABLE_COUPLE + ".together_at," + TABLE_COUPLE + ".creator_id," + TABLE_COUPLE + ".invitee_id," + TABLE_COUPLE + ".creator_name," + TABLE_COUPLE + ".invitee_name," + TABLE_COUPLE + ".creator_avatar," + TABLE_COUPLE + ".invitee_avatar," + TABLE_COUPLE_STATE + ".id," + TABLE_COUPLE_STATE + ".status," + TABLE_COUPLE_STATE + ".create_at," + TABLE_COUPLE_STATE + ".user_id," + TABLE_COUPLE_STATE + ".state").
		Form(TABLE_COUPLE).
		LeftJoin("(SELECT * FROM " + TABLE_COUPLE_STATE + " WHERE status>=?) AS " + TABLE_COUPLE_STATE).
		On(TABLE_COUPLE + ".id=" + TABLE_COUPLE_STATE + ".couple_id").
		Where(TABLE_COUPLE + ".status>=? AND (" + TABLE_COUPLE + ".creator_id=? OR " + TABLE_COUPLE + ".invitee_id=?)").
		OrderDown(TABLE_COUPLE_STATE + ".create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, entity.STATUS_VISIBLE, uid, uid).
		NextScan(&c.Id, &c.Status, &c.CreateAt, &c.UpdateAt, &c.TogetherAt, &c.CreatorId, &c.InviteeId, &c.CreatorName, &c.InviteeName, &c.CreatorAvatar, &c.InviteeAvatar, &cs.Id, &cs.Status, &cs.CreateAt, &cs.UserId, &cs.State)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	} else if cs.Id <= 0 {
		return nil, errors.New("couple_state_error")
	}
	c.State = &cs
	return &c, nil
}

// GetCoupleBy2User
func GetCoupleBy2User(uid1, uid2 int64) (*entity.Couple, error) {
	var c entity.Couple
	var cs entity.CoupleState
	db := mysqlDB().
		Select(TABLE_COUPLE + ".id," + TABLE_COUPLE + ".status," + TABLE_COUPLE + ".create_at," + TABLE_COUPLE + ".update_at," + TABLE_COUPLE + ".together_at," + TABLE_COUPLE + ".creator_id," + TABLE_COUPLE + ".invitee_id," + TABLE_COUPLE + ".creator_name," + TABLE_COUPLE + ".invitee_name," + TABLE_COUPLE + ".creator_avatar," + TABLE_COUPLE + ".invitee_avatar," + TABLE_COUPLE_STATE + ".id," + TABLE_COUPLE_STATE + ".status," + TABLE_COUPLE_STATE + ".create_at," + TABLE_COUPLE_STATE + ".user_id," + TABLE_COUPLE_STATE + ".state").
		Form(TABLE_COUPLE).
		LeftJoin("(SELECT * FROM " + TABLE_COUPLE_STATE + " WHERE status>=?) AS " + TABLE_COUPLE_STATE).
		On(TABLE_COUPLE + ".id=" + TABLE_COUPLE_STATE + ".couple_id").
		Where(TABLE_COUPLE + ".status>=? AND ((" + TABLE_COUPLE + ".creator_id=? AND " + TABLE_COUPLE + ".invitee_id=?) OR (" + TABLE_COUPLE + ".invitee_id=? AND " + TABLE_COUPLE + ".creator_id=?))").
		OrderDown(TABLE_COUPLE_STATE + ".create_at").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, entity.STATUS_VISIBLE, uid1, uid2, uid1, uid2).
		NextScan(&c.Id, &c.Status, &c.CreateAt, &c.UpdateAt, &c.TogetherAt, &c.CreatorId, &c.InviteeId, &c.CreatorName, &c.InviteeName, &c.CreatorAvatar, &c.InviteeAvatar, &cs.Id, &cs.Status, &cs.CreateAt, &cs.UserId, &cs.State)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if c.Id <= 0 {
		return nil, nil
	} else if cs.Id <= 0 {
		return nil, errors.New("couple_state_error")
	}
	c.State = &cs
	return &c, nil
}

/****************************************** admin ***************************************/

// GetCoupleList
func GetCoupleList(uid int64, offset, limit int) ([]*entity.Couple, error) {
	hasUser := uid > 0
	where := "status>=?"
	if hasUser {
		where = where + " AND (creator_id=? OR invitee_id=?)"
	}
	db := mysqlDB().
		Select("id,status,create_at,update_at,together_at,creator_id,invitee_id,creator_name,invitee_name,creator_avatar,invitee_avatar").
		Form(TABLE_COUPLE).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		db.Query(entity.STATUS_VISIBLE)
	} else {
		db.Query(entity.STATUS_VISIBLE, uid, uid)
	}
	defer db.Close()
	list := make([]*entity.Couple, 0)
	for db.Next() {
		var c entity.Couple
		db.Scan(&c.Id, &c.Status, &c.CreateAt, &c.UpdateAt, &c.TogetherAt, &c.CreatorId, &c.InviteeId, &c.CreatorName, &c.InviteeName, &c.CreatorAvatar, &c.InviteeAvatar)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &c)
	}
	return list, nil
}

// GetCoupleTotalByCreateWithDel
func GetCoupleTotalByCreateWithDel(start, end int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_COUPLE).
		Where("create_at BETWEEN ? AND ?").
		Query(start, end).
		NextScan(&total)
	defer db.Close()
	return total
}
