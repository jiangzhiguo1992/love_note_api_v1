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

// AddPost
func AddPost(uid, cid int64, p *entity.Post) (*entity.Post, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if p == nil {
		return nil, errors.New("nil_post")
	} else if !GetPostSubKindEnable(p.Kind, p.SubKind) {
		return nil, errors.New("request_model_close")
	}
	if len(strings.TrimSpace(p.Title)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(p.Title)) > GetLimit().PostTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len(strings.TrimSpace(p.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(p.ContentText)) > GetLimit().PostContentLength {
		return nil, errors.New("limit_content_text_over")
	} else if len(p.ContentImageList) > 0 {
		imgLimit := GetVipLimitByCouple(cid).TopicPostImageCount
		if imgLimit <= 0 {
			return nil, errors.New("limit_content_image_refuse")
		} else if len(p.ContentImageList) > imgLimit {
			return nil, errors.New("limit_content_image_over")
		}
	}
	// admin
	u, _ := GetUserById(uid)
	p.Official = IsAdminister(u)
	// mysql
	p.UserId = uid
	p.CoupleId = cid
	p, err := mysql.AddPost(p)
	if p == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		TopicInfoUpdatePost(p.Kind, true)
	}()
	// redis-set
	redis.SetPost(p)
	// TODO redis
	//redis.AddPostInListByAll(p)
	return p, err
}

// DelPost
func DelPost(uid, cid, pid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if pid <= 0 {
		return errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pid)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("nil_post")
	}
	// admin
	u, _ := GetUserById(uid)
	admin := IsAdminister(u)
	if !admin && p.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelPost(p)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		// message
		if admin {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(p.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_official_del_post") + p.Title
			message := CreateTopicMessage(uid, cid, p.UserId, p.CoupleId, entity.TOPIC_MESSAGE_KIND_OFFICIAL_TEXT, content, p.Id)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpdatePost(p.Kind, false)
	}()
	// redis-del
	redis.DelPost(p)
	// TODO redis
	//redis.DelPostInListByAll(p)
	//if p.Well {
	//	redis.DelPostInListByWell(p)
	//}
	//if p.Official {
	//	redis.DelPostInListByOfficial(p)
	//}
	return err
}

// UpdatePost
func UpdatePost(p *entity.Post) (*entity.Post, error) {
	if p == nil || p.Id <= 0 {
		return nil, errors.New("nil_post")
	}
	// redis_del
	redis.DelPost(p)
	// mysql
	p, err := mysql.UpdatePost(p)
	// redis-set
	redis.SetPost(p)
	// TODO redis
	//redis.UpdatePostInListByAll(p)
	//if p.Well {
	//	redis.UpdatePostInListByWell(p)
	//}
	//if p.Official {
	//	redis.UpdatePostInListByOfficial(p)
	//}
	return p, err
}

// UpdatePostCount
func UpdatePostCount(p *entity.Post, update bool) (*entity.Post, error) {
	if p == nil || p.Id <= 0 {
		return nil, errors.New("nil_post")
	} else if !GetPostSubKindEnable(p.Kind, p.SubKind) {
		return nil, errors.New("request_model_close")
	}
	// redis_del
	redis.DelPost(p)
	// mysql
	p, err := mysql.UpdatePostCount(p, update)
	// redis-set
	redis.SetPost(p)
	// TODO redis
	//redis.UpdatePostInListByAll(p)
	//if p.Well {
	//	redis.UpdatePostInListByWell(p)
	//}
	//if p.Official {
	//	redis.UpdatePostInListByOfficial(p)
	//}
	return p, err
}

// GetPostById
func GetPostById(pid int64) (*entity.Post, error) {
	// redis-get
	p, err := redis.GetPostById(pid)
	if p == nil || p.Id <= 0 || err != nil {
		p, err = mysql.GetPostById(pid)
	}
	// mysql
	if p == nil || err != nil {
		return p, err
	}
	// enable
	if !GetPostSubKindEnable(p.Kind, p.SubKind) {
		return nil, errors.New("request_model_close")
	}
	// redis-set
	redis.SetPost(p)
	return p, err
}

// GetPostByIdWithAll
func GetPostByIdWithAll(uid, cid, pid int64) (*entity.Post, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if pid <= 0 {
		return nil, errors.New("nil_post")
	}
	// mysql
	p, err := GetPostById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// 额外属性
	LoadPostWithAll(uid, cid, p)
	return p, err
}

