package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// AddPostReport
func AddPostReport(uid, cid int64, pr *entity.PostReport) (*entity.PostReport, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pr == nil || pr.PostId <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pr.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	} else if p.Top || p.Official || p.Well {
		return nil, errors.New("report_refuse")
	}
	// old
	old, err := mysql.GetPostReportByUserCouple(uid, cid, p.Id)
	if err != nil {
		return nil, err
	} else if old != nil {
		return nil, errors.New("report_repeat")
	}
	// mysql
	pr.UserId = uid
	pr.CoupleId = cid
	pr, err = mysql.AddPostReport(pr)
	if pr == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// post
		p.ReportCount = p.ReportCount + 1
		UpdatePostCount(p, false)
		// message
		if uid != p.UserId && p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(p.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_post_report") + p.Title
			message := CreateTopicMessage(uid, cid, p.UserId, p.CoupleId, entity.TOPIC_MESSAGE_KIND_POST_BE_REPORT, content, p.Id)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpReport(p.Kind)
	}()
	return pr, err
}

// IsPostReportByUserCouple
func IsPostReportByUserCouple(uid, cid, pid int64) bool {
	if uid <= 0 || cid <= 0 || pid <= 0 {
		return false
	}
	report, _ := mysql.GetPostReportByUserCouple(uid, cid, pid)
	if report == nil || report.Id <= 0 {
		return false
	}
	return true
}
