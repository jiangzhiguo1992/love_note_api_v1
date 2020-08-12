package services

import (
	"models/entity"
)

type (
	// 公共常量
	CommonConst struct {
		CompanyName   string `json:"companyName"`   // 公司名称
		CustomerQQ    string `json:"customerQQ"`    // 客服QQ
		OfficialGroup string `json:"officialGroup"` // 官方部落
		OfficialWeibo string `json:"officialWeibo"` // 官方微博
		OfficialWeb   string `json:"officialWeb"`   // 官方网站
		ContactEmail  string `json:"contactEmail"`  // 联系邮箱
		IosAppId      string `json:"iosAppId"`      // appID
	}
	// 常用数量
	CommonCount struct {
		NoticeNewCount     int `json:"noticeNewCount"`
		VersionNewCount    int `json:"versionNewCount"`
		NoteTrendsNewCount int `json:"noteTrendsNewCount"`
		TopicMsgNewCount   int `json:"topicMsgNewCount"`
	}
	// 板块显示
	ModelShow struct {
		MarketPay     bool `json:"marketPay"`
		MarketCoinAd  bool `json:"marketCoinAd"`
		Couple        bool `json:"couple"`
		CouplePlace   bool `json:"couplePlace"`
		CoupleWeather bool `json:"coupleWeather"`
		Note          bool `json:"note"`
		Topic         bool `json:"topic"`
		More          bool `json:"more"`
		MoreVip       bool `json:"moreVip"`
		MoreCoin      bool `json:"moreCoin"`
		MoreMatch     bool `json:"moreMatch"`
		MoreFeature   bool `json:"moreFeature"`
	}
	// 商务合作
	Cooperation struct {
		CatchBabyEnable bool   `json:"catchBabyEnable"`
		CatchBabyUrl    string `json:"catchBabyUrl"`
		CatchBabyDesc   string `json:"catchBabyDesc"`
	}
	// oss
	OssInfo struct {
		// sts
		AccessKeyId     string `json:"accessKeyId"`
		AccessKeySecret string `json:"accessKeySecret"`
		SecurityToken   string `json:"securityToken"`
		StsExpireTime   int64  `json:"stsExpireTime"`
		OssRefreshSec   int64  `json:"ossRefreshSec"`
		UrlExpireSec    int64  `json:"urlExpireSec"`
		// oss
		Region string `json:"region"`
		Domain string `json:"domain"`
		Bucket string `json:"bucket"`
		// path
		PathVersion        string `json:"pathVersion"`
		PathNotice         string `json:"pathNotice"`
		PathBroadcast      string `json:"pathBroadcast"`
		PathLog            string `json:"pathLog"`
		PathSuggest        string `json:"pathSuggest"`
		PathCoupleAvatar   string `json:"pathCoupleAvatar"`
		PathCoupleWall     string `json:"pathCoupleWall"`
		PathNoteVideo      string `json:"pathNoteVideo"`
		PathNoteVideoThumb string `json:"pathNoteVideoThumb"`
		PathNoteAudio      string `json:"pathNoteAudio"`
		PathNoteAlbum      string `json:"pathNoteAlbum"`
		PathNotePicture    string `json:"pathNotePicture"`
		PathNoteWhisper    string `json:"pathNoteWhisper"`
		PathNoteDiary      string `json:"pathNoteDiary"`
		PathNoteGift       string `json:"pathNoteGift"`
		PathNoteFood       string `json:"pathNoteFood"`
		PathNoteMovie      string `json:"pathNoteMovie"`
		PathTopicPost      string `json:"pathTopicPost"`
		PathMoreMatch      string `json:"pathMoreMatch"`
	}
	// 推送
	PushInfo struct {
		AliAppKey     string `json:"aliAppKey"`
		AliAppSecret  string `json:"aliAppSecret"`
		MiAppId       string `json:"miAppId"`
		MiAppKey      string `json:"miAppKey"`
		OppoAppKey    string `json:"oppoAppKey"`
		OppoAppSecret string `json:"oppoAppSecret"`
		ChannelId     string `json:"channelId"`
		NoticeLight   bool   `json:"noticeLight"`
		NoticeSound   bool   `json:"noticeSound"`
		NoticeVibrate bool   `json:"noticeVibrate"`
		NoStartHour   int    `json:"noStartHour"`
		NoEndHour     int    `json:"noEndHour"`
	}
	Push struct {
		CreateAt    int64  `json:"createAt"`    // 时间
		UserId      int64  `json:"userId"`      // 发送者
		ToUserId    int64  `json:"toUserId"`    // 接受者
		Platform    string `json:"platform"`    // 系统
		Title       string `json:"title"`       // 标题
		ContentText string `json:"contentText"` // 内容
		ContentType int    `json:"contentType"` // 类型
		ContentId   int64  `json:"contentId"`   // id
	}
	AdInfo struct {
		AppId           string `json:"appId"`
		TopicPostPosId  string `json:"topicPostPosId"`
		TopicPostStart  int    `json:"topicPostStart"`
		TopicPostJump   int    `json:"topicPostJump"`
		CoinFreePosId   string `json:"coinFreePosId"`
		CoinFreeTickSec int    `json:"coinFreeTickSec"`
	}
	// 配对信息
	PairCard struct {
		TaPhone string `json:"taPhone"` // phone
		Title   string `json:"title"`   // 标题
		Desc    string `json:"message"` // 描述
		BtnBad  string `json:"btnBad"`  // btn字数不要超过6个字
		BtnGood string `json:"btnGood"` // btn字数不要超过6个字
	}
	// 今日天气
	WeatherToday struct {
		Condition string `json:"condition"` // 状况
		Icon      string `json:"icon"`      // 图标
		Temp      string `json:"temp"`      // 温度
		Humidity  string `json:"humidity"`  // 湿度
		WindLevel string `json:"windLevel"` // 风级
		WindDir   string `json:"windDir"`   // 风向
		UpdateAt  int64  `json:"updateAt"`  // 更新时间
	}
	// 天气预报
	WeatherForecast struct {
		TimeAt         int64  `json:"timeAt"`         // 天气时间
		ConditionDay   string `json:"conditionDay"`   // 白天-状况
		ConditionNight string `json:"conditionNight"` // 夜晚-状况
		IconDay        string `json:"iconDay"`        // 白天-图标
		IconNight      string `json:"iconNight"`      // 夜晚-图标
		TempDay        string `json:"tempDay"`        // 白天-温度
		TempNight      string `json:"tempNight"`      // 夜晚-温度
		WindDay        string `json:"windDay"`        // 白天-风向+风级
		WindNight      string `json:"windNight"`      // 夜晚-风向+风级
		UpdateAt       int64  `json:"updateAt"`       // 更新时间
	}
	// 天气预报信息
	WeatherForecastInfo struct {
		Show                string             `json:"show"`
		WeatherForecastList []*WeatherForecast `json:"weatherForecastList"`
	}
	// 姨妈信息
	MensesInfo struct {
		CanMe          bool                 `json:"canMe"`
		CanTa          bool                 `json:"canTa"`
		MensesLengthMe *entity.MensesLength `json:"mensesLengthMe"`
		MensesLengthTa *entity.MensesLength `json:"mensesLengthTa"`
	}
	// 记录统计
	NoteTotal struct {
		TotalSouvenir int64 `json:"totalSouvenir"`
		TotalWord     int64 `json:"totalWord"`
		TotalDiary    int64 `json:"totalDiary"`
		TotalAlbum    int64 `json:"totalAlbum"`
		TotalPicture  int64 `json:"totalPicture"`
		TotalAudio    int64 `json:"totalAudio"`
		TotalVideo    int64 `json:"totalVideo"`
		TotalFood     int64 `json:"totalFood"`
		TotalTravel   int64 `json:"totalTravel"`
		TotalGift     int64 `json:"totalGift"`
		TotalPromise  int64 `json:"totalPromise"`
		TotalAngry    int64 `json:"totalAngry"`
		TotalDream    int64 `json:"totalDream"`
		TotalAward    int64 `json:"totalAward"`
		TotalMovie    int64 `json:"totalMovie"`
	}
	// 帖子大类
	PostKindInfo struct {
		Kind                int                `json:"kind"`
		Enable              bool               `json:"enable"`
		Name                string             `json:"name"`
		PostSubKindInfoList []*PostSubKindInfo `json:"postSubKindInfoList"`
		TopicInfo           *entity.TopicInfo  `json:"topicInfo"`
	}
	// 帖子小类
	PostSubKindInfo struct {
		Kind      int    `json:"kind"`
		Enable    bool   `json:"enable"`
		Name      string `json:"name"`
		Push      bool   `json:"push"`
		Anonymous bool   `json:"anonymous"`
	}
	// 订单信息
	OrderBefore struct {
		Platform   int         `json:"platform"`   // 支付平台
		AliOrder   string      `json:"aliOrder"`   // 阿里订单
		WXOrder    *WXOrder    `json:"wxOrder"`    // 微信订单
		AppleOrder *AppleOrder `json:"appleOrder"` // 苹果订单
	}
	// 微信订单
	WXOrder struct {
		AppId        string `json:"appId"`        // 应用id
		PartnerId    string `json:"partnerId"`    // 商户号
		PrepayId     string `json:"prepayId"`     // 预支付
		PackageValue string `json:"packageValue"` // 扩展字段
		NonceStr     string `json:"nonceStr"`     // 随机字符串
		TimeStamp    string `json:"timeStamp"`    // 时间戳
		Sign         string `json:"sign"`         // 签名
	}
	// 苹果订单
	AppleOrder struct {
		ProductId     string `json:"productId"`     // 商品id
		TransactionId string `json:"transactionId"` // 订单id
		Receipt       string `json:"receipt"`       // 收据
	}
	// 商品
	Goods struct {
		Type   int     // 货物类型
		Title  string  // 货物名称
		Amount float64 // 货物价格
	}
)

