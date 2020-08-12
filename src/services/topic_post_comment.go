package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"models/redis"
	"strings"
)

// AddPostComment
func AddPostComment(uid, taId, cid int64, pc *entity.PostComment) (*entity.PostComment, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pc == nil {
		return nil, errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		return nil, errors.New("nil_post")
	} else if pc.Kind != entity.POST_COMMENT_KIND_TEXT && pc.Kind != entity.POST_COMMENT_KIND_JAB {
		return nil, errors.New("limit_kind_nil")
	} else if pc.Kind == entity.POST_COMMENT_KIND_TEXT && len(strings.TrimSpace(pc.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(pc.ContentText)) > GetLimit().PostCommentContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// post检查
	p, err := GetPostById(pc.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// old检查
	if pc.Kind == entity.POST_COMMENT_KIND_JAB {
		old, err := mysql.GetPostCommentByUserCouplePostCommentKind(uid, cid, pc.PostId, pc.ToCommentId, pc.Kind)
		if err != nil {
			return nil, err
		} else if old != nil {
			return nil, errors.New("report_jab")
		}
	}
	// admin
	u, _ := GetUserById(uid)
	pc.Official = IsAdminister(u)
	// floor
	var toComment *entity.PostComment
	var latestComment *entity.PostComment
	if pc.ToCommentId > 0 {
		toComment, err = GetPostCommentById(pc.ToCommentId)
		if err != nil {
			return nil, err
		} else if toComment == nil {
			return nil, errors.New("nil_comment")
		}
		latestComment, err = mysql.GetPostToCommentLatest(pc.PostId, pc.ToCommentId)
	} else {
		latestComment, err = mysql.GetPostCommentLatest(pc.PostId)
	}
	if err != nil {
		return nil, err
	} else if latestComment == nil {
		pc.Floor = 1
	} else {
		pc.Floor = latestComment.Floor + 1
	}
	// mysql
	pc.UserId = uid
	pc.CoupleId = cid
	pc, err = mysql.AddPostComment(pc)
	if pc == nil || err != nil {
		return pc, err
	}
	// 同步
	go func() {
		// post
		p.CommentCount = p.CommentCount + 1
		UpdatePostCount(p, true)
		// comment
		if toComment != nil && toComment.Id > 0 {
			toComment.SubCommentCount = toComment.SubCommentCount + 1
			UpdatePostCommentCount(toComment, true)
		}
		// message
		if pc.Kind == entity.POST_COMMENT_KIND_JAB {
			// jab 不给楼主和被评论者message，只给被jab的人
			if taId > 0 {
				language := "zh-cn"
				entry, err := mysql.GetEntryLatestByUser(taId)
				if err == nil && entry != nil {
					language = entry.Language
				}
				if toComment == nil || toComment.Id <= 0 {
					// 帖子评论
					content := utils.GetLanguage(language, "topic_message_jab_post") + p.Title
					message := CreateTopicMessage(uid, cid, taId, cid, entity.TOPIC_MESSAGE_KIND_JAB_IN_POST, content, p.Id)
					AddTopicMessage(language, message)
				} else {
					// 子评论
					content := utils.GetLanguage(language, "topic_message_jab_comment") + toComment.ContentText
					message := CreateTopicMessage(uid, cid, taId, cid, entity.TOPIC_MESSAGE_KIND_JAB_IN_COMMENT, content, toComment.Id)
					AddTopicMessage(language, message)
				}
			}
		} else {
			// text 分匿名情况
			if p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {

				if toComment == nil || toComment.Id <= 0 {
					// 帖子评论
					if uid != p.UserId {
						language := "zh-cn"
						entry, err := mysql.GetEntryLatestByUser(p.UserId)
						if err == nil && entry != nil {
							language = entry.Language
						}
						content := utils.GetLanguage(language, "topic_message_post_comment") + p.Title
						message := CreateTopicMessage(uid, cid, p.UserId, p.CoupleId, entity.TOPIC_MESSAGE_KIND_POST_BE_COMMENT, content, p.Id)
						AddTopicMessage(language, message)
					}
				} else {
					// 子评论
					if uid != toComment.UserId {
						language := "zh-cn"
						entry, err := mysql.GetEntryLatestByUser(toComment.UserId)
						if err == nil && entry != nil {
							language = entry.Language
						}
						content := utils.GetLanguage(language, "topic_message_comment_reply") + toComment.ContentText
						message := CreateTopicMessage(uid, cid, toComment.UserId, toComment.CoupleId, entity.TOPIC_MESSAGE_KIND_COMMENT_BE_REPLY, content, toComment.Id)
						AddTopicMessage(language, message)
					}
				}
			}
		}
		// topicInfo
		TopicInfoUpdateComment(p.Kind, true)
	}()
	// redis-set
	redis.SetPostComment(pc)
	// TODO redis
	//if pc.ToCommentId <= 0 {
	//	redis.AddPostCommentInListByCreate(pc)
	//} else {
	//	redis.AddPostToCommentInListByCreate(pc)
	//}
	return pc, err
}

// DelPostComment
func DelPostComment(uid, cid, pcid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if pcid <= 0 {
		return errors.New("nil_comment")
	}
	// comment检查
	pc, err := GetPostCommentById(pcid)
	if err != nil {
		return err
	} else if pc == nil {
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		return errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pc.PostId)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("nil_post")
	}
	// admin
	u, _ := GetUserById(uid)
	admin := IsAdminister(u)
	if !admin && pc.UserId != u.Id {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelPostComment(pc)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		// post
		p.CommentCount = p.CommentCount - 1
		UpdatePostCount(p, false)
		// comment
		if pc.ToCommentId > 0 {
			commentParent, _ := GetPostCommentById(pc.ToCommentId)
			if commentParent != nil {
				commentParent.SubCommentCount = commentParent.SubCommentCount - 1
				UpdatePostCommentCount(commentParent, false)
			}
		}
		// message
		if admin {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(pc.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_official_del_comment") + pc.ContentText
			message := CreateTopicMessage(uid, cid, pc.UserId, pc.CoupleId, entity.TOPIC_MESSAGE_KIND_OFFICIAL_TEXT, content, pc.Id)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpdateComment(p.Kind, false)
	}()
	// redis-del
	redis.DelPostComment(pc)
	// TODO redis
	//if pc.ToCommentId <= 0 {
	//	redis.DelPostCommentInListByPoint(pc)
	//	redis.DelPostCommentInListByCreate(pc)
	//} else {
	//	redis.DelPostToCommentInListByPoint(pc)
	//	redis.DelPostToCommentInListByCreate(pc)
	//}
	return err
}

// UpdatePostCount
func UpdatePostCommentCount(pc *entity.PostComment, update bool) (*entity.PostComment, error) {
	if pc == nil || pc.Id <= 0 {
		return nil, errors.New("nil_comment")
	}
	// redis-del
	redis.DelPostComment(pc)
	// mysql
	pc, err := mysql.UpdatePostCommentCount(pc, update)
	if pc == nil || err != nil {
		return pc, err
	}
	// redis-set
	redis.SetPostComment(pc)
	// TODO redis
	//if pc.ToCommentId <= 0 {
	//	redis.UpdatePostCommentInListByPoint(pc)
	//	redis.UpdatePostCommentInListByCreate(pc)
	//} else {
	//	redis.UpdatePostToCommentInListByPoint(pc)
	//	redis.UpdatePostToCommentInListByCreate(pc)
	//}
	return pc, err
}

// GetPostCommentById
func GetPostCommentById(pcid int64) (*entity.PostComment, error) {
	if pcid <= 0 {
		return nil, errors.New("nil_comment")
	}
	// redis-get
	pc, err := redis.GetPostCommentById(pcid)
	if pc != nil && pc.Id > 0 && err == nil {
		return pc, err
	}
	// mysql
	pc, err = mysql.GetPostCommentById(pcid)
	if pc == nil || pc.Id <= 0 || err != nil {
		return pc, err
	}
	// redis-set
	redis.SetPostComment(pc)
	return pc, err
}

// GetPostCommentByIdWithAll
func GetPostCommentByIdWithAll(uid, cid, pcid int64) (*entity.PostComment, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if pcid <= 0 {
		return nil, errors.New("nil_comment")
	}
	// mysql
	pc, err := GetPostCommentById(pcid)
	if err != nil {
		return nil, err
	} else if pc == nil {
		return nil, errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pc.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// 额外属性
	LoadPostCommentWithAll(uid, cid, p.Kind, pc)
	return pc, err
}

// GetPostCommentListByPost
func GetPostCommentListByPost(uid, cid, pid int64, order int, page int) ([]*entity.PostComment, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if pid <= 0 {
		return nil, errors.New("nil_post")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// post检查
	p, err := GetPostById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	var list []*entity.PostComment
	limit := GetPageSizeLimit().PostComment
	offset := page * limit
	// TODO redis
	//if order == POST_COMMENT_ORDER_CREATE {
	//	list, _ = redis.GetPostCommentListByCreate(pid, offset, limit)
	//} else {
	//	list, _ = redis.GetPostCommentListByPoint(pid, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		limitReportCount := GetLimit().PostCommentScreenReportCount
		orderBy := PostCommentOrderList[POST_COMMENT_ORDER_POINT]
		if order > 0 && order < len(PostCommentOrderList) {
			orderBy = PostCommentOrderList[order]
		}
		list, err = mysql.GetPostCommentListByPost(pid, limitReportCount, orderBy, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if order == POST_COMMENT_ORDER_CREATE {
		//		redis.SetPostCommentListByCreate(pid, list, true)
		//	} else {
		//		redis.SetPostCommentListByPoint(pid, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		LoadPostCommentWithAll(uid, cid, p.Kind, v)
	}
	return list, err
}

// GetPostToCommentList
func GetPostToCommentList(uid, cid, pid, tcid int64, order int, page int) ([]*entity.PostComment, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if pid <= 0 {
		return nil, errors.New("nil_post")
	} else if tcid <= 0 {
		return nil, errors.New("nil_comment")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// post检查
	p, err := GetPostById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	var list []*entity.PostComment
	limit := GetPageSizeLimit().PostComment
	offset := page * limit
	// TODO redis
	//if order == POST_COMMENT_ORDER_CREATE {
	//	list, _ = redis.GetPostToCommentListByCreate(tcid, offset, limit)
	//} else {
	//	list, _ = redis.GetPostToCommentListByPoint(tcid, offset, limit)
	//}
	// mysql
	if list == nil || len(list) <= 0 {
		limitReportCount := GetLimit().PostCommentScreenReportCount
		orderBy := PostCommentOrderList[POST_COMMENT_ORDER_POINT]
		if order > 0 && order < len(PostCommentOrderList) {
			orderBy = PostCommentOrderList[order]
		}
		list, err = mysql.GetPostToCommentList(pid, tcid, limitReportCount, orderBy, offset, limit)
		// TODO redis
		//if list != nil && len(list) > 0 && err == nil {
		//	// 再存到redis里
		//	if order == POST_COMMENT_ORDER_CREATE {
		//		redis.SetPostToCommentListByCreate(tcid, list, true)
		//	} else {
		//		redis.SetPostToCommentListByPoint(tcid, list, true)
		//	}
		//}
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		LoadPostCommentWithAll(uid, cid, p.Kind, v)
	}
	// 没有同步
	return list, err
}

// GetPostCommentListByUserPost
func GetPostCommentListByUserPost(uid, cid, pid, suid int64, order int, page int) ([]*entity.PostComment, error) {
	if uid <= 0 || suid <= 0 {
		return nil, errors.New("nil_user")
	} else if pid <= 0 {
		return nil, errors.New("nil_post")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// post检查
	p, err := GetPostById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// mysql
	limit := GetPageSizeLimit().PostComment
	offset := page * limit
	limitReportCount := GetLimit().PostCommentScreenReportCount
	orderBy := PostCommentOrderList[POST_COMMENT_ORDER_POINT]
	if order > 0 && order < len(PostCommentOrderList) {
		orderBy = PostCommentOrderList[order]
	}
	list, err := mysql.GetPostCommentListByUserPost(suid, pid, limitReportCount, orderBy, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		LoadPostCommentWithAll(uid, cid, p.Kind, v)
	}
	// 没有同步
	return list, err
}

// GetPostCommentList
func GetPostCommentList(uid, pid, tcid int64, page int) ([]*entity.PostComment, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().PostComment
	offset := page * limit
	list, err := mysql.GetPostCommentList(uid, pid, tcid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	return list, nil
}

// GetPostCommentReportList
func GetPostCommentReportList(page int) ([]*entity.PostComment, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// reportList
	limit := GetPageSizeLimit().PostComment
	offset := page * limit
	reportList, err := mysql.GetPostCommentReportList(offset, limit)
	if err != nil {
		return nil, err
	} else if reportList == nil || len(reportList) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	// postList
	list := make([]*entity.PostComment, 0)
	for _, v := range reportList {
		if v == nil || v.PostCommentId <= 0 {
			continue
		}
		pc, _ := GetPostCommentById(v.PostCommentId)
		if pc != nil && pc.Id > 0 {
			list = append(list, pc)
		}
	}
	return list, nil
}

// GetPostCommentTotalByCreateWithDel
func GetPostCommentTotalByCreateWithDel(create int64) int64 {
	if create == 0 {
		return 0
	}
	// mysql
	total := mysql.GetPostCommentTotalByCreateWithDel(create)
	return total
}

// LoadPostCommentWithAll
func LoadPostCommentWithAll(uid, cid int64, pk int, pc *entity.PostComment) *entity.PostComment {
	if pc == nil || pc.Id <= 0 {
		return nil
	}
	// 额外属性
	if pc.Official {
		pc.Screen = false
	} else {
		if pc.ReportCount < GetLimit().PostCommentScreenReportCount {
			pc.Screen = false
		} else {
			pc.Screen = true
		}
	}
	if pk == entity.POST_KIND_LIMIT_UNKNOWN {
		// 匿名
		if pc.Kind == entity.POST_COMMENT_KIND_JAB {
			// 戳
			pc.Couple, _ = GetCoupleVisibleByUser(pc.UserId)
		} else {
			// 文本
			pc.Couple = nil
		}
	} else {
		pc.Couple, _ = GetCoupleVisibleByUser(pc.UserId)
	}
	if cid <= 0 {
		// 无配对
		pc.Mine = false
		pc.Our = false
		pc.SubComment = false
		pc.Report = false
		pc.Point = false
	} else {
		if pk == entity.POST_KIND_LIMIT_UNKNOWN {
			// 匿名
			pc.Mine = false
			pc.Our = false
			pc.SubComment = false
		} else {
			pc.Mine = pc.UserId == uid
			pc.Our = pc.CoupleId == cid
			if pc.SubCommentCount > 0 {
				pc.SubComment = IsPostToCommentByUserCouple(uid, cid, pc.PostId, pc.Id)
			} else {
				pc.SubComment = false
			}
		}
		pc.Report = IsPostCommentReportByUserCouple(uid, cid, pc.Id)
		pc.Point = IsPostCommentPointByUserCouple(uid, cid, pc.Id)
	}
	return pc
}

// IsPostCommentByUserCouple
func IsPostCommentByUserCouple(uid, cid, pid int64) bool {
	if uid <= 0 || cid <= 0 || pid <= 0 {
		return false
	}
	return mysql.GetPostCommentTotalByUserCouple(uid, cid, pid) > 0
}

// IsPostToCommentByUserCouple
func IsPostToCommentByUserCouple(uid, cid, pid, tcid int64) bool {
	if uid <= 0 || cid <= 0 || pid <= 0 || tcid <= 0 {
		return false
	}
	return mysql.GetPostToCommentTotalByUserCouple(uid, cid, pid, tcid) > 0
}
