package entity

import (
	"strconv"
	"strings"
)

// SplitStrByColon 分号分割str
func SplitStrByColon(str string) []string {
	strList := make([]string, 0)
	strTrim := strings.TrimSpace(str)
	if len(strTrim) <= 0 {
		return strList
	}
	urlSplit := strings.Split(strTrim, ";")
	for _, v := range urlSplit {
		trim := strings.TrimSpace(v)
		if len(trim) > 0 || trim != "" {
			strList = append(strList, trim)
		}
	}
	return strList
}

// JoinStrByColon 分号连接str
func JoinStrByColon(strList []string) string {
	if strList == nil || len(strList) <= 0 {
		return ""
	}
	b := make([]string, 0)
	for _, v := range strList {
		v := strings.TrimSpace(v)
		if len(v) > 0 && v != "" {
			b = append(b, v)
		}
	}
	join := strings.Join(b, ";")
	return join
}

// SplitInt64ByColon 分号分割int64
func SplitInt64ByColon(str string) []int64 {
	int64List := make([]int64, 0)
	strTrim := strings.TrimSpace(str)
	if len(strTrim) <= 0 {
		return int64List
	}
	urlSplit := strings.Split(strTrim, ";")
	for _, v := range urlSplit {
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if i != 0 && err == nil {
			int64List = append(int64List, i)
		}
	}
	return int64List
}

// JoinInt64ByColon 分号连接int64
func JoinInt64ByColon(int64List []int64) string {
	if int64List == nil || len(int64List) <= 0 {
		return ""
	}
	b := make([]string, 0)
	for _, v := range int64List {
		s := strings.TrimSpace(strconv.FormatInt(v, 10))
		if len(s) > 0 && s != "" {
			b = append(b, s)
		}
	}
	join := strings.Join(b, ";")
	return join
}

// GetPostSubKindMap (位置可以随意调换)
func GetPostSubKindMap() map[int]map[int]int {
	if len(postSubKindMap) > 0 {
		return postSubKindMap
	}
	postSubKindMap[POST_KIND_OPEN_LIVE] = make(map[int]int, 0)
	postSubKindMap[POST_KIND_OPEN_LIVE][0] = POST_SUB_KIND_ALL           // 全部
	postSubKindMap[POST_KIND_OPEN_LIVE][1] = POST_SUB_KIND_LIVE_LIFE     // 小日常
	postSubKindMap[POST_KIND_OPEN_LIVE][2] = POST_SUB_KIND_LIVE_PLACES   // 异地恋
	postSubKindMap[POST_KIND_OPEN_LIVE][3] = POST_SUB_KIND_LIVE_FRAGMENT // 碎碎念
	postSubKindMap[POST_KIND_OPEN_LIVE][4] = POST_SUB_KIND_LIVE_MARRY    // 结婚记
	postSubKindMap[POST_KIND_OPEN_LIVE][5] = POST_SUB_KIND_LIVE_HAPPY    // 吃喝玩乐
	postSubKindMap[POST_KIND_OPEN_LIVE][6] = POST_SUB_KIND_LIVE_KNOW     // 与TA相识
	postSubKindMap[POST_KIND_OPEN_LIVE][7] = POST_SUB_KIND_LIVE_TRAVEL   // 带tA旅行
	postSubKindMap[POST_KIND_OPEN_LIVE][8] = POST_SUB_KIND_LIVE_HOUSE    // 和TA同居
	postSubKindMap[POST_KIND_OPEN_STAR] = make(map[int]int, 0)
	postSubKindMap[POST_KIND_OPEN_STAR][0] = POST_SUB_KIND_ALL          // 全部
	postSubKindMap[POST_KIND_OPEN_STAR][1] = POST_SUB_KIND_STAR_SHEEP   // 白羊
	postSubKindMap[POST_KIND_OPEN_STAR][2] = POST_SUB_KIND_STAR_MILK    // 金牛
	postSubKindMap[POST_KIND_OPEN_STAR][3] = POST_SUB_KIND_STAR_SON     // 双子
	postSubKindMap[POST_KIND_OPEN_STAR][4] = POST_SUB_KIND_STAR_HUGE    // 巨蟹
	postSubKindMap[POST_KIND_OPEN_STAR][5] = POST_SUB_KIND_STAR_LION    // 狮子
	postSubKindMap[POST_KIND_OPEN_STAR][6] = POST_SUB_KIND_STAR_GIRL    // 处女
	postSubKindMap[POST_KIND_OPEN_STAR][7] = POST_SUB_KIND_STAR_BALANCE // 天秤
	postSubKindMap[POST_KIND_OPEN_STAR][8] = POST_SUB_KIND_STAR_SKY     // 天蝎
	postSubKindMap[POST_KIND_OPEN_STAR][9] = POST_SUB_KIND_STAR_HAND    // 射手
	postSubKindMap[POST_KIND_OPEN_STAR][10] = POST_SUB_KIND_STAR_DEVIL  // 摩羯
	postSubKindMap[POST_KIND_OPEN_STAR][11] = POST_SUB_KIND_STAR_WATER  // 水瓶
	postSubKindMap[POST_KIND_OPEN_STAR][12] = POST_SUB_KIND_STAR_FISH   // 双鱼
	postSubKindMap[POST_KIND_OPEN_ANIMAL] = make(map[int]int, 0)
	postSubKindMap[POST_KIND_OPEN_ANIMAL][0] = POST_SUB_KIND_ALL           // 全部
	postSubKindMap[POST_KIND_OPEN_ANIMAL][1] = POST_SUB_KIND_ANIMAL_CAT    // 猫系
	postSubKindMap[POST_KIND_OPEN_ANIMAL][2] = POST_SUB_KIND_ANIMAL_DOG    // 犬系
	postSubKindMap[POST_KIND_OPEN_ANIMAL][3] = POST_SUB_KIND_ANIMAL_WOLF   // 狼系
	postSubKindMap[POST_KIND_OPEN_ANIMAL][4] = POST_SUB_KIND_ANIMAL_FOX    // 狐系
	postSubKindMap[POST_KIND_OPEN_ANIMAL][5] = POST_SUB_KIND_ANIMAL_RABBIT // 兔系
	//postSubKindMap[POST_KIND_OPEN_ANIMAL][6] = POST_SUB_KIND_ANIMAL_TIGER  // 虎系
	//postSubKindMap[POST_KIND_OPEN_ANIMAL][7] = POST_SUB_KIND_ANIMAL_BEE    // 蜂系
	//postSubKindMap[POST_KIND_OPEN_ANIMAL][8] = POST_SUB_KIND_ANIMAL_MILK   // 牛系
	//postSubKindMap[POST_KIND_OPEN_ANIMAL][9] = POST_SUB_KIND_ANIMAL_MONKEY // 猴系
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN] = make(map[int]int, 0)
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN][0] = POST_SUB_KIND_ALL           // 全部
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN][1] = POST_SUB_KIND_UNKNOWN_GIFT  // 礼物
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN][2] = POST_SUB_KIND_UNKNOWN_ROMAN // 浪漫
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN][3] = POST_SUB_KIND_UNKNOWN_BRAWL // 吵架
	postSubKindMap[POST_KIND_LIMIT_UNKNOWN][4] = POST_SUB_KIND_UNKNOWN_TREE  // 树洞
	//postSubKindMap[POST_KIND_LIMIT_UNKNOWN][5] = POST_SUB_KIND_UNKNOWN_SHY   // 羞羞
	return postSubKindMap
}
