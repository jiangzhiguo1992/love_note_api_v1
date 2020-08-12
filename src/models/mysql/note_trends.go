package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTrends
func AddTrends(t *entity.Trends) (*entity.Trends, error) {
	t.Status = entity.STATUS_VISIBLE
	t.CreateAt = time.Now().Unix()
	t.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TRENDS).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,action_type=?,content_type=?,content_id=?").
		Exec(t.Status, t.CreateAt, t.UpdateAt, t.UserId, t.CoupleId, t.ActionType, t.ContentType, t.ContentId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	t.Id, _ = db.Result().LastInsertId()
	return t, nil
}

// UpdateTrends
func UpdateTrends(t *entity.Trends) (*entity.Trends, error) {
	t.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TRENDS).
		Set("update_at=?").
		Where("id=?").
		Exec(t.UpdateAt, t.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return t, nil
}

// GetTrendsByUserCoupleTypeId
func GetTrendsByUserCoupleTypeId(uid, cid int64, atp, ctp int, conId int64) (*entity.Trends, error) {
	var t entity.Trends
	t.UserId = uid
	t.CoupleId = cid
	t.ActionType = atp
	t.ContentType = ctp
	t.ContentId = conId
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_TRENDS).
		Where("status>=? AND user_id=? AND couple_id=? AND action_type=? AND content_type=? AND content_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid, atp, ctp, conId).
		NextScan(&t.Id, &t.CreateAt, &t.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if t.Id <= 0 {
		return nil, nil
	}
	return &t, nil
}

// GetTrendsListByCreateCouple
func GetTrendsListByCreateCouple(createAt, cid int64, offset, limit int) ([]*entity.Trends, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,action_type,content_type,content_id").
		Form(TABLE_TRENDS).
		Where("status>=? AND update_at<=? AND couple_id=?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, createAt, cid)
	defer db.Close()
	list := make([]*entity.Trends, 0)
	for db.Next() {
		var t entity.Trends
		t.CoupleId = cid
		db.Scan(&t.Id, &t.CreateAt, &t.UpdateAt, &t.UserId, &t.ActionType, &t.ContentType, &t.ContentId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &t)
	}
	return list, nil
}

// GetTrendsTotalByUpdateUserCouple
func GetTrendsTotalByUpdateUserCouple(update, uid, cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TRENDS).
		Where("status>=? AND update_at>=? AND couple_id=? AND user_id=?").
		Query(entity.STATUS_VISIBLE, update, cid, uid).
		NextScan(&total)
	defer db.Close()
	return total
}

/****************************************** admin ***************************************/

// GetTrendsList
func GetTrendsList(uid, cid int64, actType, conType, offset, limit int) ([]*entity.Trends, error) {
	where := "status>=?"
	hasUser := uid > 0
	hasCouple := cid > 0
	hasAct := actType > 0
	hasCon := conType > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	if hasCouple {
		where = where + " AND couple_id=?"
	}
	if hasAct {
		where = where + " AND action_type=?"
	}
	if hasCon {
		where = where + " AND content_type=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,action_type,content_type,content_id").
		Form(TABLE_TRENDS).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		if !hasCouple {
			if !hasAct {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE)
				} else {
					db.Query(entity.STATUS_VISIBLE, conType)
				}
			} else {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, actType)
				} else {
					db.Query(entity.STATUS_VISIBLE, actType, conType)
				}
			}
		} else {
			if !hasAct {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, cid)
				} else {
					db.Query(entity.STATUS_VISIBLE, cid, conType)
				}
			} else {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, cid, actType)
				} else {
					db.Query(entity.STATUS_VISIBLE, cid, actType, conType)
				}
			}
		}
	} else {
		if !hasCouple {
			if !hasAct {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, uid)
				} else {
					db.Query(entity.STATUS_VISIBLE, uid, conType)
				}
			} else {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, uid, actType)
				} else {
					db.Query(entity.STATUS_VISIBLE, uid, actType, conType)
				}
			}
		} else {
			if !hasAct {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, uid, cid)
				} else {
					db.Query(entity.STATUS_VISIBLE, uid, cid, conType)
				}
			} else {
				if !hasCon {
					db.Query(entity.STATUS_VISIBLE, uid, cid, actType)
				} else {
					db.Query(entity.STATUS_VISIBLE, uid, cid, actType, conType)
				}
			}
		}
	}
	defer db.Close()
	list := make([]*entity.Trends, 0)
	for db.Next() {
		var t entity.Trends
		db.Scan(&t.Id, &t.CreateAt, &t.UpdateAt, &t.UserId, &t.CoupleId, &t.ActionType, &t.ContentType, &t.ContentId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &t)
	}
	return list, nil
}

// GetTrendsContentTypeList
func GetTrendsContentTypeList(start, end int64) ([]*entity.FiledInfo, error) {
	db := mysqlDB().
		Select("content_type,COUNT(content_type) AS nums").
		Form(TABLE_TRENDS).
		Where("status>=? AND (create_at BETWEEN ? AND ?)").
		Group("content_type").
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
