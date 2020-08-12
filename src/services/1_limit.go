package services

import (
	"libs/utils"
)

type (
	Limit struct {
		// common
		SmsCodeLength int `json:"smsCodeLength"`
		SmsEffectSec  int `json:"smsEffectSec"`
		SmsBetweenSec int `json:"smsBetweenSec"`
		SmsMaxSec     int `json:"smsMaxSec"`
		SmsMaxCount   int `json:"smsMaxCount"`
		// settings
		SuggestTitleLength          int `json:"suggestTitleLength"`
		SuggestContentLength        int `json:"suggestContentLength"`
		SuggestCommentContentLength int `json:"suggestCommentContentLength"`
		// couple
		CoupleInviteIntervalSec int64 `json:"coupleInviteIntervalSec"`
		CoupleBreakNeedSec      int64 `json:"coupleBreakNeedSec"`
		CoupleBreakSec          int64 `json:"coupleBreakSec"`
		CoupleNameLength        int   `json:"coupleNameLength"`
		// note
		NoteResExpireSec          int64 `json:"noteResExpireSec"`
		MensesMaxPerMonth         int   `json:"mensesMaxPerMonth"`
		MensesMaxCycleDay         int   `json:"mensesMaxCycleDay"`
		MensesMaxDurationDay      int   `json:"mensesMaxDurationDay"`
		MensesDefaultCycleDay     int   `json:"mensesDefaultCycleDay"`
		MensesDefaultDurationDay  int   `json:"mensesDefaultDurationDay"`
		ShyMaxPerDay              int   `json:"shyMaxPerDay"`
		ShySafeLength             int   `json:"shySafeLength"`
		ShyDescLength             int   `json:"shyDescLength"`
		SleepMaxPerDay            int   `json:"sleepMaxPerDay"`
		NoteSleepSuccessMinSec    int64 `json:"noteSleepSuccessMinSec"`
		NoteSleepSuccessMaxSec    int64 `json:"noteSleepSuccessMaxSec"`
		NoteLockLength            int   `json:"noteLockLength"`
		SouvenirTitleLength       int   `json:"souvenirTitleLength"`
		SouvenirForeignYearCount  int   `json:"souvenirForeignYearCount"`
		TravelPlaceCount          int   `json:"travelPlaceCount"`
		TravelVideoCount          int   `json:"travelVideoCount"`
		TravelFoodCount           int   `json:"travelFoodCount"`
		TravelMovieCount          int   `json:"travelMovieCount"`
		TravelAlbumCount          int   `json:"travelAlbumCount"`
		TravelDiaryCount          int   `json:"travelDiaryCount"`
		AudioTitleLength          int   `json:"audioTitleLength"`
		VideoTitleLength          int   `json:"videoTitleLength"`
		AlbumTitleLength          int   `json:"albumTitleLength"`
		PicturePushCount          int   `json:"picturePushCount"`
		WordContentLength         int   `json:"wordContentLength"`
		WhisperChannelLength      int   `json:"whisperChannelLength"`
		WhisperContentLength      int   `json:"whisperContentLength"`
		DiaryContentLength        int   `json:"diaryContentLength"`
		AwardContentLength        int   `json:"awardContentLength"`
		AwardRuleTitleLength      int   `json:"awardRuleTitleLength"`
		AwardRuleScoreMax         int   `json:"awardRuleScoreMax"`
		DreamContentLength        int   `json:"dreamContentLength"`
		GiftTitleLength           int   `json:"giftTitleLength"`
		FoodTitleLength           int   `json:"foodTitleLength"`
		FoodContentLength         int   `json:"foodContentLength"`
		TravelTitleLength         int   `json:"travelTitleLength"`
		TravelPlaceContentLength  int   `json:"travelPlaceContentLength"`
		AngryContentLength        int   `json:"angryContentLength"`
		PromiseContentLength      int   `json:"promiseContentLength"`
		PromiseBreakContentLength int   `json:"promiseBreakContentLength"`
		MovieTitleLength          int   `json:"movieTitleLength"`
		MovieContentLength        int   `json:"movieContentLength"`
		// topic
		PostTitleLength              int `json:"postTitleLength"`
		PostContentLength            int `json:"postContentLength"`
		PostScreenReportCount        int `json:"postScreenReportCount"`
		PostCommentContentLength     int `json:"postCommentContentLength"`
		PostCommentScreenReportCount int `json:"postCommentScreenReportCount"`
		// more
		PayVipGoods1Title          string  `json:"payVipGoods1Title"`
		PayVipGoods1Days           int     `json:"payVipGoods1Days"`
		PayVipGoods1Amount         float64 `json:"payVipGoods1Amount"`
		PayVipGoods2Title          string  `json:"payVipGoods2Title"`
		PayVipGoods2Days           int     `json:"payVipGoods2Days"`
		PayVipGoods2Amount         float64 `json:"payVipGoods2Amount"`
		PayVipGoods3Title          string  `json:"payVipGoods3Title"`
		PayVipGoods3Days           int     `json:"payVipGoods3Days"`
		PayVipGoods3Amount         float64 `json:"payVipGoods3Amount"`
		PayCoinGoods1Title         string  `json:"payCoinGoods1Title"`
		PayCoinGoods1Count         int     `json:"payCoinGoods1Count"`
		PayCoinGoods1Amount        float64 `json:"payCoinGoods1Amount"`
		PayCoinGoods2Title         string  `json:"payCoinGoods2Title"`
		PayCoinGoods2Count         int     `json:"payCoinGoods2Count"`
		PayCoinGoods2Amount        float64 `json:"payCoinGoods2Amount"`
		PayCoinGoods3Title         string  `json:"payCoinGoods3Title"`
		PayCoinGoods3Count         int     `json:"payCoinGoods3Count"`
		PayCoinGoods3Amount        float64 `json:"payCoinGoods3Amount"`
		CoinSignMinCount           int     `json:"coinSignMinCount"`
		CoinSignMaxCount           int     `json:"coinSignMaxCount"`
		CoinSignIncreaseCount      int     `json:"coinSignIncreaseCount"`
		CoinAdBetweenSec           int     `json:"coinAdBetweenSec"`
		CoinAdWatchCount           int     `json:"coinAdWatchCount"`
		CoinAdClickCount           int     `json:"coinAdClickCount"`
		CoinAdMaxPerDayCount       int     `json:"coinAdMaxPerDayCount"`
		CoinWishPerDayCount        int     `json:"coinWishPerDayCount"`
		CoinCardPerDayCount        int     `json:"coinCardPerDayCount"`
		MatchWorkScreenReportCount int     `json:"matchWorkScreenReportCount"`
		MatchWorkTitleLength       int     `json:"matchWorkTitleLength"`
		MatchWorkContentLength     int     `json:"matchWorkContentLength"`
		FeatureWishMinDay          int     `json:"featureWishMinDay"`
		FeatureWishMaxDay          int     `json:"featureWishMaxDay"`
		FeatureCardMinDay          int     `json:"featureCardMinDay"`
		FeatureCardMaxDay          int     `json:"featureCardMaxDay"`
	}
	// 会员权限
	VipLimit struct {
		// common
		AdvertiseHide bool `json:"advertiseHide"`
		// couple
		WallPaperSize  int64 `json:"wallPaperSize"`
		WallPaperCount int   `json:"wallPaperCount"`
		// note
		NoteTotalEnable    bool  `json:"noteTotalEnable"`
		SouvenirCount      int   `json:"souvenirCount"`
		VideoSize          int64 `json:"videoSize"`
		AudioSize          int64 `json:"audioSize"`
		PictureSize        int64 `json:"pictureSize"`
		PictureOriginal    bool  `json:"pictureOriginal"`
		PictureTotalCount  int   `json:"pictureTotalCount"`
		DiaryImageSize     int64 `json:"diaryImageSize"`
		DiaryImageCount    int   `json:"diaryImageCount"`
		WhisperImageEnable bool  `json:"whisperImageEnable"`
		GiftImageCount     int   `json:"giftImageCount"`
		FoodImageCount     int   `json:"foodImageCount"`
		MovieImageCount    int   `json:"movieImageCount"`
		// topic
		TopicPostImageCount int `json:"topicPostImageCount"`
	}
	PageSizeLimit struct {
		Sms            int
		User           int
		Entry          int
		Api            int
		Version        int
		Notice         int
		Suggest        int
		SuggestComment int
		Couple         int
		Place          int
		Trends         int
		Souvenir       int
		Menses         int
		Shy            int
		Sleep          int
		Whisper        int
		Word           int
		Diary          int
		Album          int
		Picture        int
		Audio          int
		Video          int
		Food           int
		Travel         int
		Gift           int
		Promise        int
		PromiseBreak   int
		Angry          int
		Dream          int
		Award          int
		AwardRule      int
		Movie          int
		TopicMessage   int
		Post           int
		PostComment    int
		Broadcast      int
		Bill           int
		Vip            int
		Coin           int
		Sign           int
		MatchPeriod    int
		MatchWork      int
	}
)

