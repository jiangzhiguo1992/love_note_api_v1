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

// removeRepeatedPost
func removeRepeatedPost(list []*entity.Post) ([]*entity.Post) {
	returnList := make([]*entity.Post, 0)
	if list == nil || len(list) <= 0 {
		return returnList
	}
	// 遍历所有元素
	for i := 0; i < len(list); i++ {
		iPost := list[i]
		if iPost == nil || iPost.Id <= 0 {
			continue
		}
		// 检查重复元素 0=10 1=11 2=12 3=11 4=12 5=10 6=12
		repeat := false
		for j := i + 1; j < len(list); j++ {
			jPost := list[j]
			if jPost == nil || jPost.Id <= 0 {
				continue
			}
			if iPost.Id == jPost.Id {
				// 更新较晚的加上，相同的话 加后面的
				if iPost.UpdateAt <= jPost.UpdateAt {
					repeat = true
				} else if iPost.UpdateAt > jPost.UpdateAt {
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
			iPost.Mine = false
			iPost.Our = false
			iPost.Read = false
			iPost.Report = false
			iPost.Point = false
			iPost.Collect = false
			iPost.Comment = false
			returnList = append(returnList, iPost)
		}
	}
	return returnList
}

// SetPost
func SetPost(post *entity.Post) error {
	if post == nil || post.Id <= 0 {
		utils.LogWarn("SetPost", "无效的帖子: "+fmt.Sprintf("%+v", post))
		return errors.New("nil_post")
	}
	post.Mine = false
	post.Our = false
	post.Read = false
	post.Report = false
	post.Point = false
	post.Collect = false
	post.Comment = false
	bytes, err := json.Marshal(post)
	if err != nil {
		utils.LogErr("SetPost", err)
		// 添加失败就删除
		DelPost(post)
		return errors.New("data_decode_err")
	}
	postBody := string(bytes)
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
	// 存储post
	if post.Id > 0 {
		_, err = conn.Do("SET", KEY_ID_POST+strconv.FormatInt(post.Id, 10), postBody)
		if err != nil {
			utils.LogErr("SetPost", err)
			return err
		}
		_, err = conn.Do("EXPIRE", KEY_ID_POST+strconv.FormatInt(post.Id, 10), getRedisPostExpire())
		utils.LogErr("SetPost", err)
	} else {
		utils.LogWarn("SetPost", "pid <=0")
	}
	return err
}

// DelPost
func DelPost(post *entity.Post) error {
	if post == nil || post.Id <= 0 {
		utils.LogWarn("DelPost", "无效的帖子: "+fmt.Sprintf("%+v", post))
		return errors.New("nil_post")
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
	var err error
	if post.Id > 0 {
		_, err = conn.Do("DEL", KEY_ID_POST+strconv.FormatInt(post.Id, 10))
		utils.LogErr("DelPost", err)
	} else {
		utils.LogWarn("DelPost", "pid <=0")
	}
	return err
}

// GetPostById
func GetPostById(pid int64) (*entity.Post, error) {
	if pid <= 0 {
		utils.LogWarn("GetPostById", "pid <= 0")
		return nil, errors.New("nil_post")
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
	reply, err := conn.Do("GET", KEY_ID_POST+strconv.FormatInt(pid, 10))
	if err != nil {
		utils.LogErr("GetPostById", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析post
	post := &entity.Post{}
	err = json.Unmarshal(bytes, post)
	if err != nil {
		utils.LogErr("GetPostById", err)
		return nil, err
	}
	return post, nil
}

//------------------------------------UpdateList------------------------------------

// AddPostInListByAll
func AddPostInListByAll(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("AddPostInListByAll", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	postList, _ := GetPostListByAll(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		postList = make([]*entity.Post, 0)
	}
	// 添加指定元素
	postList = append(postList, p)
	// 重新保存
	err := SetPostListByAll(p.Kind, p.SubKind, postList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByAll(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// DelPostInListByAll
func DelPostInListByAll(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("DelPostInListByAll", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	postList, _ := GetPostListByAll(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.Post, 0)
	for _, v := range postList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != p.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostListByAll(p.Kind, p.SubKind, saveList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByAll(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// UpdatePostInListByAll
func UpdatePostInListByAll(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("UpdatePostInListByAll", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧列表
	postList, _ := GetPostListByAll(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 添加新元素
	postList = append(postList, p)
	// 重新保存
	err := SetPostListByAll(p.Kind, p.SubKind, postList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByAll(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// SetPostListByAll
// merge 是否合并已有的数据
func SetPostListByAll(kind, subKind int, list []*entity.Post, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.Post, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisPostList, _ := GetPostListByAll(kind, subKind, -1, -1)
		if redisPostList != nil && len(redisPostList) > 0 {
			for _, v := range redisPostList {
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
	returnList = removeRepeatedPost(returnList)
	// 开始降序排序
	sort.Sort(postUpdateSlice(returnList))
	// 构造数据
	saveList := &postList{
		PostList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostListByAll", err)
		// list添加失败不删除
		if merge {

		}
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
	// 存储post
	_, err = conn.Do("SET", KEY_KIND_SUB_KIND_ALL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), listBody)
	if err != nil {
		utils.LogErr("SetPostListByAll", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_KIND_SUB_KIND_ALL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), getRedisPostExpire())
	utils.LogErr("SetPostListByAll", err)
	return nil
}

// GetPostListByAll
func GetPostListByAll(kind, subKind, offset, limit int) ([]*entity.Post, error) {
	if limit == 0 {
		utils.LogWarn("GetPostListByAll", "limit == 0")
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
	reply, err := conn.Do("GET", KEY_KIND_SUB_KIND_ALL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind))
	if err != nil {
		utils.LogErr("GetPostListByAll", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析post
	list := &postList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostListByAll", err)
		return nil, err
	}
	postList := list.PostList
	// 数量不够
	if postList == nil || len(postList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return postList, nil
	}
	// 正常获取，数量够
	returnList := postList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------WellList------------------------------------

// DelPostInListByWell
func DelPostInListByWell(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("DelPostInListByWell", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	postList, _ := GetPostListByWell(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.Post, 0)
	for _, v := range postList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != p.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostListByWell(p.Kind, p.SubKind, saveList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByWell(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// UpdatePostInListByWell
func UpdatePostInListByWell(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("UpdatePostInListByWell", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧列表
	postList, _ := GetPostListByWell(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 添加新元素
	postList = append(postList, p)
	// 重新保存
	err := SetPostListByWell(p.Kind, p.SubKind, postList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByWell(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// SetPostListByWell
// merge 是否合并已有的数据
func SetPostListByWell(kind, subKind int, list []*entity.Post, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.Post, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisPostList, _ := GetPostListByWell(kind, subKind, -1, -1)
		if redisPostList != nil && len(redisPostList) > 0 {
			for _, v := range redisPostList {
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
	returnList = removeRepeatedPost(returnList)
	// 开始降序排序
	sort.Sort(postUpdateSlice(returnList))
	// 构造数据
	saveList := &postList{
		PostList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostListByWell", err)
		// list添加失败不删除
		if merge {

		}
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
	// 存储post
	_, err = conn.Do("SET", KEY_KIND_SUB_KIND_WELL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), listBody)
	if err != nil {
		utils.LogErr("SetPostListByWell", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_KIND_SUB_KIND_WELL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), getRedisPostExpire())
	utils.LogErr("SetPostListByWell", err)
	return nil
}

// GetPostListByWell
func GetPostListByWell(kind, subKind, offset, limit int) ([]*entity.Post, error) {
	if limit == 0 {
		utils.LogWarn("GetPostListByWell", "limit == 0")
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
	reply, err := conn.Do("GET", KEY_KIND_SUB_KIND_WELL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind))
	if err != nil {
		utils.LogErr("GetPostListByWell", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析post
	list := &postList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostListByWell", err)
		return nil, err
	}
	postList := list.PostList
	// 数量不够
	if postList == nil || len(postList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return postList, nil
	}
	// 正常获取，数量够
	returnList := postList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------OfficialList------------------------------------

// DelPostInListByOfficial
func DelPostInListByOfficial(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("DelPostInListByOfficial", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧数据
	postList, _ := GetPostListByOfficial(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.Post, 0)
	for _, v := range postList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != p.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetPostListByOfficial(p.Kind, p.SubKind, saveList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByOfficial(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// UpdatePostInListByOfficial
func UpdatePostInListByOfficial(p *entity.Post) error {
	if p == nil || p.Id <= 0 {
		utils.LogWarn("UpdatePostInListByOfficial", "post == nil")
		return errors.New("nil_post")
	}
	// 获取旧列表
	postList, _ := GetPostListByOfficial(p.Kind, p.SubKind, -1, -1)
	if postList == nil || len(postList) <= 0 {
		return nil
	}
	// 添加新元素
	postList = append(postList, p)
	// 重新保存
	err := SetPostListByOfficial(p.Kind, p.SubKind, postList, false)
	if err == nil {
		// 保存全部的分类
		if p.SubKind != entity.POST_SUB_KIND_ALL {
			err = SetPostListByOfficial(p.Kind, entity.POST_SUB_KIND_ALL, postList, false)
		}
	}
	return err
}

// SetPostListByOfficial
// merge 是否合并已有的数据
func SetPostListByOfficial(kind, subKind int, list []*entity.Post, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.Post, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisPostList, _ := GetPostListByOfficial(kind, subKind, -1, -1)
		if redisPostList != nil && len(redisPostList) > 0 {
			for _, v := range redisPostList {
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
	returnList = removeRepeatedPost(returnList)
	// 开始降序排序
	sort.Sort(postUpdateSlice(returnList))
	// 构造数据
	saveList := &postList{
		PostList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetPostListByOfficial", err)
		// list添加失败不删除
		if merge {

		}
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
	// 存储post
	_, err = conn.Do("SET", KEY_KIND_SUB_KIND_OFFICIAL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), listBody)
	if err != nil {
		utils.LogErr("SetPostListByOfficial", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_KIND_SUB_KIND_OFFICIAL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind), getRedisPostExpire())
	utils.LogErr("SetPostListByOfficial", err)
	return nil
}

// GetPostListByOfficial
func GetPostListByOfficial(kind, subKind, offset, limit int) ([]*entity.Post, error) {
	if limit == 0 {
		utils.LogWarn("GetPostListByOfficial", "limit == 0")
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
	reply, err := conn.Do("GET", KEY_KIND_SUB_KIND_OFFICIAL_POST_LIST+strconv.Itoa(kind)+":"+strconv.Itoa(subKind))
	if err != nil {
		utils.LogErr("GetPostListByOfficial", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析post
	list := &postList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetPostListByOfficial", err)
		return nil, err
	}
	postList := list.PostList
	// 数量不够
	if postList == nil || len(postList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return postList, nil
	}
	// 正常获取，数量够
	returnList := postList[offset : offset+limit]
	return returnList, nil
}