// GetPostListByCreate
func GetPostListByCreate(uid, cid, create int64, page int) ([]*entity.Post, error) {
	if create <= 0 {
		return nil, nil
	} else if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	var list []*entity.Post
	var err error
	limit := GetPageSizeLimit().Post
	offset := page * limit
	// TODO redis
	//if official {
	//	list, _ = redis.GetPostListByOfficial(kind, subKind, offset, limit)
	//} else if well {
	//	list, _ = redis.GetPostListByWell(kind, subKind, offset, limit)
	//} else {
	//	list, _ = redis.GetPostListByAll(kind, subKind, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		reportLimit := GetLimit().PostScreenReportCount
		list, err = mysql.GetPostListByCreate(create, reportLimit, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if official {
		//		redis.SetPostListByOfficial(kind, subKind, list, true)
		//	} else if well {
		//		redis.SetPostListByWell(kind, subKind, list, true)
		//	} else {
		//		redis.SetPostListByAll(kind, subKind, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		//LoadPostWithAll(uid, cid, v)
		LoadPostWithAll(uid, 0, v)
	}
	return list, err
}

// GetPostListBySearch
func GetPostListBySearch(uid, cid int64, search string, page int) ([]*entity.Post, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if len(strings.TrimSpace(search)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	var list []*entity.Post
	var err error
	limit := GetPageSizeLimit().Post
	offset := page * limit
	// TODO redis
	//if official {
	//	list, _ = redis.GetPostListByOfficial(kind, subKind, offset, limit)
	//} else if well {
	//	list, _ = redis.GetPostListByWell(kind, subKind, offset, limit)
	//} else {
	//	list, _ = redis.GetPostListByAll(kind, subKind, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		reportLimit := GetLimit().PostScreenReportCount
		list, err = mysql.GetPostListBySearch(search, reportLimit, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if official {
		//		redis.SetPostListByOfficial(kind, subKind, list, true)
		//	} else if well {
		//		redis.SetPostListByWell(kind, subKind, list, true)
		//	} else {
		//		redis.SetPostListByAll(kind, subKind, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		//LoadPostWithAll(uid, cid, v)
		LoadPostWithAll(uid, 0, v)
	}
	return list, err
}

// GetPostListByCreateKindOfficialWell
func GetPostListByCreateKindOfficialWell(uid, cid, create int64, kind, subKind int, official, well bool, page int) ([]*entity.Post, error) {
	if create <= 0 {
		return nil, nil
	} else if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if !GetPostSubKindEnable(kind, subKind) {
		return nil, errors.New("request_model_close")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	var list []*entity.Post
	var err error
	limit := GetPageSizeLimit().Post
	offset := page * limit
	// TODO redis
	//if official {
	//	list, _ = redis.GetPostListByOfficial(kind, subKind, offset, limit)
	//} else if well {
	//	list, _ = redis.GetPostListByWell(kind, subKind, offset, limit)
	//} else {
	//	list, _ = redis.GetPostListByAll(kind, subKind, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		reportLimit := GetLimit().PostScreenReportCount
		list, err = mysql.GetPostListByCreateKindOfficialWell(create, kind, subKind, official, well, reportLimit, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if official {
		//		redis.SetPostListByOfficial(kind, subKind, list, true)
		//	} else if well {
		//		redis.SetPostListByWell(kind, subKind, list, true)
		//	} else {
		//		redis.SetPostListByAll(kind, subKind, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		//LoadPostWithAll(uid, cid, v)
		LoadPostWithAll(uid, 0, v)
	}
	return list, err
}

// GetPostListByUserCouple
func GetPostListByUserCouple(uid, cid int64, page int) ([]*entity.Post, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Post
	offset := page * limit
	list, err := mysql.GetPostListByUserCouple(uid, cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		//LoadPostWithAll(uid, cid, v)
		LoadPostWithAll(uid, 0, v)
	}
	return list, err
}

// GetPostList
func GetPostList(uid int64, page int) ([]*entity.Post, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Post
	offset := page * limit
	list, err := mysql.GetPostList(uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetPostListByUserCoupleCollect
func GetPostListByUserCoupleCollect(meid, suid, cid int64, page int) ([]*entity.Post, error) {
	if meid <= 0 || suid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// collectList
	limit := GetPageSizeLimit().Post
	offset := page * limit
	collectList, err := mysql.GetPostCollectListByUserCouple(suid, cid, offset, limit)
	if err != nil {
		return nil, err
	} else if collectList == nil || len(collectList) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// postList
	list := make([]*entity.Post, 0)
	for _, v := range collectList {
		if v == nil || v.PostId <= 0 {
			continue
		}
		p, err := GetPostById(v.PostId)
		if p == nil || p.Id <= 0 || p.Status <= entity.STATUS_DELETE || err != nil {
			// 删除或异常的帖子
			p = &entity.Post{}
			p.Id = v.PostId
			p.Status = entity.STATUS_DELETE
		} else {
			// 额外属性
			//LoadPostWithAll(uid, cid, p)
			LoadPostWithAll(meid, 0, p)
		}
		list = append(list, p)
	}
	return list, nil
}

// GetPostReportList
func GetPostReportList(page int) ([]*entity.Post, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// reportList
	limit := GetPageSizeLimit().Post
	offset := page * limit
	reportList, err := mysql.GetPostReportList(offset, limit)
	if err != nil {
		return nil, err
	} else if reportList == nil || len(reportList) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_post")
		} else {
			return nil, nil
		}
	}
	// postList
	list := make([]*entity.Post, 0)
	for _, v := range reportList {
		if v == nil || v.PostId <= 0 {
			continue
		}
		p, _ := GetPostById(v.PostId)
		if p != nil && p.Id > 0 {
			list = append(list, p)
		}
	}
	return list, nil
}

// GetPostTotalByCreateWithDel
func GetPostTotalByCreateWithDel(create int64) int64 {
	if create == 0 {
		return 0
	}
	// mysql
	total := mysql.GetPostTotalByCreateWithDel(create)
	return total
}

// LoadPostWithAll
func LoadPostWithAll(uid, cid int64, p *entity.Post) (*entity.Post, error) {
	if p == nil || p.Id <= 0 {
		return nil, nil
	}
	// 额外属性
	if p.Official {
		p.Screen = false
	} else {
		if p.ReportCount < GetLimit().PostScreenReportCount {
			p.Screen = false
		} else {
			p.Screen = true
		}
	}
	p.Hot = IsPostHot(p)
	if p.Kind == entity.POST_KIND_LIMIT_UNKNOWN {
		// 匿名
		p.Couple = nil
	} else {
		p.Couple, _ = GetCoupleVisibleByUser(p.UserId)
	}
	p.Read = IsPostReadByUserCouple(uid, p.Id)
	if cid <= 0 {
		// 没配对
		p.Mine = false
		p.Our = false
		p.Report = false
		p.Point = false
		p.Collect = false
		p.Comment = false
	} else {
		if p.Kind == entity.POST_KIND_LIMIT_UNKNOWN {
			// 匿名
			p.Mine = false
			p.Our = false
			p.Comment = false
		} else {
			p.Mine = p.UserId == uid
			p.Our = p.CoupleId == cid
			p.Comment = IsPostCommentByUserCouple(uid, cid, p.Id)
		}
		p.Report = IsPostReportByUserCouple(uid, cid, p.Id)
		p.Point = IsPostPointByUserCouple(uid, cid, p.Id)
		p.Collect = IsPostCollectByUserCouple(uid, cid, p.Id)
	}
	return p, nil
}

// GetPostKindEnable
func GetPostKindEnable(kind int) bool {
	if entity.PostKindList == nil || len(entity.PostKindList) <= 0 {
		return false
	}
	for _, v := range entity.PostKindList {
		if v == kind {
			return true
		}
	}
	return false
}

// GetPostSubKindEnable
func GetPostSubKindEnable(kind, subKind int) bool {
	kindEnable := GetPostKindEnable(kind)
	if !kindEnable {
		return false
	}
	subKindList := entity.GetPostSubKindMap()[kind]
	if subKindList == nil || len(subKindList) <= 0 {
		return false
	}
	for _, v := range subKindList {
		if v == subKind {
			return true
		}
	}
	return false
}

// IsPostHot
func IsPostHot(p *entity.Post) bool {
	if p == nil || p.Id <= 0 {
		return false
	}
	existHour := (time.Now().Unix() - p.CreateAt) / (60 * 60)
	postHotMinCreateHour := utils.GetConfigInt64("conf", "limit.conf", "time", "topic_post_hot_min_create_hour")
	if existHour < postHotMinCreateHour {
		// 发帖时间小于一小时
		return false
	}
	postHotPointPerHour := utils.GetConfigInt64("conf", "limit.conf", "time", "topic_post_hot_point_per_hour")
	postHotCollectPerHour := utils.GetConfigInt64("conf", "limit.conf", "time", "topic_post_hot_collect_per_hour")
	postHotCommentPerHour := utils.GetConfigInt64("conf", "limit.conf", "time", "topic_post_hot_comment_per_hour")
	if int64(p.PointCount) < existHour*postHotPointPerHour {
		return false
	} else if int64(p.CollectCount) < existHour*postHotCollectPerHour {
		return false
	} else if int64(p.CommentCount) < existHour*postHotCommentPerHour {
		return false
	}
	return true
}
