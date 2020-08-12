package controllers

import (
	"net/http"

	"libs/utils"
	"models/entity"
	"services"
	"strconv"
	"strings"
)

func HandlerEntry(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "entry")
	if r.Method == http.MethodPost {
		PostEntry(w, r)
	} else if r.Method == http.MethodGet {
		GetEntry(w, r)
	} else {
		response405(w, r)
	}
}

// PostEntry
// 1.app进入时调用 所有app端不耗时能获取到的初始化数据都这里处理
func PostEntry(w http.ResponseWriter, r *http.Request) {
	// params
	entry := &entity.Entry{}
	checkRequestBody(w, r, entry)
	entryPlatform := strings.ToLower(strings.TrimSpace(entry.Platform))
	platformAndroid := utils.GetConfigStr("conf", "app.conf", "phone", "platform_android")
	platformIos := utils.GetConfigStr("conf", "app.conf", "phone", "platform_ios")
	// commonConst
	commonConst := &services.CommonConst{
		CompanyName:   utils.GetConfigStr("conf", "app.conf", "common", "company_name"),
		CustomerQQ:    utils.GetConfigStr("conf", "app.conf", "common", "customer_qq"),
		OfficialGroup: utils.GetConfigStr("conf", "app.conf", "common", "official_group"),
		OfficialWeibo: utils.GetConfigStr("conf", "app.conf", "common", "official_weibo"),
		OfficialWeb:   utils.GetConfigStr("conf", "app.conf", "common", "official_web"),
		ContactEmail:  utils.GetConfigStr("conf", "app.conf", "common", "contact_email"),
		IosAppId:      utils.GetConfigStr("conf", "app.conf", "common", "ios_app_id"),
	}
	// limit
	limit := services.GetLimit()
	// user
	user, _ := getTokenCouple(r)
	user.Password = ""
	taId := services.GetTaId(user)
	cid := services.GetCoupleIdByUser(user)
	// oss
	var ossInfo *services.OssInfo
	if getModelStatus("oss") {
		ossInfo, _ = services.GetOssInfoByUserCouple(user.Id, cid, user.UserToken)
	}
	// version
	versionList := make([]*entity.Version, 0)
	if getModelStatus("version") && entry != nil && entry.AppVersion > 0 {
		// 检查最低的强制更新版本
		if entryPlatform == platformAndroid {
			versionList, _ = services.GetVersionListByPlatformCode(entry.Platform, entry.AppVersion)
			limitAndroid := utils.GetConfigInt("conf", "model.conf", "app_version", "android")
			if entry.AppVersion < limitAndroid {
				response409(w, r, "", struct {
					OssInfo     *services.OssInfo `json:"ossInfo"`
					VersionList []*entity.Version `json:"versionList"`
				}{
					ossInfo,
					versionList,
				})
			}
		} else if entryPlatform == platformIos {
			limitIos := utils.GetConfigInt("conf", "model.conf", "app_version", "ios")
			if entry.AppVersion < limitIos {
				response409(w, r, "", struct{}{})
			}
		}
	}
	// entry 需要在version之后
	entry.UserId = user.Id
	entry, _ = services.AddEntry(entry)
	// modelShow
	marketPayShow := true
	if len(strings.TrimSpace(entry.Market)) > 0 {
		markets := utils.GetConfigStr("conf", "model.conf", "model", "markets_hide_pay")
		if len(markets) > 0 && len(strings.Split(markets, ";")) > 0 {
			marketList := strings.Split(markets, ";")
			for _, v := range marketList {
				if strings.TrimSpace(v) == strings.TrimSpace(entry.Market) {
					marketPayShow = false
					break
				}
			}
		}
		if marketPayShow == true {
			marketsTestPhone := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_phone")
			marketsTest := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_hide_pay")
			if len(marketsTestPhone) > 0 && len(strings.Split(marketsTestPhone, ";")) > 0 && len(marketsTest) > 0 && len(strings.Split(marketsTest, ";")) > 0 {
				phoneList := strings.Split(marketsTestPhone, ";")
				for _, phone := range phoneList {
					if strings.TrimSpace(phone) == strings.TrimSpace(user.Phone) {
						marketList := strings.Split(marketsTest, ";")
						for _, market := range marketList {
							if strings.TrimSpace(market) == strings.TrimSpace(entry.Market) {
								marketPayShow = false
								break
							}
						}
						break
					}
				}
			}
		}
	}
	marketCoinAd := true
	if len(strings.TrimSpace(entry.Market)) > 0 {
		markets := utils.GetConfigStr("conf", "model.conf", "model", "markets_hide_coin_ad")
		if len(markets) > 0 && len(strings.Split(markets, ";")) > 0 {
			marketList := strings.Split(markets, ";")
			for _, v := range marketList {
				if strings.TrimSpace(v) == strings.TrimSpace(entry.Market) {
					marketCoinAd = false
					break
				}
			}
		}
		if marketCoinAd == true {
			marketsTestPhone := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_phone")
			marketsTest := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_hide_coin_ad")
			if len(marketsTestPhone) > 0 && len(strings.Split(marketsTestPhone, ";")) > 0 && len(marketsTest) > 0 && len(strings.Split(marketsTest, ";")) > 0 {
				phoneList := strings.Split(marketsTestPhone, ";")
				for _, phone := range phoneList {
					if strings.TrimSpace(phone) == strings.TrimSpace(user.Phone) {
						marketList := strings.Split(marketsTest, ";")
						for _, market := range marketList {
							if strings.TrimSpace(market) == strings.TrimSpace(entry.Market) {
								marketCoinAd = false
								break
							}
						}
						break
					}
				}
			}
		}
	}
	modelShow := &services.ModelShow{
		MarketPay:     marketPayShow,
		MarketCoinAd:  marketCoinAd,
		Couple:        getModelStatus("couple"),
		CouplePlace:   getModelStatus("couple_place"),
		CoupleWeather: getModelStatus("couple_weather"),
		Note:          getModelStatus("note"),
		Topic:         getModelStatus("topic"),
		More:          getModelStatus("more"),
		MoreVip:       getModelStatus("more_vip"),
		MoreCoin:      getModelStatus("more_coin"),
		MoreMatch:     getModelStatus("more_match"),
		MoreFeature:   getModelStatus("more_feature"),
	}
	// cooperation
	marketBaby := utils.GetConfigBool("conf", "third.conf", "cooperation", "catch_baby_enable")
	if len(strings.TrimSpace(entry.Market)) > 0 {
		marketsTestPhone := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_phone")
		marketsTest := utils.GetConfigStr("conf", "model.conf", "model", "markets_test_hide_baby")
		if len(marketsTestPhone) > 0 && len(strings.Split(marketsTestPhone, ";")) > 0 && len(marketsTest) > 0 && len(strings.Split(marketsTest, ";")) > 0 {
			phoneList := strings.Split(marketsTestPhone, ";")
			for _, phone := range phoneList {
				if strings.TrimSpace(phone) == strings.TrimSpace(user.Phone) {
					marketList := strings.Split(marketsTest, ";")
					for _, market := range marketList {
						if strings.TrimSpace(market) == strings.TrimSpace(entry.Market) {
							marketBaby = false
							break
						}
					}
					break
				}
			}
		}
	}
	cooperation := &services.Cooperation{
		CatchBabyEnable: marketBaby,
		CatchBabyUrl:    "http://www.woyaozhua.com/h5/index.html#/loginHome?useInviteCode=019727", // 有#还是在这里粘贴吧
		CatchBabyDesc:   utils.GetConfigStr("conf", "third.conf", "cooperation", "catch_baby_desc"),
	}
	// push
	var pushInfo *services.PushInfo
	if getModelStatus("push") {
		pushInfo = services.GetPushInfo(user.Id)
	}
	// ad
	var adInfo *services.AdInfo
	if getModelStatus("ad") {
		if entryPlatform == platformIos {
			adInfo = &services.AdInfo{
				AppId:           utils.GetConfigStr("conf", "third.conf", "qq-ad", "ios_app_id"),
				TopicPostPosId:  utils.GetConfigStr("conf", "third.conf", "qq-ad", "ios_topic_post_pos_id"),
				TopicPostStart:  utils.GetConfigInt("conf", "third.conf", "qq-ad", "ios_topic_post_start"),
				TopicPostJump:   utils.GetConfigInt("conf", "third.conf", "qq-ad", "ios_topic_post_jump"),
				CoinFreePosId:   utils.GetConfigStr("conf", "third.conf", "qq-ad", "ios_coin_free_pos_id"),
				CoinFreeTickSec: utils.GetConfigInt("conf", "third.conf", "qq-ad", "ios_coin_free_tick_sec"),
			}
		} else {
			adInfo = &services.AdInfo{
				AppId:           utils.GetConfigStr("conf", "third.conf", "qq-ad", "android_app_id"),
				TopicPostPosId:  utils.GetConfigStr("conf", "third.conf", "qq-ad", "android_topic_post_pos_id"),
				TopicPostStart:  utils.GetConfigInt("conf", "third.conf", "qq-ad", "android_topic_post_start"),
				TopicPostJump:   utils.GetConfigInt("conf", "third.conf", "qq-ad", "android_topic_post_jump"),
				CoinFreePosId:   utils.GetConfigStr("conf", "third.conf", "qq-ad", "android_coin_free_pos_id"),
				CoinFreeTickSec: utils.GetConfigInt("conf", "third.conf", "qq-ad", "android_coin_free_tick_sec"),
			}
		}
	}
	// vipLimit
	vipLimit := services.GetVipLimitByCouple(cid)
	// commonCount
	commonCount := &services.CommonCount{
		NoticeNewCount:     services.GetNoticeCountByNoRead(user.Id),
		VersionNewCount:    0,
		NoteTrendsNewCount: services.GetTrendsCountByUserCouple(user.Id, taId, cid),
		TopicMsgNewCount:   services.GetTopicMessageCountByUserCouple(user.Id, cid),
	}
	if versionList != nil && len(versionList) > 0 {
		commonCount.VersionNewCount = len(versionList)
	}
	// 返回
	response200Data(w, r, struct {
		CommonConst *services.CommonConst `json:"commonConst"`
		Limit       *services.Limit       `json:"limit"`
		User        *entity.User          `json:"user"`
		VersionList []*entity.Version     `json:"versionList"`
		ModelShow   *services.ModelShow   `json:"modelShow"`
		Cooperation *services.Cooperation `json:"cooperation"`
		OssInfo     *services.OssInfo     `json:"ossInfo"`
		PushInfo    *services.PushInfo    `json:"pushInfo"`
		AdInfo      *services.AdInfo      `json:"adInfo"`
		VipLimit    *services.VipLimit    `json:"vipLimit"`
		CommonCount *services.CommonCount `json:"commonCount"`
	}{
		commonConst,
		limit,
		user,
		versionList,
		modelShow,
		cooperation,
		ossInfo,
		pushInfo,
		adInfo,
		vipLimit,
		commonCount,
	})
}

