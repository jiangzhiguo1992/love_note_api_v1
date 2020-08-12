package services

import (
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// GetPostKindInfoListWithTopicInfo
func GetPostKindInfoListWithTopicInfo(me *entity.User) []*PostKindInfo {
	postKindInfoList := GetPostKindInfoList(me)
	// 循环postKindInfo 取 topicInfo
	if postKindInfoList != nil && len(postKindInfoList) > 0 {
		for _, kindInfo := range postKindInfoList {
			if kindInfo == nil || !kindInfo.Enable {
				continue
			}
			now := utils.GetCSTDateByUnix(time.Now().Unix())
			topicInfo, err := mysql.GetTopicInfoByKindYearDays(kindInfo.Kind, now.Year(), now.YearDay())
			if err == nil {
				// err不是nil不能操作 防止重复
				if topicInfo == nil {
					// 添加新的一天
					topicInfo = &entity.TopicInfo{Kind: kindInfo.Kind, Year: now.Year(), DayOfYear: now.YearDay()}
					topicInfo, _ = mysql.AddTopicInfo(topicInfo)
				}
				kindInfo.TopicInfo = topicInfo
			}
		}
	}
	return postKindInfoList
}

// GetPostKindInfoList
// 1.记得修改kind和subKind之后也修改这里
func GetPostKindInfoList(me *entity.User) []*PostKindInfo {
	// language
	language := "zh-cn"
	entry, err := mysql.GetEntryLatestByUser(me.Id)
	if err == nil && entry != nil {
		language = entry.Language
	}
	// kindInfo
	kindInfoList := make([]*PostKindInfo, 0)
	for i := 0; i < len(entity.PostKindList); i++ {
		kindId := entity.PostKindList[i]
		kindInfo := &PostKindInfo{}
		kindInfo.Kind = kindId
		kindInfo.Enable = GetPostKindEnable(kindId)
		kindInfo.Name = GetPostKindName(language, kindId, -1)
		// subKindInfo
		kindInfo.PostSubKindInfoList = make([]*PostSubKindInfo, 0)
		subKindMap := entity.GetPostSubKindMap()[kindId]
		for j := 0; j < len(subKindMap); j++ {
			subKindId := subKindMap[j]
			subKindInfo := &PostSubKindInfo{}
			subKindInfo.Kind = subKindId
			subKindInfo.Enable = GetPostSubKindEnable(kindId, subKindId)
			subKindInfo.Name = GetPostKindName(language, kindId, subKindId)
			subKindInfo.Push = subKindId != entity.POST_SUB_KIND_ALL
			subKindInfo.Anonymous = kindId == entity.POST_KIND_LIMIT_UNKNOWN
			kindInfo.PostSubKindInfoList = append(kindInfo.PostSubKindInfoList, subKindInfo)
		}
		kindInfoList = append(kindInfoList, kindInfo)
	}
	return kindInfoList
}

// GetPostKindName
func GetPostKindName(language string, kind, subKind int) string {
	name := ""
	if kind <= 0 {
		return name
	}
	nameSuffix := ""
	if kind == entity.POST_KIND_OPEN_LIVE {
		nameSuffix = "_live"
		if subKind >= 0 {
			if subKind == entity.POST_SUB_KIND_ALL {
				nameSuffix = "_all"
			} else if subKind == entity.POST_SUB_KIND_LIVE_LIFE {
				nameSuffix = nameSuffix + "_life"
			} else if subKind == entity.POST_SUB_KIND_LIVE_PLACES {
				nameSuffix = nameSuffix + "_places"
			} else if subKind == entity.POST_SUB_KIND_LIVE_FRAGMENT {
				nameSuffix = nameSuffix + "_fragment"
			} else if subKind == entity.POST_SUB_KIND_LIVE_MARRY {
				nameSuffix = nameSuffix + "_marry"
			} else if subKind == entity.POST_SUB_KIND_LIVE_HAPPY {
				nameSuffix = nameSuffix + "_happy"
			} else if subKind == entity.POST_SUB_KIND_LIVE_KNOW {
				nameSuffix = nameSuffix + "_know"
			} else if subKind == entity.POST_SUB_KIND_LIVE_TRAVEL {
				nameSuffix = nameSuffix + "_travel"
			} else if subKind == entity.POST_SUB_KIND_LIVE_HOUSE {
				nameSuffix = nameSuffix + "_house"
			}
		}
	} else if kind == entity.POST_KIND_OPEN_STAR {
		nameSuffix = "_star"
		if subKind == entity.POST_SUB_KIND_ALL {
			nameSuffix = "_all"
		} else if subKind == entity.POST_SUB_KIND_STAR_SHEEP {
			nameSuffix = nameSuffix + "_sheep"
		} else if subKind == entity.POST_SUB_KIND_STAR_MILK {
			nameSuffix = nameSuffix + "_milk"
		} else if subKind == entity.POST_SUB_KIND_STAR_SON {
			nameSuffix = nameSuffix + "_son"
		} else if subKind == entity.POST_SUB_KIND_STAR_HUGE {
			nameSuffix = nameSuffix + "_huge"
		} else if subKind == entity.POST_SUB_KIND_STAR_LION {
			nameSuffix = nameSuffix + "_lion"
		} else if subKind == entity.POST_SUB_KIND_STAR_GIRL {
			nameSuffix = nameSuffix + "_girl"
		} else if subKind == entity.POST_SUB_KIND_STAR_BALANCE {
			nameSuffix = nameSuffix + "_balance"
		} else if subKind == entity.POST_SUB_KIND_STAR_SKY {
			nameSuffix = nameSuffix + "_sky"
		} else if subKind == entity.POST_SUB_KIND_STAR_HAND {
			nameSuffix = nameSuffix + "_hand"
		} else if subKind == entity.POST_SUB_KIND_STAR_DEVIL {
			nameSuffix = nameSuffix + "_devil"
		} else if subKind == entity.POST_SUB_KIND_STAR_WATER {
			nameSuffix = nameSuffix + "_water"
		} else if subKind == entity.POST_SUB_KIND_STAR_FISH {
			nameSuffix = nameSuffix + "_fish"
		}
	} else if kind == entity.POST_KIND_OPEN_ANIMAL {
		nameSuffix = "_animal"
		if subKind == entity.POST_SUB_KIND_ALL {
			nameSuffix = "_all"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_CAT {
			nameSuffix = nameSuffix + "_cat"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_DOG {
			nameSuffix = nameSuffix + "_dog"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_WOLF {
			nameSuffix = nameSuffix + "_wolf"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_FOX {
			nameSuffix = nameSuffix + "_fox"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_RABBIT {
			nameSuffix = nameSuffix + "_rabbit"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_TIGER {
			nameSuffix = nameSuffix + "_tiger"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_BEE {
			nameSuffix = nameSuffix + "_bee"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_MILK {
			nameSuffix = nameSuffix + "_milk"
		} else if subKind == entity.POST_SUB_KIND_ANIMAL_MONKEY {
			nameSuffix = nameSuffix + "_monkey"
		}
	} else if kind == entity.POST_KIND_LIMIT_UNKNOWN {
		nameSuffix = "_unknown"
		if subKind == entity.POST_SUB_KIND_ALL {
			nameSuffix = "_all"
		} else if subKind == entity.POST_SUB_KIND_UNKNOWN_GIFT {
			nameSuffix = nameSuffix + "_gift"
		} else if subKind == entity.POST_SUB_KIND_UNKNOWN_ROMAN {
			nameSuffix = nameSuffix + "_roman"
		} else if subKind == entity.POST_SUB_KIND_UNKNOWN_BRAWL {
			nameSuffix = nameSuffix + "_brawl"
		} else if subKind == entity.POST_SUB_KIND_UNKNOWN_TREE {
			nameSuffix = nameSuffix + "_tree"
		} else if subKind == entity.POST_SUB_KIND_UNKNOWN_SHY {
			nameSuffix = nameSuffix + "_shy"
		}
	}
	return utils.GetLanguage(language, "topic_post_name"+nameSuffix)
}

// TopicInfoUpBrowse
func TopicInfoUpBrowse(kind int) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			info.BrowseCount = info.BrowseCount + 1
			mysql.UpdateTopicInfo(info)
		}
	}
}

// TopicInfoUpdatePost
func TopicInfoUpdatePost(kind int, up bool) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			if up {
				info.PostCount = info.PostCount + 1
			} else {
				info.PostCount = info.PostCount - 1
			}
			mysql.UpdateTopicInfo(info)
		}
	}
}

// TopicInfoUpdateComment
func TopicInfoUpdateComment(kind int, up bool) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			if up {
				info.CommentCount = info.CommentCount + 1
			} else {
				info.CommentCount = info.CommentCount - 1
			}
			mysql.UpdateTopicInfo(info)
		}
	}
}

// TopicInfoUpReport
func TopicInfoUpReport(kind int) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			info.ReportCount = info.ReportCount + 1
			mysql.UpdateTopicInfo(info)
		}
	}
}

// TopicInfoUpdatePoint
func TopicInfoUpdatePoint(kind int, up bool) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			if up {
				info.PointCount = info.PointCount + 1
			} else {
				info.PointCount = info.PointCount - 1
			}
			mysql.UpdateTopicInfo(info)
		}
	}
}

// TopicInfoUpdateCollect
func TopicInfoUpdateCollect(kind int, up bool) {
	if kind <= 0 {
		return
	}
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	info, err := mysql.GetTopicInfoByKindYearDays(kind, now.Year(), now.YearDay())
	if err == nil {
		// err不是nil不能操作 防止重复
		if info == nil || info.Id <= 0 {
			info = &entity.TopicInfo{Kind: kind, Year: now.Year(), DayOfYear: now.YearDay()}
			info, _ = mysql.AddTopicInfo(info)
		}
		if info != nil {
			if up {
				info.CollectCount = info.CollectCount + 1
			} else {
				info.CollectCount = info.CollectCount - 1
			}
			mysql.UpdateTopicInfo(info)
		}
	}
}