const (
	OSS_PREDIX = "acs:oss:*:*:"
	// 公共写，不能读
	OSS_PATH_LOG = "app-log/"
	// 公共读，不能写
	OSS_PATH_VERSION   = "app-version/"
	OSS_PATH_NOTICE    = "app-notice/"
	OSS_PATH_BROADCAST = "app-broadcast/"
	// 公共读，私有写
	OSS_PATH_SUGGEST       = "app-suggest/"
	OSS_PATH_SUGGEST_UID   = OSS_PATH_SUGGEST + "uid-?/"
	OSS_PATH_COUPLE        = "app-couple/"
	OSS_PATH_COUPLE_CID    = OSS_PATH_COUPLE + "cid-?/"
	OSS_PATH_COUPLE_AVATAR = OSS_PATH_COUPLE_CID + "avatar/"
	OSS_PATH_COUPLE_WALL   = OSS_PATH_COUPLE_CID + "wall/"
	OSS_PATH_TOPIC         = "app-topic/"
	OSS_PATH_TOPIC_CID     = OSS_PATH_TOPIC + "cid-?/"
	OSS_PATH_TOPIC_POST    = OSS_PATH_TOPIC_CID + "post/"
	OSS_PATH_MORE          = "app-more/"
	OSS_PATH_MORE_CID      = OSS_PATH_MORE + "cid-?/"
	OSS_PATH_MORE_MATCH    = OSS_PATH_MORE_CID + "match/"
	// 私有读写
	OSS_PATH_NOTE             = "app-note/"
	OSS_PATH_NOTE_CID         = OSS_PATH_NOTE + "cid-?/"
	OSS_PATH_NOTE_AUDIO       = OSS_PATH_NOTE_CID + "audio/"
	OSS_PATH_NOTE_VIDEO       = OSS_PATH_NOTE_CID + "video/"
	OSS_PATH_NOTE_VIDEO_THUMB = OSS_PATH_NOTE_CID + "video-thumb/"
	OSS_PATH_NOTE_ALBUM       = OSS_PATH_NOTE_CID + "album/"
	OSS_PATH_NOTE_PICTURE     = OSS_PATH_NOTE_CID + "picture/"
	OSS_PATH_NOTE_WHISPER     = OSS_PATH_NOTE_CID + "whisper/"
	OSS_PATH_NOTE_DIARY       = OSS_PATH_NOTE_CID + "diary/"
	OSS_PATH_NOTE_GIFT        = OSS_PATH_NOTE_CID + "gift/"
	OSS_PATH_NOTE_FOOD        = OSS_PATH_NOTE_CID + "food/"
	OSS_PATH_NOTE_MOVIE       = OSS_PATH_NOTE_CID + "movie/"
)