// GetEntry
func GetEntry(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	total, _ := strconv.ParseBool(values.Get("total"))
	group, _ := strconv.ParseBool(values.Get("group"))
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	if list {
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		entryList, err := services.GetEntryList(uid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			EntryList []*entity.Entry `json:"entryList"`
		}{entryList})
	} else if total {
		create, _ := strconv.ParseBool(values.Get("create"))
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		var count int64
		if create {
			count = services.GetEntryTotalByCreate(start, end)
		} else {
			count = services.GetEntryTotalByUpdate(start, end)
		}
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{count})
	} else if group {
		filed := values.Get("filed")
		create, _ := strconv.ParseBool(values.Get("create"))
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		var at string
		if create {
			at = "create_at"
		} else {
			at = "update_at"
		}
		infoList := make([]*entity.FiledInfo, 0)
		if filed == "deviceName" {
			infoList, _ = services.GetEntryGroupDeviceNameList(at, start, end)
		} else if filed == "market" {
			infoList, _ = services.GetEntryGroupMarketList(at, start, end)
		} else if filed == "language" {
			infoList, _ = services.GetEntryGroupLanguageList(at, start, end)
		} else if filed == "platform" {
			infoList, _ = services.GetEntryGroupPlatformList(at, start, end)
		} else if filed == "osVersion" {
			infoList, _ = services.GetEntryGroupOsVersionList(at, start, end)
		} else if filed == "appVersion" {
			infoList, _ = services.GetEntryGroupAppVersionList(at, start, end)
		}
		// 返回
		response200Data(w, r, struct {
			InfoList []*entity.FiledInfo `json:"infoList"`
		}{infoList})
	} else {
		response405(w, r)
	}
}
