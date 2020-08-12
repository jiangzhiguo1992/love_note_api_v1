package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"libs/utils"

	"github.com/gomodule/redigo/redis"
	"models/entity"
	"sort"
)

// removeRepeatedPostComment
func removeRepeatedPostComment(list []*entity.PostComment) ([]*entity.PostComment) {
	returnList := make([]*entity.PostComment, 0)
	if list == nil || len(list) <= 0 {
		return returnList
	}
	// 遍历所有元素
	for i := 0; i < len(list); i++ {
		iComment := list[i]
		if iComment == nil || iComment.Id <= 0 {
			continue
		}
		// 检查重复元素 0=10 1=11 2=12 3=11 4=12 5=10 6=12
		repeat := false
		for j := i + 1; j < len(list); j++ {
			jComment := list[j]
			if jComment == nil || jComment.Id <= 0 {
				continue
			}
			if iComment.Id == jComment.Id {
				// 更新较晚的加上，相同的话 加后面的
				if iComment.UpdateAt <= jComment.UpdateAt {
					repeat = true
				} else if iComment.UpdateAt > jComment.UpdateAt {
					repeat = false
				}
				// 防止多重
				if repeat {
					break
				}
			}
		}
		if !repeat {
			// 并且去掉user信息
			iComment.Mine = false
			iComment.Our = false
			iComment.SubComment = false
			iComment.Report = false
			iComment.Point = false
			returnList = append(returnList, iComment)
		}
	}
	return returnList
}