var (
	limit         *Limit
	pageSizeLimit *PageSizeLimit
)

func GetLimit() *Limit {
	if limit != nil {
		return limit
	}
	limit = &Limit{
		SmsCodeLength:                utils.GetConfigInt("conf", "third.conf", "sms", "code_length"),
		SmsEffectSec:                 utils.GetConfigInt("conf", "third.conf", "sms", "hold_effect_min") * 60,
		SmsBetweenSec:                utils.GetConfigInt("conf", "third.conf", "sms", "send_between_min") * 60,
		SmsMaxSec:                    utils.GetConfigInt("conf", "third.conf", "sms", "send_max_min") * 60,
		SmsMaxCount:                  utils.GetConfigInt("conf", "third.conf", "sms", "send_max_count"),
		SuggestTitleLength:           utils.GetConfigInt("conf", "limit.conf", "input", "suggest_title_length"),
		SuggestContentLength:         utils.GetConfigInt("conf", "limit.conf", "input", "suggest_content_length"),
		SuggestCommentContentLength:  utils.GetConfigInt("conf", "limit.conf", "input", "suggest_comment_content_length"),
		CoupleInviteIntervalSec:      utils.GetConfigInt64("conf", "limit.conf", "time", "couple_invite_interval_min") * 60,
		CoupleBreakNeedSec:           utils.GetConfigInt64("conf", "limit.conf", "time", "couple_break_need_day") * 24 * 60 * 60,
		CoupleBreakSec:               utils.GetConfigInt64("conf", "limit.conf", "time", "couple_break_hou") * 60 * 60,
		CoupleNameLength:             utils.GetConfigInt("conf", "limit.conf", "input", "couple_name_length"),
		NoteResExpireSec:             utils.GetConfigInt64("conf", "limit.conf", "time", "note_res_expire_day") * 24 * 60 * 60,
		MensesMaxPerMonth:            utils.GetConfigInt("conf", "limit.conf", "input", "menses_max_per_month"),
		MensesMaxCycleDay:            utils.GetConfigInt("conf", "limit.conf", "input", "menses_max_cycle_day"),
		MensesMaxDurationDay:         utils.GetConfigInt("conf", "limit.conf", "input", "menses_max_duration_day"),
		MensesDefaultCycleDay:        utils.GetConfigInt("conf", "limit.conf", "input", "menses_default_cycle_day"),
		MensesDefaultDurationDay:     utils.GetConfigInt("conf", "limit.conf", "input", "menses_default_duration_day"),
		ShyMaxPerDay:                 utils.GetConfigInt("conf", "limit.conf", "input", "shy_max_per_day"),
		ShySafeLength:                utils.GetConfigInt("conf", "limit.conf", "input", "shy_safe_length"),
		ShyDescLength:                utils.GetConfigInt("conf", "limit.conf", "input", "shy_desc_length"),
		SleepMaxPerDay:               utils.GetConfigInt("conf", "limit.conf", "input", "sleep_max_per_day"),
		NoteSleepSuccessMinSec:       utils.GetConfigInt64("conf", "limit.conf", "time", "note_sleep_success_min_min") * 60,
		NoteSleepSuccessMaxSec:       utils.GetConfigInt64("conf", "limit.conf", "time", "note_sleep_success_max_hour") * 60 * 60,
		NoteLockLength:               utils.GetConfigInt("conf", "limit.conf", "input", "note_lock_length"),
		SouvenirTitleLength:          utils.GetConfigInt("conf", "limit.conf", "input", "souvenir_title_length"),
		SouvenirForeignYearCount:     utils.GetConfigInt("conf", "limit.conf", "input", "souvenir_foreign_year_count"),
		TravelPlaceCount:             utils.GetConfigInt("conf", "limit.conf", "input", "travel_place_count"),
		TravelVideoCount:             utils.GetConfigInt("conf", "limit.conf", "input", "travel_video_count"),
		TravelFoodCount:              utils.GetConfigInt("conf", "limit.conf", "input", "travel_food_count"),
		TravelMovieCount:             utils.GetConfigInt("conf", "limit.conf", "input", "travel_movie_count"),
		TravelAlbumCount:             utils.GetConfigInt("conf", "limit.conf", "input", "travel_album_count"),
		TravelDiaryCount:             utils.GetConfigInt("conf", "limit.conf", "input", "travel_diary_count"),
		AudioTitleLength:             utils.GetConfigInt("conf", "limit.conf", "input", "audio_title_length"),
		VideoTitleLength:             utils.GetConfigInt("conf", "limit.conf", "input", "video_title_length"),
		AlbumTitleLength:             utils.GetConfigInt("conf", "limit.conf", "input", "album_title_length"),
		PicturePushCount:             utils.GetConfigInt("conf", "limit.conf", "input", "picture_push_count"),
		WordContentLength:            utils.GetConfigInt("conf", "limit.conf", "input", "word_content_length"),
		WhisperChannelLength:         utils.GetConfigInt("conf", "limit.conf", "input", "whisper_channel_length"),
		WhisperContentLength:         utils.GetConfigInt("conf", "limit.conf", "input", "whisper_content_length"),
		DiaryContentLength:           utils.GetConfigInt("conf", "limit.conf", "input", "diary_content_length"),
		AwardContentLength:           utils.GetConfigInt("conf", "limit.conf", "input", "award_content_length"),
		AwardRuleTitleLength:         utils.GetConfigInt("conf", "limit.conf", "input", "award_rule_title_length"),
		AwardRuleScoreMax:            utils.GetConfigInt("conf", "limit.conf", "input", "award_rule_score_max"),
		DreamContentLength:           utils.GetConfigInt("conf", "limit.conf", "input", "dream_content_length"),
		GiftTitleLength:              utils.GetConfigInt("conf", "limit.conf", "input", "gift_title_length"),
		FoodTitleLength:              utils.GetConfigInt("conf", "limit.conf", "input", "food_title_length"),
		FoodContentLength:            utils.GetConfigInt("conf", "limit.conf", "input", "food_content_length"),
		TravelTitleLength:            utils.GetConfigInt("conf", "limit.conf", "input", "travel_title_length"),
		TravelPlaceContentLength:     utils.GetConfigInt("conf", "limit.conf", "input", "travel_place_content_length"),
		AngryContentLength:           utils.GetConfigInt("conf", "limit.conf", "input", "angry_content_length"),
		PromiseContentLength:         utils.GetConfigInt("conf", "limit.conf", "input", "promise_content_length"),
		PromiseBreakContentLength:    utils.GetConfigInt("conf", "limit.conf", "input", "promise_break_content_length"),
		MovieTitleLength:             utils.GetConfigInt("conf", "limit.conf", "input", "movie_title_length"),
		MovieContentLength:           utils.GetConfigInt("conf", "limit.conf", "input", "movie_content_length"),
		PostTitleLength:              utils.GetConfigInt("conf", "limit.conf", "input", "post_title_length"),
		PostContentLength:            utils.GetConfigInt("conf", "limit.conf", "input", "post_content_length"),
		PostScreenReportCount:        utils.GetConfigInt("conf", "limit.conf", "input", "post_screen_report_count"),
		PostCommentContentLength:     utils.GetConfigInt("conf", "limit.conf", "input", "post_comment_content_length"),
		PostCommentScreenReportCount: utils.GetConfigInt("conf", "limit.conf", "input", "post_comment_screen_report_count"),
		PayVipGoods1Title:            utils.GetConfigStr("conf", "limit.conf", "gold", "pay_vip_goods_1_title"),
		PayVipGoods1Days:             utils.GetConfigInt("conf", "limit.conf", "gold", "pay_vip_goods_1_days"),
		PayVipGoods1Amount:           utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_vip_goods_1_amount"),
		PayVipGoods2Title:            utils.GetConfigStr("conf", "limit.conf", "gold", "pay_vip_goods_2_title"),
		PayVipGoods2Days:             utils.GetConfigInt("conf", "limit.conf", "gold", "pay_vip_goods_2_days"),
		PayVipGoods2Amount:           utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_vip_goods_2_amount"),
		PayVipGoods3Title:            utils.GetConfigStr("conf", "limit.conf", "gold", "pay_vip_goods_3_title"),
		PayVipGoods3Days:             utils.GetConfigInt("conf", "limit.conf", "gold", "pay_vip_goods_3_days"),
		PayVipGoods3Amount:           utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_vip_goods_3_amount"),
		PayCoinGoods1Title:           utils.GetConfigStr("conf", "limit.conf", "gold", "pay_coin_goods_1_title"),
		PayCoinGoods1Count:           utils.GetConfigInt("conf", "limit.conf", "gold", "pay_coin_goods_1_count"),
		PayCoinGoods1Amount:          utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_coin_goods_1_amount"),
		PayCoinGoods2Title:           utils.GetConfigStr("conf", "limit.conf", "gold", "pay_coin_goods_2_title"),
		PayCoinGoods2Count:           utils.GetConfigInt("conf", "limit.conf", "gold", "pay_coin_goods_2_count"),
		PayCoinGoods2Amount:          utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_coin_goods_2_amount"),
		PayCoinGoods3Title:           utils.GetConfigStr("conf", "limit.conf", "gold", "pay_coin_goods_3_title"),
		PayCoinGoods3Count:           utils.GetConfigInt("conf", "limit.conf", "gold", "pay_coin_goods_3_count"),
		PayCoinGoods3Amount:          utils.GetConfigFloat64("conf", "limit.conf", "gold", "pay_coin_goods_3_amount"),
		CoinSignMinCount:             utils.GetConfigInt("conf", "limit.conf", "gold", "coin_sign_min_count"),
		CoinSignMaxCount:             utils.GetConfigInt("conf", "limit.conf", "gold", "coin_sign_max_count"),
		CoinSignIncreaseCount:        utils.GetConfigInt("conf", "limit.conf", "gold", "coin_sign_increase_count"),
		CoinAdBetweenSec:             utils.GetConfigInt("conf", "limit.conf", "gold", "coin_ad_between_min") * 60,
		CoinAdWatchCount:             utils.GetConfigInt("conf", "limit.conf", "gold", "coin_ad_watch_count"),
		CoinAdClickCount:             utils.GetConfigInt("conf", "limit.conf", "gold", "coin_ad_click_count"),
		CoinAdMaxPerDayCount:         utils.GetConfigInt("conf", "limit.conf", "gold", "coin_ad_max_per_day_count"),
		CoinWishPerDayCount:          utils.GetConfigInt("conf", "limit.conf", "gold", "coin_wish_per_day_count"),
		CoinCardPerDayCount:          utils.GetConfigInt("conf", "limit.conf", "gold", "coin_card_per_day_count"),
		MatchWorkScreenReportCount:   utils.GetConfigInt("conf", "limit.conf", "input", "match_work_screen_report_count"),
		MatchWorkTitleLength:         utils.GetConfigInt("conf", "limit.conf", "input", "match_work_title_length"),
		MatchWorkContentLength:       utils.GetConfigInt("conf", "limit.conf", "input", "match_work_content_length"),
		FeatureWishMinDay:            utils.GetConfigInt("conf", "limit.conf", "time", "feature_wish_min_day"),
		FeatureWishMaxDay:            utils.GetConfigInt("conf", "limit.conf", "time", "feature_wish_max_day"),
		FeatureCardMinDay:            utils.GetConfigInt("conf", "limit.conf", "time", "feature_card_min_day"),
		FeatureCardMaxDay:            utils.GetConfigInt("conf", "limit.conf", "time", "feature_card_max_day"),
	}
	return limit
}

