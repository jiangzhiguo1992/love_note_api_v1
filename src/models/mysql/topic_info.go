package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTopicInfo
func AddTopicInfo(ti *entity.TopicInfo) (*entity.TopicInfo, error) {
	ti.Status = entity.STATUS_VISIBLE
	ti.CreateAt = time.Now().Unix()
	ti.UpdateAt = time.Now().Unix()
	ti.PostCount = 0
	ti.BrowseCount = 0
	ti.CommentCount = 0
	ti.ReportCount = 0
	ti.PointCount = 0
	ti.CollectCount = 0
	db := mysqlDB().
		Insert(TABLE_TOPIC_INFO).
		Set("status=?,create_at=?,update_at=?,kind=?,year=?,day_of_year=?,post_count=?,browse_count=?,comment_count=?,report_count=?,point_count=?,collect_count=?").
		Exec(ti.Status, ti.CreateAt, ti.UpdateAt, ti.Kind, ti.Year, ti.DayOfYear, ti.PostCount, ti.BrowseCount, ti.CommentCount, ti.ReportCount, ti.PointCount, ti.CollectCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	ti.Id, _ = db.Result().LastInsertId()
	return ti, nil
}

// UpdateTopicInfo
func UpdateTopicInfo(ti *entity.TopicInfo) (*entity.TopicInfo, error) {
	if ti.BrowseCount < 0 {
		ti.BrowseCount = 0
	}
	if ti.PostCount < 0 {
		ti.PostCount = 0
	}
	if ti.CommentCount < 0 {
		ti.CommentCount = 0
	}
	if ti.ReportCount < 0 {
		ti.ReportCount = 0
	}
	if ti.PointCount < 0 {
		ti.PointCount = 0
	}
	if ti.CollectCount < 0 {
		ti.CollectCount = 0
	}
	ti.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TOPIC_INFO).
		Set("update_at=?,post_count=?,browse_count=?,comment_count=?,report_count=?,point_count=?,collect_count=?").
		Where("id=?").
		Exec(ti.UpdateAt, ti.PostCount, ti.BrowseCount, ti.CommentCount, ti.ReportCount, ti.PointCount, ti.CollectCount, ti.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return ti, nil
}

// GetTopicInfoByKindYearDays
func GetTopicInfoByKindYearDays(kind, year, days int) (*entity.TopicInfo, error) {
	var ti entity.TopicInfo
	ti.Kind = kind
	ti.Year = year
	ti.DayOfYear = days
	db := mysqlDB().
		Select("id,create_at,update_at,post_count,browse_count,comment_count,report_count,point_count,collect_count").
		Form(TABLE_TOPIC_INFO).
		Where("status >=? AND kind=? AND year=? AND day_of_year=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, kind, year, days).
		NextScan(&ti.Id, &ti.CreateAt, &ti.UpdateAt, &ti.PostCount, &ti.BrowseCount, &ti.CommentCount, &ti.ReportCount, &ti.PointCount, &ti.CollectCount)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if ti.Id <= 0 {
		return nil, nil
	}
	return &ti, nil
}
