package entity

type (
	// 非mysql实体类
	FiledInfo struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	BaseObj struct {
		Id       int64 `json:"id"`       // --
		Status   int   `json:"status"`   // 状态，一般见上
		CreateAt int64 `json:"createAt"` // 注册时间，(web生成)
		UpdateAt int64 `json:"updateAt"` // 更新时间，最后一次是删除时间(web生成)
	}
	BaseCp struct {
		UserId   int64 `json:"userId"`   // 用户
		CoupleId int64 `json:"coupleId"` // 配对
	}
	Api struct {
		BaseObj
		UserId   int64   `json:"userId"`
		Platform string  `json:"platform"`
		Language string  `json:"language"`
		URI      string  `json:"uri"`
		Method   string  `json:"method"`
		Params   string  `json:"params"`
		Body     string  `json:"body"`
		Result   string  `json:"result"`
		Duration float64 `json:"duration"`
	}
	User struct {
		BaseObj
		Phone     string `json:"phone"`     // 手机号
		Password  string `json:"password"`  // 密码
		Sex       int8   `json:"sex"`       // 0女生，1男生
		Birthday  int64  `json:"birthday"`  // --
		UserToken string `json:"userToken"` // 用户令牌
		// 关联
		Couple *Couple `json:"couple"`
	}
	Administer struct {
		BaseObj
		UserId     int64 `json:"userId"`
		Permission int   `json:"permission"`
	}
	Sms struct {
		BaseObj
		Phone    string `json:"phone"`    // 手机号
		SendType int    `json:"sendType"` // 类型
		Content  string `json:"content"`  // 内容
	}
	Entry struct {
		BaseObj
		UserId     int64  `json:"userId"`     // --
		DeviceId   string `json:"deviceId"`   // 设备标识，用于身份验证
		DeviceName string `json:"deviceName"` // 设备名称+型号，用于兼容设备
		Market     string `json:"market"`     // 渠道，用于统计下载来源
		Language   string `json:"language"`   // 语言
		Platform   string `json:"platform"`   // WeChat，IOS，Android，用于统计平台
		OsVersion  string `json:"osVersion"`  // weChat/android/ios版本，用于统计兼容版本
		AppVersion int    `json:"appVersion"` // 软件versionCode，用于统计升级率，409 低版本升级
	}
	Version struct {
		BaseObj
		Platform    string `json:"platform"`    // 平台
		VersionName string `json:"versionName"` // 版本名
		VersionCode int    `json:"versionCode"` // 版本标识
		UpdateLog   string `json:"updateLog"`   // 升级日志
		UpdateUrl   string `json:"updateUrl"`   // 升级地址
	}
	Notice struct {
		BaseObj
		Title       string `json:"title"`       // 标题
		ContentType int    `json:"contentType"` // 类型
		ContentText string `json:"contentText"` // 文字
		// 关联
		Read bool `json:"read"` // 已读
	}
	NoticeRead struct {
		BaseObj
		UserId   int64 `json:"userId"`
		NoticeId int64 `json:"noticeId"`
	}
	Suggest struct {
		BaseObj
		UserId       int64  `json:"userId"`       // 用户
		Kind         int    `json:"kind"`         // 类型
		Title        string `json:"title"`        // 标题
		ContentText  string `json:"contentText"`  // 内容文本
		ContentImage string `json:"contentImage"` // 内容图片
		Top          bool   `json:"top"`          // 置顶
		Official     bool   `json:"official"`     // 官方
		FollowCount  int    `json:"followCount"`  // 关注数量
		CommentCount int    `json:"commentCount"` // 评论数量
		// 关联
		Mine    bool `json:"mine"`    // 我的
		Follow  bool `json:"follow"`  // 关注
		Comment bool `json:"comment"` // 评论
	}
	SuggestComment struct {
		BaseObj
		UserId      int64  `json:"userId"`      // 用户
		SuggestId   int64  `json:"suggestId"`   // 建议id
		ContentText string `json:"contentText"` // 内容
		Official    bool   `json:"official"`    // 官方
		// 关联
		Mine bool `json:"mine"` // 我的
	}
	SuggestFollow struct {
		BaseObj
		UserId    int64 `json:"userId"`    // 用户
		SuggestId int64 `json:"suggestId"` // 建议id
	}
	Couple struct {
		BaseObj
		TogetherAt    int64  `json:"togetherAt"`    // 在一起
		CreatorId     int64  `json:"creatorId"`     // 创建者id
		InviteeId     int64  `json:"inviteeId"`     // 受邀者id
		CreatorName   string `json:"creatorName"`   // 创建者昵称
		InviteeName   string `json:"inviteeName"`   // 受邀者昵称
		CreatorAvatar string `json:"creatorAvatar"` // 创建者头像
		InviteeAvatar string `json:"inviteeAvatar"` // 受邀者头像
		// 关联
		State *CoupleState `json:"state"` // cp状态
	}
	CoupleState struct {
		BaseObj
		BaseCp
		State int `json:"state"`
	}
	WallPaper struct {
		BaseObj
		CoupleId      int64  `json:"coupleId"`      // 配对
		ContentImages string `json:"contentImages"` // 图片集合
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
	}
	Place struct {
		BaseObj
		BaseCp
		Longitude float64 `json:"longitude"` // 经度
		Latitude  float64 `json:"latitude"`  // 纬度
		Address   string  `json:"address"`   // 综合地址
		Country   string  `json:"country"`   // 国家
		Province  string  `json:"province"`  // 省
		City      string  `json:"city"`      // 城/市
		District  string  `json:"district"`  // 区
		Street    string  `json:"street"`    // 街道
		CityId    string  `json:"cityId"`    // 城市编号
	}
	Lock struct {
		BaseObj
		BaseCp
		Password string `json:"password"` // 密码
		IsLock   bool   `json:"isLock"`   // 上锁
	}
	Trends struct {
		BaseObj
		BaseCp
		ActionType  int   `json:"actionType"`  // 操作类型
		ContentType int   `json:"contentType"` // 动态类型
		ContentId   int64 `json:"contentId"`   // 动态id
	}
	TrendsBrowse struct {
		BaseObj
		BaseCp
	}
	Shy struct {
		BaseObj
		BaseCp
		Year        int    `json:"year"`        // 年
		MonthOfYear int    `json:"monthOfYear"` // 月
		DayOfMonth  int    `json:"dayOfMonth"`  // 日
		HappenAt    int64  `json:"happenAt"`    // 发生时间
		EndAt       int64  `json:"endAt"`       // 结束时间
		Safe        string `json:"safe"`        // 安全措施
		Desc        string `json:"desc"`        // 描述
	}
	Menses struct {
		BaseObj
		BaseCp
		IsMe        bool `json:"isMe"`        // 我的
		Year        int  `json:"year"`        // 年
		MonthOfYear int  `json:"monthOfYear"` // 月
		DayOfMonth  int  `json:"dayOfMonth"`  // 日
		IsStart     bool `json:"isStart"`     // 是否开始
	}
	MensesLength struct {
		BaseObj
		BaseCp
		CycleDay    int `json:"cycleDay"`    // 周期长度
		DurationDay int `json:"durationDay"` // 经期长度
	}
	Menses2 struct {
		BaseObj
		BaseCp
		StartAt          int64 `json:"startAt"`          // 开始
		EndAt            int64 `json:"endAt"`            // 结束
		StartYear        int   `json:"startYear"`        // 姨妈开始年份
		StartMonthOfYear int   `json:"startMonthOfYear"` // 姨妈开始月份
		StartDayOfMonth  int   `json:"startDayOfMonth"`  // 姨妈开始日期
		EndYear          int   `json:"endYear"`          // 姨妈结束年份
		EndMonthOfYear   int   `json:"endMonthOfYear"`   // 姨妈结束月份
		EndDayOfMonth    int   `json:"endDayOfMonth"`    // 姨妈结束日期
		// 关联
		IsReal                 bool         `json:"isReal"`                 // 开始
		MensesDayList          []*MensesDay `json:"mensesDayList"`          // day集合
		SafeStartYear          int          `json:"safeStartYear"`          // 安全期开始年份
		SafeStartMonthOfYear   int          `json:"safeStartMonthOfYear"`   // 安全期开始月份
		SafeStartDayOfMonth    int          `json:"safeStartDayOfMonth"`    // 安全期开始日期
		SafeEndYear            int          `json:"safeEndYear"`            // 安全期结束年份
		SafeEndMonthOfYear     int          `json:"safeEndMonthOfYear"`     // 安全期结束月份
		SafeEndDayOfMonth      int          `json:"safeEndDayOfMonth"`      // 安全期结束日期
		DangerStartYear        int          `json:"dangerStartYear"`        // 危险期开始年份
		DangerStartMonthOfYear int          `json:"dangerStartMonthOfYear"` // 危险期开始月份
		DangerStartDayOfMonth  int          `json:"dangerStartDayOfMonth"`  // 危险期开始日期
		DangerEndYear          int          `json:"dangerEndYear"`          // 危险期结束年份
		DangerEndMonthOfYear   int          `json:"dangerEndMonthOfYear"`   // 危险期结束月份
		DangerEndDayOfMonth    int          `json:"dangerEndDayOfMonth"`    // 危险期结束日期
		OvulationYear          int          `json:"ovulationYear"`          // 排卵年份
		OvulationMonthOfYear   int          `json:"ovulationMonthOfYear"`   // 排卵月份
		OvulationDayOfMonth    int          `json:"ovulationDayOfMonth"`    // 排卵日期
	}
	MensesDay struct {
		BaseObj
		BaseCp
		Menses2Id   int64 `json:"menses2Id"`   // 姨妈
		Year        int   `json:"year"`        // 年
		MonthOfYear int   `json:"monthOfYear"` // 月
		DayOfMonth  int   `json:"dayOfMonth"`  // 日
		Blood       int   `json:"blood"`       // 血量
		Pain        int   `json:"pain"`        // 痛经
		Mood        int   `json:"mood"`        // 心情
	}
	Sleep struct {
		BaseObj
		BaseCp
		Year        int  `json:"year"`        // 年
		MonthOfYear int  `json:"monthOfYear"` // 月
		DayOfMonth  int  `json:"dayOfMonth"`  // 日
		IsSleep     bool `json:"isSleep"`     // 睡觉/醒来
	}
	Audio struct {
		BaseObj
		BaseCp
		HappenAt     int64  `json:"happenAt"`     // 发生
		Title        string `json:"title"`        // 标题
		ContentAudio string `json:"contentAudio"` // 内容
		Duration     int    `json:"duration"`     // 时长
	}
	Video struct {
		BaseObj
		BaseCp
		HappenAt     int64   `json:"happenAt"`     // 发生
		Title        string  `json:"title"`        // 标题
		ContentThumb string  `json:"contentThumb"` // 缩略图
		ContentVideo string  `json:"contentVideo"` // 内容
		Duration     int     `json:"duration"`     // 时长
		Longitude    float64 `json:"longitude"`    // 经度
		Latitude     float64 `json:"latitude"`     // 纬度
		Address      string  `json:"address"`      // 综合地址
		CityId       string  `json:"cityId"`       // 留作备用
	}
	Album struct {
		BaseObj
		BaseCp
		Title        string `json:"title"`        // 相册名
		Cover        string `json:"cover"`        // 封面
		StartAt      int64  `json:"startAt"`      // 起始时间
		EndAt        int64  `json:"endAt"`        // 终止时间
		PictureCount int    `json:"pictureCount"` // 数量
	}
	Picture struct {
		BaseObj
		BaseCp
		AlbumId      int64   `json:"albumId"`      // 相册id
		HappenAt     int64   `json:"happenAt"`     // 发生
		ContentImage string  `json:"contentImage"` // 照片内容
		Longitude    float64 `json:"longitude"`    // 经度
		Latitude     float64 `json:"latitude"`     // 纬度
		Address      string  `json:"address"`      // 综合地址
		CityId       string  `json:"cityId"`       // 留作备用
	}
	Souvenir struct {
		BaseObj
		BaseCp
		HappenAt  int64   `json:"happenAt"`  // 实现/期望时间
		Title     string  `json:"title"`     // 标题
		Done      bool    `json:"done"`      // 实现/期望
		Longitude float64 `json:"longitude"` // 实现/期望经度
		Latitude  float64 `json:"latitude"`  // 实现/期望纬度
		Address   string  `json:"address"`   // 实现/期望综合地址
		CityId    string  `json:"cityId"`    // 实现/期望留作备用
		// 关联
		SouvenirGiftList   []*SouvenirGift   `json:"souvenirGiftList"`   // 礼物
		SouvenirTravelList []*SouvenirTravel `json:"souvenirTravelList"` // 游记
		SouvenirAlbumList  []*SouvenirAlbum  `json:"souvenirAlbumList"`  // 相册
		SouvenirVideoList  []*SouvenirVideo  `json:"souvenirVideoList"`  // 视频
		SouvenirFoodList   []*SouvenirFood   `json:"souvenirFoodList"`   // 美食
		SouvenirMovieList  []*SouvenirMovie  `json:"souvenirMovieList"`  // 电影
		SouvenirDiaryList  []*SouvenirDiary  `json:"souvenirDiaryList"`  // 日记
	}
	SouvenirTravel struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		TravelId   int64 `json:"travelId"`   // 关联游记
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Travel *Travel `json:"travel"`
	}
	SouvenirVideo struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		VideoId    int64 `json:"videoId"`    // 关联视频
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Video *Video `json:"video"`
	}
	SouvenirAlbum struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		AlbumId    int64 `json:"albumId"`    // 关联相册
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Album *Album `json:"album"`
	}
	SouvenirDiary struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		DiaryId    int64 `json:"diaryId"`    // 关联日记
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Diary *Diary `json:"diary"`
	}
	SouvenirFood struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		FoodId     int64 `json:"foodId"`     // 关联美食
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Food *Food `json:"food"`
	}
	SouvenirMovie struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		MovieId    int64 `json:"movieId"`    // 关联电影
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Movie *Movie `json:"movie"`
	}
	SouvenirGift struct {
		BaseObj
		BaseCp
		SouvenirId int64 `json:"souvenirId"` // 关联纪念日
		GiftId     int64 `json:"giftId"`     // 关联礼物
		Year       int   `json:"year"`       // 哪一年的
		// 关联
		Gift *Gift `json:"gift"`
	}
	Travel struct {
		BaseObj
		BaseCp
		HappenAt int64  `json:"happenAt"` // 时间
		Title    string `json:"title"`    // 标题
		// 关联
		TravelPlaceList []*TravelPlace `json:"travelPlaceList"` // 行程列表
		TravelAlbumList []*TravelAlbum `json:"travelAlbumList"` // 相册
		TravelVideoList []*TravelVideo `json:"travelVideoList"` // 视频
		TravelFoodList  []*TravelFood  `json:"travelFoodList"`  // 美食
		TravelMovieList []*TravelMovie `json:"travelMovieList"` // 电影
		TravelDiaryList []*TravelDiary `json:"travelDiaryList"` // 日记
	}
	TravelPlace struct {
		BaseObj
		BaseCp
		TravelId    int64   `json:"travelId"`    // 关联游记
		HappenAt    int64   `json:"happenAt"`    // 时间
		ContentText string  `json:"contentText"` // 介绍
		Longitude   float64 `json:"longitude"`   // 经度
		Latitude    float64 `json:"latitude"`    // 纬度
		Address     string  `json:"address"`     // 综合地址
		CityId      string  `json:"cityId"`      // 城市编号
	}
	TravelVideo struct {
		BaseObj
		BaseCp
		TravelId int64 `json:"travelId"` // 关联游记
		VideoId  int64 `json:"videoId"`  // 关联日记
		// 关联
		Video *Video `json:"video"`
	}
	TravelDiary struct {
		BaseObj
		BaseCp
		TravelId int64 `json:"travelId"` // 关联游记
		DiaryId  int64 `json:"diaryId"`  // 关联日记
		// 关联
		Diary *Diary `json:"diary"`
	}
	TravelAlbum struct {
		BaseObj
		BaseCp
		TravelId int64 `json:"travelId"` // 关联游记
		AlbumId  int64 `json:"albumId"`  // 关联日记
		// 关联
		Album *Album `json:"album"`
	}
	TravelFood struct {
		BaseObj
		BaseCp
		TravelId int64 `json:"travelId"` // 关联游记
		FoodId   int64 `json:"foodId"`   // 关联美食
		// 关联
		Food *Food `json:"food"`
	}
	TravelMovie struct {
		BaseObj
		BaseCp
		TravelId int64 `json:"travelId"` // 关联游记
		MovieId  int64 `json:"movieId"`  // 关联电影
		// 关联
		Movie *Movie `json:"movie"`
	}
	Word struct {
		BaseObj
		BaseCp
		ContentText string `json:"contentText"` // 内容
	}
	Whisper struct {
		BaseObj
		BaseCp
		Channel string `json:"channel"` // 频道
		IsImage bool   `json:"isImage"` // 是否是图片
		Content string `json:"content"` // 内容
	}
	Diary struct {
		BaseObj
		BaseCp
		HappenAt      int64  `json:"happenAt"`      // 发生
		ContentText   string `json:"contentText"`   // 内容
		ContentImages string `json:"contentImages"` // 图片
		ReadCount     int    `json:"readCount"`     // 阅读数
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
	}
	Award struct {
		BaseObj
		BaseCp
		HappenId    int64  `json:"happenId"`    // 发生者
		AwardRuleId int64  `json:"awardRuleId"` // 约定
		HappenAt    int64  `json:"happenAt"`    // 时间
		ContentText string `json:"contentText"` // 内容
		ScoreChange int    `json:"scoreChange"` // 分数
	}
	AwardRule struct {
		BaseObj
		BaseCp
		Title    string `json:"title"`    // 标题
		Score    int    `json:"score"`    // 分数
		UseCount int    `json:"useCount"` // 引用次数
	}
	AwardScore struct {
		BaseObj
		BaseCp
		ChangeCount int   `json:"changeCount"` // 变动
		TotalScore  int64 `json:"totalScore"`  // 总分
	}
	Dream struct {
		BaseObj
		BaseCp
		HappenAt    int64  `json:"happenAt"`    // 发生
		ContentText string `json:"contentText"` // 内容
	}
	Gift struct {
		BaseObj
		BaseCp
		ReceiveId     int64  `json:"receiveId"`     // 接受者
		HappenAt      int64  `json:"happenAt"`      // 发生
		Title         string `json:"title"`         // 标题
		ContentImages string `json:"contentImages"` // 图片
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
	}
	Food struct {
		BaseObj
		BaseCp
		HappenAt      int64   `json:"happenAt"`      // 发生
		Title         string  `json:"title"`         // 内容
		ContentImages string  `json:"contentImages"` // 图片
		ContentText   string  `json:"contentText"`   // 文字
		Longitude     float64 `json:"longitude"`     // 经度
		Latitude      float64 `json:"latitude"`      // 纬度
		Address       string  `json:"address"`       // 综合地址
		CityId        string  `json:"cityId"`        // 留作备用
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
	}
	Angry struct {
		BaseObj
		BaseCp
		HappenId    int64  `json:"happenId"`    // 生气者
		HappenAt    int64  `json:"happenAt"`    // 发生
		ContentText string `json:"contentText"` // 内容
		GiftId      int64  `json:"giftId"`      // 礼物
		PromiseId   int64  `json:"promiseId"`   // 承诺
		// 关联
		Gift    *Gift    `json:"gift"`
		Promise *Promise `json:"promise"`
	}
	Promise struct {
		BaseObj
		BaseCp
		HappenId    int64  `json:"happenId"`    // 承诺者
		HappenAt    int64  `json:"happenAt"`    // 开始时间
		ContentText string `json:"contentText"` // 内容
		BreakCount  int    `json:"breakCount"`  // 违反次数
	}
	PromiseBreak struct {
		BaseObj
		BaseCp
		PromiseId   int64  `json:"promiseId"`   // 关联承诺
		HappenAt    int64  `json:"happenAt"`    // 违反时间
		ContentText string `json:"contentText"` // 内容
	}
	Movie struct {
		BaseObj
		BaseCp
		HappenAt      int64   `json:"happenAt"`      // 发生
		Title         string  `json:"title"`         // 内容
		ContentImages string  `json:"contentImages"` // 图片
		ContentText   string  `json:"contentText"`   // 文字
		Longitude     float64 `json:"longitude"`     // 经度
		Latitude      float64 `json:"latitude"`      // 纬度
		Address       string  `json:"address"`       // 综合地址
		CityId        string  `json:"cityId"`        // 留作备用
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
	}
	TopicInfo struct {
		BaseObj
		Kind         int `json:"kind"`
		Year         int `json:"year"`
		DayOfYear    int `json:"dayOfYear"`
		PostCount    int `json:"postCount"`
		BrowseCount  int `json:"browseCount"`
		CommentCount int `json:"commentCount"`
		ReportCount  int `json:"reportCount"`
		PointCount   int `json:"pointCount"`
		CollectCount int `json:"collectCount"`
	}
	TopicMessage struct {
		BaseObj
		BaseCp
		ToUserId    int64  `json:"toUserId"`    // 给谁的
		ToCoupleId  int64  `json:"toCoupleId"`  // 给谁们的
		Kind        int    `json:"kind"`        // 类型
		ContentText string `json:"contentText"` // 内容
		ContentId   int64  `json:"contentId"`   // 内容id
		// 关联
		Couple *Couple `json:"couple"` // 目标信息
	}
	TopicMessageBrowse struct {
		BaseObj
		BaseCp
	}
	Post struct {
		BaseObj
		BaseCp
		Kind          int    `json:"kind"`          // 种类
		SubKind       int    `json:"subKind"`       // 子类
		Title         string `json:"title"`         // 标题
		ContentText   string `json:"contentText"`   // 内容文本
		ContentImages string `json:"contentImages"` // 内容图片
		Top           bool   `json:"top"`           // 置顶
		Official      bool   `json:"official"`      // 官方
		Well          bool   `json:"well"`          // 精华
		ReportCount   int    `json:"reportCount"`   // 举报数
		PointCount    int    `json:"pointCount"`    // 点赞数
		CollectCount  int    `json:"collectCount"`  // 关注数
		CommentCount  int    `json:"commentCount"`  // 评论数
		// 关联
		ContentImageList []string `json:"contentImageList"` // 图片集合
		Screen           bool     `json:"screen"`           // 屏蔽
		Hot              bool     `json:"hot"`              // 热门
		Couple           *Couple  `json:"couple"`           // 信息
		Mine             bool     `json:"mine"`             // 我的
		Our              bool     `json:"our"`              // 我们的
		Read             bool     `json:"read"`             // 阅读
		Report           bool     `json:"report"`           // 举报
		Point            bool     `json:"point"`            // 点赞
		Collect          bool     `json:"collect"`          // 关注
		Comment          bool     `json:"comment"`          // 评论
	}
	PostRead struct {
		BaseObj
		UserId int64 `json:"userId"`
		PostId int64 `json:"postId"` // 帖子id
	}
	PostReport struct {
		BaseObj
		BaseCp
		PostId int64 `json:"postId"` // 帖子id
	}
	PostPoint struct {
		BaseObj
		BaseCp
		PostId int64 `json:"postId"` // 帖子id
	}
	PostCollect struct {
		BaseObj
		BaseCp
		PostId int64 `json:"postId"` // 帖子id
	}
	PostComment struct {
		BaseObj
		BaseCp
		PostId          int64  `json:"postId"`          // 帖子id
		ToCommentId     int64  `json:"toCommentId"`     // 被评论id
		Floor           int    `json:"floor"`           // 楼
		Kind            int    `json:"kind"`            // 种类
		ContentText     string `json:"contentText"`     // 内容文本
		Official        bool   `json:"official"`        // 官方
		SubCommentCount int    `json:"subCommentCount"` // 子评论数量
		ReportCount     int    `json:"reportCount"`     // 举报数
		PointCount      int    `json:"pointCount"`      // 点赞数
		// 关联
		Screen     bool    `json:"screen"`     // 屏蔽
		Couple     *Couple `json:"couple"`     // 信息
		Mine       bool    `json:"mine"`       // 我的
		Our        bool    `json:"our"`        // 我们的
		SubComment bool    `json:"subComment"` // 评论
		Report     bool    `json:"report"`     // 举报
		Point      bool    `json:"point"`      // 点赞
	}
	PostCommentReport struct {
		BaseObj
		BaseCp
		PostCommentId int64 `json:"postCommentId"` // 帖子id
	}
	PostCommentPoint struct {
		BaseObj
		BaseCp
		PostCommentId int64 `json:"postCommentId"` // 评论id
	}
	Broadcast struct {
		BaseObj
		Title       string `json:"title"`       // 标题
		Cover       string `json:"cover"`       // 封面
		StartAt     int64  `json:"startAt"`     // 开始
		EndAt       int64  `json:"endAt"`       // 结束
		ContentType int    `json:"contentType"` // 类型
		ContentText string `json:"contentText"` // 文字
		IsEnd       bool   `json:"isEnd"`       // 是否结束
	}
	Bill struct {
		BaseObj
		BaseCp
		PlatformOs   string  `json:"platformOs"`   // 设备平台
		PlatformPay  int     `json:"platformPay"`  // 支付平台
		PayType      int     `json:"payType"`      // 支付方式
		PayAmount    float64 `json:"payAmount"`    // 支付金额
		TradeNo      string  `json:"tradeNo"`      // 订单号
		TradeReceipt string  `json:"tradeReceipt"` // 订单收据
		GoodsType    int     `json:"goodsType"`    // 货物类型
		GoodsId      int64   `json:"goodsId"`      // 货物id
		//TradePay     bool    `json:"tradePay"`     // 订单状态
		//GoodsOut     bool    `json:"goodsOut"`     // 是否发货
	}
	Vip struct {
		BaseObj
		BaseCp
		FromType   int   `json:"fromType"`   // 来源 购买/系统赠送/活动
		ExpireDays int   `json:"expireDays"` // 到期天数
		ExpireAt   int64 `json:"expireAt"`   // 到期时间
		//BillId     int64 `json:"billId"`     // 账单信息
	}
	Coin struct {
		BaseObj
		BaseCp
		Kind   int `json:"kind"`   // 种类
		Change int `json:"change"` // 变化
		Count  int `json:"count"`  // 数目
		//BillId int64 `json:"billId"` // 账单信息
	}
	Sign struct {
		BaseObj
		BaseCp
		Year        int `json:"year"`        // 年
		MonthOfYear int `json:"monthOfYear"` // 月
		DayOfMonth  int `json:"dayOfMonth"`  // 日
		ContinueDay int `json:"continueDay"` // 连续天数
	}
	MatchPeriod struct {
		BaseObj
		StartAt     int64  `json:"startAt"`     // 开始
		EndAt       int64  `json:"endAt"`       // 结束
		Period      int    `json:"period"`      // 期数
		Kind        int    `json:"kind"`        // 种类
		Title       string `json:"title"`       // 标题
		CoinChange  int    `json:"coinChange"`  // 金币变化
		WorksCount  int    `json:"worksCount"`  // 作品数量
		ReportCount int    `json:"reportCount"` // 举报数量
		PointCount  int    `json:"pointCount"`  // 点赞总数
		CoinCount   int    `json:"coinCount"`   // 金币总数
	}
	MatchWork struct {
		BaseObj
		BaseCp
		MatchPeriodId int64  `json:"matchPeriodId"` // 期数
		Kind          int    `json:"kind"`          // 种类
		Title         string `json:"title"`         // 标题
		ContentText   string `json:"contentText"`   // 内容文字
		ContentImage  string `json:"contentImage"`  // 内容图片
		ReportCount   int    `json:"reportCount"`   // 举报数
		PointCount    int    `json:"pointCount"`    // 点赞数
		CoinCount     int    `json:"coinCount"`     // 金币数
		// 关联
		Screen bool    `json:"screen"` // 屏蔽
		Couple *Couple `json:"couple"` // 信息
		Mine   bool    `json:"mine"`   // 我的
		Our    bool    `json:"our"`    // 我们的
		Report bool    `json:"report"` // 举报
		Point  bool    `json:"point"`  // 点赞
		Coin   bool    `json:"coin"`   // 金币
	}
	MatchReport struct {
		BaseObj
		BaseCp
		MatchPeriodId int64  `json:"matchPeriodId"` // 期数id
		MatchWorkId   int64  `json:"matchWorkId"`   // 作品id
		Reason        string `json:"reason"`        // 原因
	}
	MatchPoint struct {
		BaseObj
		BaseCp
		MatchPeriodId int64 `json:"matchPeriodId"` // 期数id
		MatchWorkId   int64 `json:"matchWorkId"`   // 作品id
	}
	MatchCoin struct {
		BaseObj
		BaseCp
		MatchPeriodId int64 `json:"matchPeriodId"` // 期数id
		MatchWorkId   int64 `json:"matchWorkId"`   // 作品id
		CoinId        int64 `json:"coinId"`        // 金币id
		CoinCount     int   `json:"coinCount"`     // 金币数
	}
)
