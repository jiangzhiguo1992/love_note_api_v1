package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"models/redis"
	"strings"
	"time"
)

// CoupleNoBreakTime
// 1.create大于这个时间的都还没完全分手
func CoupleNoBreakTime() int64 {
	return time.Now().Unix() - GetLimit().CoupleBreakSec
}

// AddCouple 添加配对
// 1.状态从邀请开始
// 2.曾经没有过邀请
func AddCouple(c *entity.Couple) (*entity.Couple, error) {
	if c == nil {
		return nil, errors.New("nil_couple")
	} else if c.CreatorId == c.InviteeId {
		return nil, errors.New("couple_creator_same_self")
	} else if c.InviteeId <= 0 {
		return nil, errors.New("couple_invitee_id_nil")
	} else if c.CreatorId == c.InviteeId {
		return nil, errors.New("couple_creator_same_self")
	} else if c.InviteeId <= 0 {
		return nil, errors.New("couple_invitee_id_nil")
	}
	// 验证双方可见的cp(两个人之外，其他人员的关系检查)
	if cp, _ := GetCoupleSelfByUser(c.CreatorId); cp != nil {
		return nil, errors.New("couple_invite_twice")
	}
	if cp, _ := GetCoupleSelfByUser(c.InviteeId); cp != nil {
		return nil, errors.New("couple_invite_late")
	}
	// mysql 事务操作不用抽出来
	c, err := mysql.AddCouple(c)
	return c, err
}

// AddCoupleStateInvitee 添加配对邀请
// 1.曾经有过邀请
// 2.状态从邀请开始
func AddCoupleStateInvitee(uid int64, c *entity.Couple) (*entity.Couple, error) {
	if c == nil || c.Id <= 0 {
		return c, errors.New("nil_couple")
	}
	// 验证双方可见的cp(两个人之外，其他人员的关系检查)
	if cp, _ := GetCoupleSelfByUser(c.CreatorId); cp != nil {
		return c, errors.New("couple_invite_twice")
	}
	if cp, _ := GetCoupleSelfByUser(c.InviteeId); cp != nil {
		return c, errors.New("couple_invite_late")
	}
	cs := c.State
	if cs != nil && cs.State == entity.COUPLE_STATE_INVITE {
		// 重复邀请
		return c, errors.New("couple_invite_repeat")
	} else if cs != nil && (cs.State == entity.COUPLE_STATE_INVITE_CANCEL ||
		cs.State == entity.COUPLE_STATE_INVITE_REJECT ||
		cs.State == entity.COUPLE_STATE_BREAK_ACCEPT) {
		// 检查邀请间隔时间
		inviteIntervalSec := GetLimit().CoupleInviteIntervalSec
		goTime := time.Now().Unix() - cs.CreateAt
		if goTime < int64(inviteIntervalSec) {
			return c, errors.New("couple_invite_frequent")
		}
	} else if cs != nil && cs.State == entity.COUPLE_STATE_BREAK {
		// 检查是否已分手成功
		noBreakTime := CoupleNoBreakTime()
		if cs.CreateAt >= noBreakTime {
			return c, errors.New("couple_status_together_always")
		}
	} else if cs != nil && cs.State == entity.COUPLE_STATE_520 {
		// 已经在一起了
		return c, errors.New("couple_status_together_always")
	} else {
		// 其他异常情况清零
	}
	// data
	cs.State = entity.COUPLE_STATE_INVITE
	cs.UserId = uid
	cs.CoupleId = c.Id
	// mysql
	cs, err := AddCoupleState(cs)
	if cs != nil && err == nil {
		c.State = cs
	}
	return c, err
}

// AddCoupleState 添加新的couple的state
func AddCoupleState(cs *entity.CoupleState) (*entity.CoupleState, error) {
	if cs == nil {
		return nil, errors.New("couple_state_error")
	} else if cs.UserId <= 0 {
		return nil, errors.New("nil_user")
	} else if cs.CoupleId <= 0 {
		return nil, errors.New("nil_couple")
	} else if cs.State != entity.COUPLE_STATE_INVITE &&
		cs.State != entity.COUPLE_STATE_INVITE_CANCEL &&
		cs.State != entity.COUPLE_STATE_INVITE_REJECT &&
		cs.State != entity.COUPLE_STATE_BREAK &&
		cs.State != entity.COUPLE_STATE_BREAK_ACCEPT &&
		cs.State != entity.COUPLE_STATE_520 {
		return nil, errors.New("couple_state_error")
	}
	// couple检查
	couple, err := mysql.GetCoupleById(cs.CoupleId)
	if err != nil {
		return nil, err
	} else if couple == nil || couple.Id <= 0 {
		return nil, errors.New("nil_couple")
	}
	// redis-del
	redis.DelCoupleByUser(couple.CreatorId)
	redis.DelCoupleByUser(couple.InviteeId)
	// mysql
	cs, err = mysql.AddCoupleState(cs)
	// 同步
	go func() {
		// push
		toUid := couple.InviteeId
		if toUid == cs.UserId {
			toUid = couple.CreatorId
		}
		entry, err := mysql.GetEntryLatestByUser(toUid)
		if err != nil {
			return
		} else if entry == nil {
			return
		}
		title := utils.GetLanguage(entry.Language, "push_title_new_notice")
		content := utils.GetLanguage(entry.Language, "push_content_couple_state_change")
		push := CreatePush(cs.UserId, toUid, couple.Id, title, content, entity.PUSH_TYPE_APP)
		AddPush(push)
	}()
	return cs, err
}