const (
	// user手机号长度
	PHONE_LENGTH = 11
	// user 修改类型
	USER_TYPE_FORGET_PWD          = 11
	USER_TYPE_UPDATE_PWD          = 12
	USER_TYPE_UPDATE_PHONE        = 13
	USER_TYPE_UPDATE_INFO         = 14
	USER_TYPE_ADMIN_UPDATE_INFO   = 21
	USER_TYPE_ADMIN_UPDATE_STATUS = 22
	// user 登录类型
	USER_TYPE_LOG_PWD = 1
	USER_TYPE_LOG_VER = 2

	// couple put参数
	COUPLE_UPDATE_GOOD = 1 // 更好
	COUPLE_UPDATE_BAD  = 2 // 更坏
	COUPLE_UPDATE_INFO = 3 // 信息
	// couple list参数
	LIST_WHO_BY_CP = 0
	LIST_WHO_BY_ME = 1
	LIST_WHO_BY_TA = 2

	// weather time格式化
	WEATHER_TIME_FORMAT_ALL = "2006-01-02 15:04:05"
	WEATHER_TIME_FORMAT_DAY = "2006-01-02"

	// postComment order
	POST_COMMENT_ORDER_POINT  = 0
	POST_COMMENT_ORDER_CREATE = 1

	// matchWork order
	MATCH_WORK_ORDER_COIN   = 0
	MATCH_WORK_ORDER_POINT  = 1
	MATCH_WORK_ORDER_CREATE = 2
)

var (
	// PostCommentOrderList 0.point 1.create
	PostCommentOrderList = []string{"point_count DESC", "create_at DESC"}
	// MatchWorkOrderList 0.coin_count 1.point_count 2.create_at
	MatchWorkOrderList = []string{"coin_count DESC", "point_count DESC", "create_at DESC"}
)