// GetVipLimitByCouple 获取couple的会员权限
func GetVipLimitByCouple(cid int64) *VipLimit {
	isVip := cid > 0 && IsVip(cid)
	return GetVipLimit(isVip)
}

// GetVipLimit
func GetVipLimit(vip bool) *VipLimit {
	node := ""
	if vip {
		// 是会员
		node = "vip_yes"
	} else {
		// 不是会员
		node = "vip_no"
	}
	v := &VipLimit{
		AdvertiseHide:       utils.GetConfigBool("conf", "limit.conf", node, "advertise_hide"),
		WallPaperSize:       utils.GetConfigInt64("conf", "limit.conf", node, "wall_paper_mb") * utils.MB,
		WallPaperCount:      utils.GetConfigInt("conf", "limit.conf", node, "wall_paper_count"),
		NoteTotalEnable:     utils.GetConfigBool("conf", "limit.conf", node, "note_total_enable"),
		SouvenirCount:       utils.GetConfigInt("conf", "limit.conf", node, "souvenir_count"),
		VideoSize:           utils.GetConfigInt64("conf", "limit.conf", node, "video_mb") * utils.MB,
		AudioSize:           utils.GetConfigInt64("conf", "limit.conf", node, "audio_mb") * utils.MB,
		PictureSize:         utils.GetConfigInt64("conf", "limit.conf", node, "picture_mb") * utils.MB,
		PictureOriginal:     utils.GetConfigBool("conf", "limit.conf", node, "picture_original"),
		PictureTotalCount:   utils.GetConfigInt("conf", "limit.conf", node, "picture_total_count"),
		DiaryImageSize:      utils.GetConfigInt64("conf", "limit.conf", node, "diary_image_mb") * utils.MB,
		DiaryImageCount:     utils.GetConfigInt("conf", "limit.conf", node, "diary_image_count"),
		WhisperImageEnable:  utils.GetConfigBool("conf", "limit.conf", node, "whisper_image_enable"),
		GiftImageCount:      utils.GetConfigInt("conf", "limit.conf", node, "gift_image_count"),
		FoodImageCount:      utils.GetConfigInt("conf", "limit.conf", node, "food_image_count"),
		MovieImageCount:     utils.GetConfigInt("conf", "limit.conf", node, "movie_image_count"),
		TopicPostImageCount: utils.GetConfigInt("conf", "limit.conf", node, "topic_post_image_count"),
	}
	return v
}

