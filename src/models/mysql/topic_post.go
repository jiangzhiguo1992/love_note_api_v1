package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddPost
func AddPost(p *entity.Post) (*entity.Post, error) {
	p.Status = entity.STATUS_VISIBLE
	p.CreateAt = time.Now().Unix()
	p.UpdateAt = time.Now().Unix()
	p.ContentImages = entity.JoinStrByColon(p.ContentImageList)
	p.Top = false
	p.Well = false
	p.ReportCount = 0
	p.PointCount = 0
	p.CollectCount = 0
	p.CommentCount = 0
	db := mysqlDB().
		Insert(TABLE_POST).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,kind=?,sub_kind=?,title=?,content_text=?,content_images=?,top=?,official=?,well=?,report_count=?,point_count=?,collect_count=?,comment_count=?").
		Exec(p.Status, p.CreateAt, p.UpdateAt, p.UserId, p.CoupleId, p.Kind, p.SubKind, p.Title, p.ContentText, p.ContentImages, p.Top, p.Official, p.Well, p.ReportCount, p.PointCount, p.CollectCount, p.CommentCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	p.Id, _ = db.Result().LastInsertId()
	p.ContentImages = ""
	return p, nil
}

// DelPost
func DelPost(p *entity.Post) error {
	p.Status = entity.STATUS_DELETE
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(p.Status, p.UpdateAt, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

// UpdatePost
func UpdatePost(p *entity.Post) (*entity.Post, error) {
	p.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_POST).
		Set("update_at=?,top=?,official=?,well=?").
		Where("id=?").
		Exec(p.UpdateAt, p.Top, p.Official, p.Well, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return p, nil
}

// UpdatePostCount
func UpdatePostCount(p *entity.Post, update bool) (*entity.Post, error) {
	if p.Official || p.ReportCount < 0 {
		p.ReportCount = 0
	}
	if p.PointCount < 0 {
		p.PointCount = 0
	}
	if p.CollectCount < 0 {
		p.CollectCount = 0
	}
	if p.CommentCount < 0 {
		p.CommentCount = 0
	}
	if update {
		p.UpdateAt = time.Now().Unix()
	}
	db := mysqlDB().
		Update(TABLE_POST).
		Set("update_at=?,report_count=?,point_count=?,collect_count=?,comment_count=?").
		Where("id=?").
		Exec(p.UpdateAt, p.ReportCount, p.PointCount, p.CollectCount, p.CommentCount, p.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	p.ContentImages = ""
	return p, nil
}

// GetPostById
func GetPostById(pid int64) (*entity.Post, error) {
	var p entity.Post
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where("id=?").
		Query(pid).
		NextScan(&p.Id, &p.Status, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if p.Id <= 0 {
		return nil, nil
	} else if p.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
	p.ContentImages = ""
	return &p, nil
}

// GetPostListBySearch
func GetPostListBySearch(search string, reportLimit, offset, limit int) ([]*entity.Post, error) {
	search = string('%') + search + string('%')
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where("status>=? AND title LIKE ? AND report_count<?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, search, reportLimit)
	defer db.Close()
	list := make([]*entity.Post, 0)
	for db.Next() {
		var p entity.Post
		p.Screen = false
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
		p.ContentImages = ""
		list = append(list, &p)
	}
	return list, nil
}

// GetPostListByCreate
func GetPostListByCreate(create int64, reportLimit, offset, limit int) ([]*entity.Post, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where("status>=? AND create_at<=? AND report_count<?").
		Order("top DESC,update_at DESC").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, create, reportLimit)
	defer db.Close()
	list := make([]*entity.Post, 0)
	for db.Next() {
		var p entity.Post
		p.Screen = false
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
		p.ContentImages = ""
		list = append(list, &p)
	}
	return list, nil
}

// GetPostListByCreateKindOfficialWell
func GetPostListByCreateKindOfficialWell(create int64, kind, subKind int, official, well bool, reportLimit, offset, limit int) ([]*entity.Post, error) {
	// 条件判断
	hasSubKind := subKind > 0
	hasOfficial := official
	hasWell := well
	// 构造where和args
	where := "status>=? AND create_at<=? AND report_count<? AND kind=?"
	if hasSubKind {
		where = where + " AND sub_kind=?"
	}
	if hasOfficial {
		where = where + " AND official=?"
	} else if hasWell {
		where = where + " AND well=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where(where).
		Order("top DESC,update_at DESC").
		Limit(offset, limit)
	if hasSubKind {
		if hasOfficial || hasWell {
			db.Query(entity.STATUS_VISIBLE, create, reportLimit, kind, subKind, true)
		} else {
			db.Query(entity.STATUS_VISIBLE, create, reportLimit, kind, subKind)
		}
	} else {
		if hasOfficial || hasWell {
			db.Query(entity.STATUS_VISIBLE, create, reportLimit, kind, true)
		} else {
			db.Query(entity.STATUS_VISIBLE, create, reportLimit, kind)
		}
	}
	defer db.Close()
	list := make([]*entity.Post, 0)
	for db.Next() {
		var p entity.Post
		p.Screen = false
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
		p.ContentImages = ""
		list = append(list, &p)
	}
	return list, nil
}

// GetPostListByUserCouple
func GetPostListByUserCouple(uid, cid int64, offset, limit int) ([]*entity.Post, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where("status>=? AND user_id=? AND couple_id=? AND kind<>?").
		OrderDown("update_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, uid, cid, entity.POST_KIND_LIMIT_UNKNOWN)
	defer db.Close()
	list := make([]*entity.Post, 0)
	for db.Next() {
		var p entity.Post
		p.UserId = uid
		p.CoupleId = cid
		p.Screen = false
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
		p.ContentImages = ""
		list = append(list, &p)
	}
	return list, nil
}

/****************************************** admin ***************************************/

// GetPostList
func GetPostList(uid int64, offset, limit int) ([]*entity.Post, error) {
	where := "status>=?"
	hasUser := uid > 0
	if hasUser {
		where = where + " AND user_id=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,sub_kind,title,content_text,content_images,top,official,well,report_count,point_count,collect_count,comment_count").
		Form(TABLE_POST).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasUser {
		db.Query(entity.STATUS_VISIBLE, )
	} else {
		db.Query(entity.STATUS_VISIBLE, uid)
	}
	defer db.Close()
	list := make([]*entity.Post, 0)
	for db.Next() {
		var p entity.Post
		db.Scan(&p.Id, &p.CreateAt, &p.UpdateAt, &p.UserId, &p.CoupleId, &p.Kind, &p.SubKind, &p.Title, &p.ContentText, &p.ContentImages, &p.Top, &p.Official, &p.Well, &p.ReportCount, &p.PointCount, &p.CollectCount, &p.CommentCount)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		p.ContentImageList = entity.SplitStrByColon(p.ContentImages)
		p.ContentImages = ""
		list = append(list, &p)
	}
	return list, nil
}

// GetPostTotalByCreateWithDel
func GetPostTotalByCreateWithDel(create int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_POST).
		Where("create_at>?").
		Query(create).
		NextScan(&total)
	defer db.Close()
	return total
}