// UpdateCouple 修改配对信息
func UpdateCouple(c *entity.Couple) (*entity.Couple, error) {
	if c == nil || c.Id <= 0 {
		return c, errors.New("nil_couple")
	} else if c.TogetherAt > time.Now().Unix() {
		return c, errors.New("limit_happen_err")
	}
	c.CreatorName = strings.TrimSpace(c.CreatorName)
	c.InviteeName = strings.TrimSpace(c.InviteeName)
	c.CreatorAvatar = strings.TrimSpace(c.CreatorAvatar)
	c.InviteeAvatar = strings.TrimSpace(c.InviteeAvatar)
	limitName := GetLimit().CoupleNameLength
	cnChina := []rune(c.CreatorName)
	if len(cnChina) > limitName {
		c.CreatorName = string(cnChina[:limitName])
	}
	inChina := []rune(c.InviteeName)
	if len(inChina) > limitName {
		c.InviteeName = string(inChina[:limitName])
	}
	// redis-get
	redis.DelCoupleByUser(c.CreatorId)
	redis.DelCoupleByUser(c.InviteeId)
	// mysql
	c, err := mysql.UpdateCouple(c)
	if c == nil || err != nil {
		return c, err
	}
	// redis-set
	redis.SetCoupleByUser(c.CreatorId, c)
	redis.SetCoupleByUser(c.InviteeId, c)
	return c, err
}

// GetCoupleTogetherDay 获取couple的最近Together
//func GetCoupleTogetherDay(cid int64) int {
//	days := 1
//	if cid <= 0 {
//		return days
//	}
//	// mysql
//	cs, _ := mysql.GetCoupleStateLatestByState(cid, entity.COUPLE_STATE_520)
//	if cs != nil {
//		togetherSec := time.Now().Unix() - cs.CreateAt
//		if togetherSec >= 0 {
//			days = (int(togetherSec) / 60 / 60 / 24) + 1
//		}
//	}
//	return days
//}

// GetCoupleById 返回最新的state的couple
func GetCoupleById(cid int64) (*entity.Couple, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	c, err := mysql.GetCoupleById(cid)
	return c, err
}

// GetCoupleBy2User 返回最新的state的couple
func GetCoupleBy2User(uid1, uid2 int64) (*entity.Couple, error) {
	if uid1 <= 0 || uid2 <= 0 {
		return nil, errors.New("nil_user")
	}
	// mysql
	c, err := mysql.GetCoupleBy2User(uid1, uid2)
	return c, err
}

// GetCoupleSelfByUser 正在邀请 + 配对成功的 + 正在分手
// 1.用于验证(逻辑上，多个cp时，只有一个cp是可self的，也是state最近的那个)
// 2.返回最新的state的couple(多个不同ta的couple，返回最后操作的cp)
func GetCoupleSelfByUser(uid int64) (*entity.Couple, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	}
	// mysql
	couple, err := mysql.GetCoupleByUser(uid)
	if couple == nil || err != nil {
		return couple, err
	}
	// check
	noBreakTime := CoupleNoBreakTime()
	state := couple.State
	if state.State == entity.COUPLE_STATE_INVITE || state.State == entity.COUPLE_STATE_520 {
		// 邀请 + 在一起
		return couple, nil
	} else if state.State == entity.COUPLE_STATE_BREAK && state.CreateAt >= noBreakTime {
		// 正在分手
		return couple, nil
	}
	return nil, nil
}

// GetCoupleVisibleByUser 配对成功 + 正在分手
// 1.初始化信息
func GetCoupleVisibleByUser(uid int64) (*entity.Couple, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	}
	// redis-get
	couple, _ := redis.GetCoupleByUser(uid)
	if couple != nil && couple.Id > 0 {
		return couple, nil
	}
	// mysql
	couple, err := GetCoupleSelfByUser(uid)
	if couple == nil || err != nil {
		return nil, err
	}
	state := couple.State
	if state == nil || state.State == entity.COUPLE_STATE_INVITE {
		// 邀请
		return nil, nil
	}
	// redis-set
	redis.SetCoupleByUser(couple.CreatorId, couple)
	redis.SetCoupleByUser(couple.InviteeId, couple)
	return couple, nil
}

// GetCoupleList
func GetCoupleList(uid int64, page int) ([]*entity.Couple, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Couple
	offset := page * limit
	list, err := mysql.GetCoupleList(uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_common")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetCoupleStateListByCouple
func GetCoupleStateListByCouple(cid int64, page int) ([]*entity.CoupleState, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Couple
	offset := page * limit
	list, err := mysql.GetCoupleStateListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_common")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetCoupleTotalByCreateWithDel
func GetCoupleTotalByCreateWithDel(start, end int64) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetCoupleTotalByCreateWithDel(start, end)
	return total
}

// GetCoupleStateStateListByCreate
func GetCoupleStateStateListByCreate(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	// mysql
	list, err := mysql.GetCoupleStateStateListByCreate(start, end)
	return list, err
}
