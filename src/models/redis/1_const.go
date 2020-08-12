package redis

import (
	"libs/utils"
	"models/entity"
)

const (
	// user
	KEY_ID_USER    = "id->user<-"
	KEY_PHONE_USER = "phone->user<-"
	KEY_TOKEN_USER = "token->user<-"
	// couple
	KEY_UID_COUPLE = "uid->couple<-"
	// post
	KEY_ID_POST                          = "id->post<-"
	KEY_KIND_SUB_KIND_ALL_POST_LIST      = "kind:sub_kind->all_post_list<-"
	KEY_KIND_SUB_KIND_WELL_POST_LIST     = "kind:sub_kind->well_post_list<-"
	KEY_KIND_SUB_KIND_OFFICIAL_POST_LIST = "kind:sub_kind->official_post_list<-"
	// postComment
	KEY_ID_POST_COMMENT              = "id->post_comment<-"
	KEY_PID_POINT_POST_COMMENT_LIST  = "pid->point_post_comment_list<-"
	KEY_PID_CREATE_POST_COMMENT_LIST = "pid->create_post_comment_list<-"
	KEY_CID_POINT_POST_COMMENT_LIST  = "cid->point_post_comment_list<-"
	KEY_CID_CREATE_POST_COMMENT_LIST = "cid->create_post_comment_list<-"
	// matchCork
	KEY_ID_WORK                      = "id->work<-"
	KEY_PERIOD_KIND_COIN_WORK_LIST   = "period:kind->coin_work_list<-"
	KEY_PERIOD_KIND_POINT_WORK_LIST  = "period:kind->point_work_list<-"
	KEY_PERIOD_KIND_CREATE_WORK_LIST = "period:kind->create_work_list<-"
)

// getRedisUserExpire
func getRedisUserExpire() int64 {
	redisUserExpireSec := utils.GetConfigInt64("conf", "limit.conf", "time", "redis_user_expire_hour") * 60 * 60
	return redisUserExpireSec
}

// getRedisCoupleExpire
func getRedisCoupleExpire() int64 {
	redisCoupleExpireSec := utils.GetConfigInt64("conf", "limit.conf", "time", "redis_couple_expire_hour") * 60 * 60
	return redisCoupleExpireSec
}

// getRedisPostExpire
func getRedisPostExpire() int64 {
	redisPostExpireSec := utils.GetConfigInt64("conf", "limit.conf", "time", "redis_post_expire_hour") * 60 * 60
	return redisPostExpireSec
}

// getRedisPostCommentExpire
func getRedisPostCommentExpire() int64 {
	redisPostCommentExpireSec := utils.GetConfigInt64("conf", "limit.conf", "time", "redis_post_comment_expire_hour") * 60 * 60
	return redisPostCommentExpireSec
}

// getRedisMatchWorkExpire
func getRedisMatchWorkExpire() int64 {
	redisMatchWorkExpireSec := utils.GetConfigInt64("conf", "limit.conf", "time", "redis_match_work_expire_hour") * 60 * 60
	return redisMatchWorkExpireSec
}

type (
	// 按照 Post.Update 从大到小排序
	postUpdateSlice []*entity.Post
	// redis存的结构体
	postList struct {
		PostList []*entity.Post
	}
	// 按照 PostComment.Point 从大到小排序
	postCommentPointSlice []*entity.PostComment
	// 按照 PostComment.Create 从大到小排序
	postCommentCreateSlice []*entity.PostComment
	// redis存的结构体
	postCommentList struct {
		PostCommentList []*entity.PostComment
	}
	// 按照 Work.Coin 从大到小排序
	workCoinSlice []*entity.MatchWork
	// 按照 Work.Point 从大到小排序
	workPointSlice []*entity.MatchWork
	// 按照 Work.Create 从大到小排序
	workCreateSlice []*entity.MatchWork
	// redis存的结构体
	matchWorkList struct {
		MatchWorkList []*entity.MatchWork
	}
)

// update 排序方法重写
func (p postUpdateSlice) Len() int {
	return len(p)
}

func (p postUpdateSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p postUpdateSlice) Less(i, j int) bool {
	return (p[i].Top && !p[j].Top) || (p[i].UpdateAt > p[j].UpdateAt)
}

// point 排序方法重写
func (pc postCommentPointSlice) Len() int {
	return len(pc)
}

func (pc postCommentPointSlice) Swap(i, j int) {
	pc[i], pc[j] = pc[j], pc[i]
}

func (pc postCommentPointSlice) Less(i, j int) bool {
	return (pc[i].Official && !pc[j].Official) || (pc[i].PointCount > pc[j].PointCount)
}

// create 排序方法重写
func (pc postCommentCreateSlice) Len() int {
	return len(pc)
}

func (pc postCommentCreateSlice) Swap(i, j int) {
	pc[i], pc[j] = pc[j], pc[i]
}

func (pc postCommentCreateSlice) Less(i, j int) bool {
	return (pc[i].Official && !pc[j].Official) || (pc[i].CreateAt > pc[j].CreateAt)
}

// coin 排序方法重写
func (w workCoinSlice) Len() int {
	return len(w)
}

func (w workCoinSlice) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w workCoinSlice) Less(i, j int) bool {
	return w[i].CoinCount > w[j].CoinCount
}

// point 排序方法重写
func (w workPointSlice) Len() int {
	return len(w)
}

func (w workPointSlice) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w workPointSlice) Less(i, j int) bool {
	return w[i].PointCount > w[j].PointCount
}

// create 排序方法重写
func (w workCreateSlice) Len() int {
	return len(w)
}

func (w workCreateSlice) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w workCreateSlice) Less(i, j int) bool {
	return w[i].CreateAt > w[j].CreateAt
}