// SetPostComment
func SetPostComment(comment *entity.PostComment) error {
	if comment == nil || comment.Id <= 0 {
		utils.LogWarn("SetPostComment", "无效的评论: "+fmt.Sprintf("%+v", comment))
		return errors.New("nil_comment")
	}
	comment.Mine = false
	comment.Our = false
	comment.SubComment = false
	comment.Report = false
	comment.Point = false
	bytes, err := json.Marshal(comment)
	if err != nil {
		utils.LogErr("SetPostComment", err)
		// 添加失败就删除
		DelPostComment(comment)
		return errors.New("data_decode_err")
	}
	commentBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储comment
	_, err = conn.Do("SET", KEY_ID_POST_COMMENT+strconv.FormatInt(comment.Id, 10), commentBody)
	if err != nil {
		utils.LogErr("SetPostComment", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_ID_POST_COMMENT+strconv.FormatInt(comment.Id, 10), getRedisPostCommentExpire())
	utils.LogErr("SetPostComment", err)
	return err
}

// DelPostComment
func DelPostComment(comment *entity.PostComment) error {
	if comment == nil || comment.Id <= 0 {
		utils.LogWarn("DelPostComment", "无效的评论: "+fmt.Sprintf("%+v", comment))
		return errors.New("nil_comment")
	}
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 开始删除-id
	_, err := conn.Do("DEL", KEY_ID_POST_COMMENT+strconv.FormatInt(comment.Id, 10))
	utils.LogErr("DelPostComment", err)
	return err
}

// GetPostCommentById
func GetPostCommentById(cid int64) (*entity.PostComment, error) {
	if cid <= 0 {
		utils.LogWarn("GetPostCommentById", "cid <= 0")
		return nil, errors.New("nil_comment")
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_ID_POST_COMMENT+strconv.FormatInt(cid, 10))
	if err != nil {
		utils.LogErr("GetPostCommentById", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析comment
	comment := &entity.PostComment{}
	err = json.Unmarshal(bytes, comment)
	if err != nil {
		utils.LogErr("GetPostCommentById", err)
		return nil, err
	}
	return comment, nil
}

//------------------------------------PointList------------------------------------

// DelPostCommentInListByPoint
func DelPostCommentInListByPoint(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("DelPostCommentInListByPoint", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		utils.LogWarn("DelPostCommentInListByPoint", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	commentList, _ := GetPostCommentListByPoint(pc.PostId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.PostComment, 0)
	for _, v := range commentList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != pc.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostCommentListByPoint(pc.PostId, saveList, false)
	return err
}

// UpdatePostCommentInListByPoint
func UpdatePostCommentInListByPoint(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("UpdatePostCommentInListByPoint", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		utils.LogWarn("UpdatePostCommentInListByPoint", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧列表
	commentList, _ := GetPostCommentListByPoint(pc.PostId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 添加新元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostCommentListByPoint(pc.PostId, commentList, false)
	return err
}

// SetPostCommentListByPoint
// merge 是否合并已有的数据
func SetPostCommentListByPoint(pid int64, list []*entity.PostComment, merge bool) error {
	if pid <= 0 {
		utils.LogWarn("SetPostCommentListByPoint", "pid <= 0")
		return errors.New("nil_post")
	} else if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.PostComment, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisCommentList, _ := GetPostCommentListByPoint(pid, -1, -1)
		if redisCommentList != nil && len(redisCommentList) > 0 {
			for _, v := range redisCommentList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedPostComment(returnList)
	// 开始降序排序
	sort.Sort(postCommentPointSlice(returnList))
	// 构造数据
	saveList := &postCommentList{
		PostCommentList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostCommentListByPoint", err)
		// list添加失败不删除
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储comment
	_, err = conn.Do("SET", KEY_PID_POINT_POST_COMMENT_LIST+strconv.FormatInt(pid, 10), listBody)
	if err != nil {
		utils.LogErr("SetPostCommentListByPoint", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_PID_POINT_POST_COMMENT_LIST+strconv.FormatInt(pid, 10), getRedisPostCommentExpire())
	utils.LogErr("SetPostCommentListByPoint", err)
	return nil
}

// GetPostCommentListByPoint
func GetPostCommentListByPoint(pid int64, offset, limit int) ([]*entity.PostComment, error) {
	if pid <= 0 {
		utils.LogWarn("GetPostCommentListByPoint", "pid <= 0")
		return nil, errors.New("nil_post")
	} else if limit == 0 {
		utils.LogWarn("GetPostCommentListByPoint", "limit == 0")
		return nil, nil
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_PID_POINT_POST_COMMENT_LIST+strconv.FormatInt(pid, 10))
	if err != nil {
		utils.LogErr("GetPostCommentListByPoint", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析comment
	list := &postCommentList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostCommentListByPoint", err)
		return nil, err
	}
	commentList := list.PostCommentList
	// 数量不够
	if commentList == nil || len(commentList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return commentList, nil
	}
	// 正常获取，数量够
	returnList := commentList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------CreateList------------------------------------

// AddPostCommentInListByCreate
func AddPostCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("AddPostCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		utils.LogWarn("AddPostCommentInListByCreate", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	commentList, _ := GetPostCommentListByCreate(pc.PostId, -1, -1)
	if commentList == nil {
		commentList = make([]*entity.PostComment, 0)
	}
	// 添加指定元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostCommentListByCreate(pc.PostId, commentList, false)
	return err
}

// DelPostCommentInListByCreate
func DelPostCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("DelPostCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		utils.LogWarn("DelPostCommentInListByCreate", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	commentList, _ := GetPostCommentListByCreate(pc.PostId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.PostComment, 0)
	for _, v := range commentList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != pc.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostCommentListByCreate(pc.PostId, saveList, false)
	return err
}

// UpdatePostCommentInListByCreate
func UpdatePostCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("UpdatePostCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		utils.LogWarn("UpdatePostCommentInListByCreate", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧列表
	commentList, _ := GetPostCommentListByCreate(pc.PostId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 添加新元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostCommentListByCreate(pc.PostId, commentList, false)
	return err
}

// SetPostCommentListByCreate
// merge 是否合并已有的数据
func SetPostCommentListByCreate(pid int64, list []*entity.PostComment, merge bool) error {
	if pid <= 0 {
		utils.LogWarn("SetPostCommentListByCreate", "pid <= 0")
		return errors.New("nil_post")
	} else if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.PostComment, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisCommentList, _ := GetPostCommentListByCreate(pid, -1, -1)
		if redisCommentList != nil && len(redisCommentList) > 0 {
			for _, v := range redisCommentList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedPostComment(returnList)
	// 开始降序排序
	sort.Sort(postCommentCreateSlice(returnList))
	// 构造数据
	saveList := &postCommentList{
		PostCommentList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostCommentListByCreate", err)
		// list添加失败不删除
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储comment
	_, err = conn.Do("SET", KEY_PID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(pid, 10), listBody)
	if err != nil {
		utils.LogErr("SetPostCommentListByCreate", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_PID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(pid, 10), getRedisPostCommentExpire())
	utils.LogErr("SetPostCommentListByCreate", err)
	return nil
}

// RedisGetPostCommentListByCreate
func GetPostCommentListByCreate(pid int64, offset, limit int) ([]*entity.PostComment, error) {
	if pid <= 0 {
		utils.LogWarn("GetPostCommentListByCreate", "pid <= 0")
		return nil, errors.New("nil_post")
	} else if limit == 0 {
		utils.LogWarn("GetPostCommentListByCreate", "limit == 0")
		return nil, nil
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_PID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(pid, 10))
	if err != nil {
		utils.LogErr("GetPostCommentListByCreate", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析comment
	list := &postCommentList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostCommentListByCreate", err)
		return nil, err
	}
	commentList := list.PostCommentList
	// 数量不够
	if commentList == nil || len(commentList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return commentList, nil
	}
	// 正常获取，数量够
	returnList := commentList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------SubPointList------------------------------------

// DelPostToCommentInListByPoint
func DelPostToCommentInListByPoint(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("DelPostToCommentInListByPoint", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.ToCommentId <= 0 {
		utils.LogWarn("DelPostToCommentInListByPoint", "toComment == nil")
		return errors.New("nil_comment")
	}
	// 获取旧数据
	commentList, _ := GetPostToCommentListByPoint(pc.ToCommentId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.PostComment, 0)
	for _, v := range commentList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != pc.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostToCommentListByPoint(pc.ToCommentId, saveList, false)
	return err
}

// UpdatePostToCommentInListByPoint
func UpdatePostToCommentInListByPoint(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("UpdatePostToCommentInListByPoint", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.ToCommentId <= 0 {
		utils.LogWarn("UpdatePostToCommentInListByPoint", "toComment == nil")
		return errors.New("nil_comment")
	}
	// 获取旧列表
	commentList, _ := GetPostToCommentListByPoint(pc.ToCommentId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 添加新元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostToCommentListByPoint(pc.ToCommentId, commentList, false)
	return err
}

// SetPostToCommentListByPoint
// merge 是否合并已有的数据
func SetPostToCommentListByPoint(cid int64, list []*entity.PostComment, merge bool) error {
	if cid <= 0 {
		utils.LogWarn("SetPostToCommentListByPoint", "cid <= 0")
		return errors.New("nil_comment")
	} else if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.PostComment, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisCommentList, _ := GetPostToCommentListByPoint(cid, -1, -1)
		if redisCommentList != nil && len(redisCommentList) > 0 {
			for _, v := range redisCommentList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedPostComment(returnList)
	// 开始降序排序
	sort.Sort(postCommentPointSlice(returnList))
	// 构造数据
	saveList := &postCommentList{
		PostCommentList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostToCommentListByPoint", err)
		// list添加失败不删除
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储comment
	_, err = conn.Do("SET", KEY_CID_POINT_POST_COMMENT_LIST+strconv.FormatInt(cid, 10), listBody)
	if err != nil {
		utils.LogErr("SetPostToCommentListByPoint", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_CID_POINT_POST_COMMENT_LIST+strconv.FormatInt(cid, 10), getRedisPostCommentExpire())
	utils.LogErr("SetPostToCommentListByPoint", err)
	return nil
}

// GetPostToCommentListByPoint
func GetPostToCommentListByPoint(cid int64, offset, limit int) ([]*entity.PostComment, error) {
	if cid <= 0 {
		utils.LogWarn("GetPostToCommentListByPoint", "cid <= 0")
		return nil, errors.New("nil_comment")
	} else if limit == 0 {
		utils.LogWarn("GetPostToCommentListByPoint", "limit == 0")
		return nil, nil
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_CID_POINT_POST_COMMENT_LIST+strconv.FormatInt(cid, 10))
	if err != nil {
		utils.LogErr("GetPostToCommentListByPoint", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析comment
	list := &postCommentList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostToCommentListByPoint", err)
		return nil, err
	}
	commentList := list.PostCommentList
	// 数量不够
	if commentList == nil || len(commentList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return commentList, nil
	}
	// 正常获取，数量够
	returnList := commentList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------SubCreateList------------------------------------

// AddPostToCommentInListByCreate
func AddPostToCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("AddPostToCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.ToCommentId <= 0 {
		utils.LogWarn("AddPostToCommentInListByCreate", "toComment == nil")
		return errors.New("nil_comment")
	}
	// 获取旧数据
	commentList, _ := GetPostToCommentListByCreate(pc.ToCommentId, -1, -1)
	if commentList == nil {
		commentList = make([]*entity.PostComment, 0)
	}
	// 添加指定元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostToCommentListByCreate(pc.ToCommentId, commentList, false)
	return err
}

// DelPostToCommentInListByCreate
func DelPostToCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("DelPostToCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.ToCommentId <= 0 {
		utils.LogWarn("DelPostToCommentInListByCreate", "toComment == nil")
		return errors.New("nil_comment")
	}
	// 获取旧数据
	commentList, _ := GetPostToCommentListByCreate(pc.ToCommentId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.PostComment, 0)
	for _, v := range commentList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != pc.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostToCommentListByCreate(pc.ToCommentId, saveList, false)
	return err
}

// UpdatePostToCommentInListByCreate
func UpdatePostToCommentInListByCreate(pc *entity.PostComment) error {
	if pc == nil || pc.Id <= 0 {
		utils.LogWarn("UpdatePostToCommentInListByCreate", "comment == nil")
		return errors.New("nil_comment")
	} else if pc.ToCommentId <= 0 {
		utils.LogWarn("UpdatePostToCommentInListByCreate", "toComment == nil")
		return errors.New("nil_comment")
	}
	// 获取旧列表
	commentList, _ := GetPostToCommentListByCreate(pc.ToCommentId, -1, -1)
	if commentList == nil || len(commentList) <= 0 {
		return nil
	}
	// 添加新元素
	commentList = append(commentList, pc)
	// 重新保存
	err := SetPostToCommentListByCreate(pc.ToCommentId, commentList, false)
	return err
}

// SetPostToCommentListByCreate
// merge 是否合并已有的数据
func SetPostToCommentListByCreate(cid int64, list []*entity.PostComment, merge bool) error {
	if cid <= 0 {
		utils.LogWarn("SetPostToCommentListByCreate", "cid <= 0")
		return errors.New("nil_couple")
	} else if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.PostComment, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisCommentList, _ := GetPostToCommentListByCreate(cid, -1, -1)
		if redisCommentList != nil && len(redisCommentList) > 0 {
			for _, v := range redisCommentList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedPostComment(returnList)
	// 开始降序排序
	sort.Sort(postCommentCreateSlice(returnList))
	// 构造数据
	saveList := &postCommentList{
		PostCommentList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostToCommentListByCreate", err)
		// list添加失败不删除
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储comment
	_, err = conn.Do("SET", KEY_CID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(cid, 10), listBody)
	if err != nil {
		utils.LogErr("SetPostToCommentListByCreate", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_CID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(cid, 10), getRedisPostCommentExpire())
	utils.LogErr("SetPostToCommentListByCreate", err)
	return nil
}

// GetPostToCommentListByCreate
func GetPostToCommentListByCreate(cid int64, offset, limit int) ([]*entity.PostComment, error) {
	if cid <= 0 {
		utils.LogWarn("GetPostToCommentListByCreate", "cid <= 0")
		return nil, errors.New("nil_couple")
	} else if limit == 0 {
		utils.LogWarn("GetPostToCommentListByCreate", "limit == 0")
		return nil, nil
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_CID_CREATE_POST_COMMENT_LIST+strconv.FormatInt(cid, 10))
	if err != nil {
		utils.LogErr("GetPostToCommentListByCreate", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析comment
	list := &postCommentList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostToCommentListByCreate", err)
		return nil, err
	}
	commentList := list.PostCommentList
	// 数量不够
	if commentList == nil || len(commentList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return commentList, nil
	}
	// 正常获取，数量够
	returnList := commentList[offset : offset+limit]
	return returnList, nil
}