func GetPageSizeLimit() *PageSizeLimit {
	if pageSizeLimit != nil {
		return pageSizeLimit
	}
	pageSizeLimit = &PageSizeLimit{
		Sms:            utils.GetConfigInt("conf", "limit.conf", "page_size", "sms"),
		User:           utils.GetConfigInt("conf", "limit.conf", "page_size", "user"),
		Entry:          utils.GetConfigInt("conf", "limit.conf", "page_size", "entry"),
		Api:            utils.GetConfigInt("conf", "limit.conf", "page_size", "api"),
		Version:        utils.GetConfigInt("conf", "limit.conf", "page_size", "version"),
		Notice:         utils.GetConfigInt("conf", "limit.conf", "page_size", "notice"),
		Suggest:        utils.GetConfigInt("conf", "limit.conf", "page_size", "suggest"),
		SuggestComment: utils.GetConfigInt("conf", "limit.conf", "page_size", "suggest_comment"),
		Couple:         utils.GetConfigInt("conf", "limit.conf", "page_size", "couple"),
		Place:          utils.GetConfigInt("conf", "limit.conf", "page_size", "place"),
		Trends:         utils.GetConfigInt("conf", "limit.conf", "page_size", "trends"),
		Souvenir:       utils.GetConfigInt("conf", "limit.conf", "page_size", "souvenir"),
		Menses:         utils.GetConfigInt("conf", "limit.conf", "page_size", "menses"),
		Shy:            utils.GetConfigInt("conf", "limit.conf", "page_size", "shy"),
		Sleep:          utils.GetConfigInt("conf", "limit.conf", "page_size", "sleep"),
		Whisper:        utils.GetConfigInt("conf", "limit.conf", "page_size", "whisper"),
		Word:           utils.GetConfigInt("conf", "limit.conf", "page_size", "word"),
		Diary:          utils.GetConfigInt("conf", "limit.conf", "page_size", "diary"),
		Album:          utils.GetConfigInt("conf", "limit.conf", "page_size", "album"),
		Picture:        utils.GetConfigInt("conf", "limit.conf", "page_size", "picture"),
		Audio:          utils.GetConfigInt("conf", "limit.conf", "page_size", "audio"),
		Video:          utils.GetConfigInt("conf", "limit.conf", "page_size", "video"),
		Food:           utils.GetConfigInt("conf", "limit.conf", "page_size", "food"),
		Travel:         utils.GetConfigInt("conf", "limit.conf", "page_size", "travel"),
		Gift:           utils.GetConfigInt("conf", "limit.conf", "page_size", "gift"),
		Promise:        utils.GetConfigInt("conf", "limit.conf", "page_size", "promise"),
		PromiseBreak:   utils.GetConfigInt("conf", "limit.conf", "page_size", "promise_break"),
		Angry:          utils.GetConfigInt("conf", "limit.conf", "page_size", "angry"),
		Dream:          utils.GetConfigInt("conf", "limit.conf", "page_size", "dream"),
		Award:          utils.GetConfigInt("conf", "limit.conf", "page_size", "award"),
		AwardRule:      utils.GetConfigInt("conf", "limit.conf", "page_size", "award_rule"),
		Movie:          utils.GetConfigInt("conf", "limit.conf", "page_size", "movie"),
		TopicMessage:   utils.GetConfigInt("conf", "limit.conf", "page_size", "topic_message"),
		Post:           utils.GetConfigInt("conf", "limit.conf", "page_size", "post"),
		PostComment:    utils.GetConfigInt("conf", "limit.conf", "page_size", "post_comment"),
		Broadcast:      utils.GetConfigInt("conf", "limit.conf", "page_size", "broadcast"),
		Bill:           utils.GetConfigInt("conf", "limit.conf", "page_size", "bill"),
		Vip:            utils.GetConfigInt("conf", "limit.conf", "page_size", "vip"),
		Coin:           utils.GetConfigInt("conf", "limit.conf", "page_size", "coin"),
		Sign:           utils.GetConfigInt("conf", "limit.conf", "page_size", "sign"),
		MatchPeriod:    utils.GetConfigInt("conf", "limit.conf", "page_size", "match_period"),
		MatchWork:      utils.GetConfigInt("conf", "limit.conf", "page_size", "match_work"),
	}
	return pageSizeLimit
}
