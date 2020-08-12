package entity

const (
	// statue (app同步) 一般0位正常，-1为删除切不可逆，>0的是用来扩展的
	STATUS_DELETE  = -1 // 不可见
	STATUS_VISIBLE = 0  // 可见，小本本里是仅couple可见

	// from
	APP_FROM_1 = 1 // 鱼笙

	// phone
	PHONE_AREA_CHINA = "86"

	// user-sex (app同步)
	USER_SEX_GIRL = 1
	USER_SEX_BOY  = 2

	// user-birth
	USER_BIRTH_LIGHT = 0
	USER_BIRTH_NIGHT = 1

	// admin
	ADMIN_PERMISSION_NO  = 0
	ADMIN_PERMISSION_ALL = 100

	// sms-type (app同步)
	SMS_TYPE_REGISTER = 10 // 注册
	SMS_TYPE_LOGIN    = 11 // 登录
	SMS_TYPE_FORGET   = 12 // 忘记密码
	SMS_TYPE_PHONE    = 13 // 换手机
	SMS_TYPE_LOCK     = 30 // 密码锁

	// push-type (app同步)
	PUSH_TYPE_APP                = 0
	PUSH_TYPE_SUGGEST            = 50
	PUSH_TYPE_COUPLE             = 100
	PUSH_TYPE_COUPLE_INFO        = 110
	PUSH_TYPE_COUPLE_WALL        = 120
	PUSH_TYPE_COUPLE_PLACE       = 130
	PUSH_TYPE_COUPLE_WEATHER     = 140
	PUSH_TYPE_NOTE               = 200
	PUSH_TYPE_NOTE_LOCK          = 210
	PUSH_TYPE_NOTE_TRENDS        = 211
	PUSH_TYPE_NOTE_TOTAL         = 212
	PUSH_TYPE_NOTE_SHY           = 220
	PUSH_TYPE_NOTE_MENSES        = 221
	PUSH_TYPE_NOTE_SLEEP         = 222
	PUSH_TYPE_NOTE_AUDIO         = 230
	PUSH_TYPE_NOTE_VIDEO         = 231
	PUSH_TYPE_NOTE_ALBUM         = 232
	PUSH_TYPE_NOTE_PICTURE       = 233
	PUSH_TYPE_NOTE_SOUVENIR      = 240
	PUSH_TYPE_NOTE_WISH          = 241
	PUSH_TYPE_NOTE_WORD          = 250
	PUSH_TYPE_NOTE_AWARD         = 251
	PUSH_TYPE_NOTE_AWARD_RULE    = 252
	PUSH_TYPE_NOTE_DIARY         = 253
	PUSH_TYPE_NOTE_DREAM         = 260
	PUSH_TYPE_NOTE_ANGRY         = 261
	PUSH_TYPE_NOTE_GIFT          = 262
	PUSH_TYPE_NOTE_PROMISE       = 263
	PUSH_TYPE_NOTE_PROMISE_BREAK = 264
	PUSH_TYPE_NOTE_TRAVEL        = 270
	PUSH_TYPE_NOTE_MOVIE         = 271
	PUSH_TYPE_NOTE_FOOD          = 272
	PUSH_TYPE_TOPIC              = 300
	PUSH_TYPE_TOPIC_MINE         = 310
	PUSH_TYPE_TOPIC_COLLECT      = 320
	PUSH_TYPE_TOPIC_MESSAGE      = 330
	PUSH_TYPE_TOPIC_POST         = 340
	PUSH_TYPE_TOPIC_COMMENT      = 350
	PUSH_TYPE_MORE               = 400

	// notice-type (app同步)
	NOTICE_TYPE_TEXT  = 0 // 文章
	NOTICE_TYPE_URL   = 1 // 网址
	NOTICE_TYPE_IMAGE = 2 // 图片

	// suggest-status (app同步)
	SUGGEST_STATUS_REPLY_NO    = 0  // 未回复
	SUGGEST_STATUS_REPLY_YES   = 10 // 已回复
	SUGGEST_STATUS_ACCEPT_NO   = 20 // 未接受
	SUGGEST_STATUS_ACCEPT_YES  = 30 // 已接受
	SUGGEST_STATUS_HANGLE_ING  = 40 // 处理中
	SUGGEST_STATUS_HANGLE_OVER = 50 // 处理完
	// suggest-kind (app同步)
	SUGGEST_KIND_ALL      = 0  // 全部
	SUGGEST_KIND_ERROR    = 10 // 程序错误
	SUGGEST_KIND_FUNCTION = 20 // 功能添加
	SUGGEST_KIND_OPTIMISE = 30 // 体验优化
	SUGGEST_KIND_DEBUNK   = 40 // 纯粹吐槽

	// coupleState-state (app同步)
	// 1.wo邀请(1)-wo取消(-2)
	// 2.wo邀请(1)-ta拒绝(-2)
	// 3.wo邀请(1)-ta接受(3)-wo分手(0)-wo取消(3)
	// 4.wo邀请(1)-ta接受(3)-wo分手(0)-ta默认(0)
	// 6.wo邀请(1)-ta接受(3)-wo分手(0)-ta接受(-1)-wo复合(2)-wo取消(-1)
	// 7.wo邀请(1)-ta接受(3)-wo分手(0)-ta接受(-1)-wo复合(2)-ta拒绝(-1)
	// 8.wo邀请(1)-ta接受(3)-wo分手(0)-ta接受(-1)-wo复合(2)-ta接受(3)-循环
	COUPLE_STATE_INVITE        = 0   // 正在邀请(SelfVisible)
	COUPLE_STATE_INVITE_CANCEL = 110 // 邀请者撤回(NoVisible)
	COUPLE_STATE_INVITE_REJECT = 120 // 被邀请者拒绝(NoVisible)
	COUPLE_STATE_BREAK         = 210 // 正在分手(Visible)/已分手(NoVisible)
	COUPLE_STATE_BREAK_ACCEPT  = 220 // 被分手者同意(NoVisible)
	COUPLE_STATE_520           = 520 // 在一起(Visible)

	// trends-actType (app同步)
	TRENDS_ACT_TYPE_INSERT = 1 // 添加
	TRENDS_ACT_TYPE_DELETE = 2 // 删除
	TRENDS_ACT_TYPE_UPDATE = 3 // 修改
	TRENDS_ACT_TYPE_QUERY  = 4 // 查看
	// trends-conType (app同步)
	TRENDS_CON_TYPE_SOUVENIR   = 100 // 纪念日
	TRENDS_CON_TYPE_WISH       = 110 // 愿望清单
	TRENDS_CON_TYPE_SHY        = 200 // 羞羞
	TRENDS_CON_TYPE_MENSES     = 210 // 姨妈
	TRENDS_CON_TYPE_SLEEP      = 220 // 睡眠
	TRENDS_CON_TYPE_AUDIO      = 300 // 音频
	TRENDS_CON_TYPE_VIDEO      = 310 // 视频
	TRENDS_CON_TYPE_ALBUM      = 320 // 相册
	TRENDS_CON_TYPE_WORD       = 500 // 留言
	TRENDS_CON_TYPE_WHISPER    = 510 // 耳语
	TRENDS_CON_TYPE_DIARY      = 520 // 日记
	TRENDS_CON_TYPE_AWARD      = 530 // 打卡
	TRENDS_CON_TYPE_AWARD_RULE = 540 // 约定
	TRENDS_CON_TYPE_DREAM      = 550 // 梦境
	TRENDS_CON_TYPE_FOOD       = 560 // 美食
	TRENDS_CON_TYPE_TRAVEL     = 570 // 游记
	TRENDS_CON_TYPE_GIFT       = 580 // 礼物
	TRENDS_CON_TYPE_PROMISE    = 590 // 承诺
	TRENDS_CON_TYPE_ANGRY      = 600 // 生气
	TRENDS_CON_TYPE_MOVIE      = 610 // 电影
	// trends-id (app同步)
	TRENDS_CON_ID_LIST = 0 // 列表信息

	// topicMessage-kind (app同步)
	TOPIC_MESSAGE_KIND_ALL               = 0
	TOPIC_MESSAGE_KIND_OFFICIAL_TEXT     = 1
	TOPIC_MESSAGE_KIND_JAB_IN_POST       = 10
	TOPIC_MESSAGE_KIND_JAB_IN_COMMENT    = 11
	TOPIC_MESSAGE_KIND_POST_BE_REPORT    = 20
	TOPIC_MESSAGE_KIND_POST_BE_POINT     = 21
	TOPIC_MESSAGE_KIND_POST_BE_COLLECT   = 22
	TOPIC_MESSAGE_KIND_POST_BE_COMMENT   = 23
	TOPIC_MESSAGE_KIND_COMMENT_BE_REPLY  = 30
	TOPIC_MESSAGE_KIND_COMMENT_BE_REPORT = 31
	TOPIC_MESSAGE_KIND_COMMENT_BE_POINT  = 32

	// post-kind
	POST_KIND_OPEN_LIVE     = 110 // 生活
	POST_KIND_OPEN_STAR     = 120 // 星座
	POST_KIND_OPEN_ANIMAL   = 130 // 动物
	POST_KIND_LIMIT_UNKNOWN = 220 // 匿名
	// post-subKind
	POST_SUB_KIND_ALL = 0
	// post-subKind-live
	POST_SUB_KIND_LIVE_KNOW     = 10 // 与TA相识
	POST_SUB_KIND_LIVE_PLACES   = 20 // 异地恋
	POST_SUB_KIND_LIVE_TRAVEL   = 30 // 带TA旅行
	POST_SUB_KIND_LIVE_HOUSE    = 40 // 和TA同居
	POST_SUB_KIND_LIVE_MARRY    = 50 // 结婚记
	POST_SUB_KIND_LIVE_LIFE     = 60 // 小日常
	POST_SUB_KIND_LIVE_FRAGMENT = 70 // 碎碎念
	POST_SUB_KIND_LIVE_HAPPY    = 80 // 吃喝玩乐
	// post-subKind-star
	POST_SUB_KIND_STAR_SHEEP   = 10  // 白羊
	POST_SUB_KIND_STAR_MILK    = 20  // 金牛
	POST_SUB_KIND_STAR_SON     = 30  // 双子
	POST_SUB_KIND_STAR_HUGE    = 40  // 巨蟹
	POST_SUB_KIND_STAR_LION    = 50  // 狮子
	POST_SUB_KIND_STAR_GIRL    = 60  // 处女
	POST_SUB_KIND_STAR_BALANCE = 70  // 天秤
	POST_SUB_KIND_STAR_SKY     = 80  // 天蝎
	POST_SUB_KIND_STAR_HAND    = 90  // 射手
	POST_SUB_KIND_STAR_DEVIL   = 100 // 摩羯
	POST_SUB_KIND_STAR_WATER   = 110 // 水瓶
	POST_SUB_KIND_STAR_FISH    = 120 // 双鱼
	// post-subKind-animal
	POST_SUB_KIND_ANIMAL_CAT    = 10 // 猫系
	POST_SUB_KIND_ANIMAL_DOG    = 20 // 犬系
	POST_SUB_KIND_ANIMAL_WOLF   = 30 // 狼系
	POST_SUB_KIND_ANIMAL_FOX    = 40 // 狐系
	POST_SUB_KIND_ANIMAL_RABBIT = 50 // 兔系
	POST_SUB_KIND_ANIMAL_TIGER  = 60 // 虎系
	POST_SUB_KIND_ANIMAL_BEE    = 70 // 蜂系
	POST_SUB_KIND_ANIMAL_MILK   = 80 // 牛系
	POST_SUB_KIND_ANIMAL_MONKEY = 90 // 猴系
	// post-subKind-group
	POST_SUB_KIND_GROUP_2BOY  = 110 // 男同
	POST_SUB_KIND_GROUP_2GIRL = 120 // 女同
	POST_SUB_KIND_GROUP_BOY   = 210 // 男性
	POST_SUB_KIND_GROUP_GIRL  = 220 // 女性
	POST_SUB_KIND_GROUP_LOVER = 310 // 未婚
	POST_SUB_KIND_GROUP_MARRY = 320 // 已婚
	POST_SUB_KIND_GROUP_STUDY = 410 // 学生
	POST_SUB_KIND_GROUP_WORK  = 420 // 职场
	// post-subKind-unknown
	POST_SUB_KIND_UNKNOWN_GIFT  = 10 // 礼物讨论
	POST_SUB_KIND_UNKNOWN_ROMAN = 20 // 浪漫攻略
	POST_SUB_KIND_UNKNOWN_BRAWL = 30 // 吵架帮忙
	POST_SUB_KIND_UNKNOWN_SHY   = 40 // 羞羞的事
	POST_SUB_KIND_UNKNOWN_TREE  = 50 // 树洞倾诉

	// postComment-kind (app同步)
	POST_COMMENT_KIND_TEXT = 0 // 文本
	POST_COMMENT_KIND_JAB  = 1 // 戳TA

	// broadcast-type (app同步)
	BROADCAST_TYPE_TEXT  = 0 // 文章
	BROADCAST_TYPE_URL   = 1 // 网址
	BROADCAST_TYPE_IMAGE = 2 // 图片

	// bill-payPlatform (app同步)
	BILL_PLATFORM_PAY_ALI    = 100 // 阿里
	BILL_PLATFORM_PAY_WX     = 200 // 微信
	BILL_PLATFORM_PAY_APPLE  = 300 // 苹果
	BILL_PLATFORM_PAY_GOOGLE = 400 // 谷歌
	// bill-payType
	BILL_PAY_TYPE_APP = 1 // app支付
	// bill-goodsType (app同步) 其中价格不能修改，只能新开type
	BILL_AND_GOODS_TYPE_VIP_1  = 1101   // 一月
	BILL_IOS_GOODS_TYPE_VIP_1  = 110100 // 一月
	BILL_AND_GOODS_TYPE_VIP_2  = 1201   // 一年
	BILL_IOS_GOODS_TYPE_VIP_2  = 12010  // 一年
	BILL_AND_GOODS_TYPE_VIP_3  = 1301   // 百年
	BILL_IOS_GOODS_TYPE_VIP_3  = 13010  // 百年
	BILL_AND_GOODS_TYPE_COIN_1 = 2101   // 100
	BILL_IOS_GOODS_TYPE_COIN_1 = 21010  // 100
	BILL_AND_GOODS_TYPE_COIN_2 = 2201   // 1500
	BILL_IOS_GOODS_TYPE_COIN_2 = 22010  // 1500
	BILL_AND_GOODS_TYPE_COIN_3 = 2301   // 20000
	BILL_IOS_GOODS_TYPE_COIN_3 = 23010  // 20000

	// vip-fromType
	VIP_FROM_TYPE_SYS_SEND = 10  // 系统赠送
	VIP_FROM_TYPE_USER_BUY = 100 // 用户购买

	// coin-kind-add (app同步)
	COIN_KIND_ADD_BY_SYS        = 10   // +系统变更
	COIN_KIND_ADD_BY_PLAY_PAY   = 100  // +商店充值
	COIN_KIND_ADD_BY_SIGN_DAY   = 200  // +每日签到
	COIN_KIND_ADD_BY_AD_WATCH   = 210  // +广告观看
	COIN_KIND_ADD_BY_AD_CLICK   = 211  // +广告点击
	COIN_KIND_ADD_BY_MATCH_POST = 300  // +参加比拼
	COIN_KIND_SUB_BY_MATCH_UP   = -300 // -比拼投币
	COIN_KIND_SUB_BY_WISH_UP    = -410 // -许愿投币
	COIN_KIND_SUB_BY_CARD_UP    = -420 // -明信投币

	// match-kind (app同步)
	MATCH_KIND_WIFE_PICTURE = 100 // 照片墙
	MATCH_KIND_LETTER_SHOW  = 200 // 情话集
	MATCH_KIND_DISCUSS_MEET = 300 // 讨论会
)

var (
	// kind(add同步subKind(func)、topicInfo、language，del注释掉下面item)
	PostKindList = []int{
		POST_KIND_OPEN_LIVE,
		POST_KIND_OPEN_STAR,
		POST_KIND_OPEN_ANIMAL,
		POST_KIND_LIMIT_UNKNOWN,
	}
	// subKind(add同步func、topicInfo、language，del注释掉func里的item)
	postSubKindMap = make(map[int]map[int]int, 0)
)
