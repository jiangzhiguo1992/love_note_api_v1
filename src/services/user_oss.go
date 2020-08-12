package services

import (
	"strconv"
	"strings"
	"time"

	"libs/aliyun"
	"libs/utils"
)

// 策略中可读取的路径
func getOssReadPath(cid int64) []string {
	bucket := utils.GetConfigStr("conf", "app.conf", "oss", "bucket")
	paths := make([]string, 0)
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_VERSION+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_NOTICE+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_BROADCAST+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_SUGGEST+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_COUPLE+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_TOPIC+"*")
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_MORE+"*")
	if cid > 0 { // 没有配对不给小本本读取权限
		cidStr := strconv.FormatInt(cid, 10)
		pathNoteCid := strings.Replace(OSS_PATH_NOTE_CID, "?", cidStr, -1)
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathNoteCid+"*")
	}
	return paths
}

// 策略中可供查询的路径
func getOssInfoPath(cid int64) []string {
	bucket := utils.GetConfigStr("conf", "app.conf", "oss", "bucket")
	paths := make([]string, 0)
	if cid > 0 { // 没有配对不给小本本查询权限
		cidStr := strconv.FormatInt(cid, 10)
		pathNoteCid := strings.Replace(OSS_PATH_NOTE_CID, "?", cidStr, -1)
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathNoteCid+"*")
	}
	return paths
}

// 策略中可供上传的路径
func getOssWritePath(uid, cid int64) []string {
	bucket := utils.GetConfigStr("conf", "app.conf", "oss", "bucket")
	uidStr := strconv.FormatInt(uid, 10)
	cidStr := strconv.FormatInt(cid, 10)

	paths := make([]string, 0)
	paths = append(paths, OSS_PREDIX+bucket+"/"+OSS_PATH_LOG+"*")
	if uid > 0 { // 需要登录才行的
		pathSuggestUid := strings.Replace(OSS_PATH_SUGGEST_UID, "?", uidStr, -1)
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathSuggestUid+"*")
	}
	if cid > 0 { // 需要配对才行的
		pathCoupleCid := strings.Replace(OSS_PATH_COUPLE_CID, "?", cidStr, -1)
		pathNoteCid := strings.Replace(OSS_PATH_NOTE_CID, "?", cidStr, -1)
		pathTopicCid := strings.Replace(OSS_PATH_TOPIC_CID, "?", cidStr, -1)
		pathMoreCid := strings.Replace(OSS_PATH_MORE_CID, "?", cidStr, -1)
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathCoupleCid+"*")
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathNoteCid+"*")
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathTopicCid+"*")
		paths = append(paths, OSS_PREDIX+bucket+"/"+pathMoreCid+"*")
	}
	return paths
}

// GetOssInfoByUserCouple
func GetOssInfoByUserCouple(uid, cid int64, userToken string) (*OssInfo, error) {
	// 获取policy策略信息
	readPaths := getOssReadPath(cid)
	infoPaths := getOssInfoPath(cid)
	writePaths := getOssWritePath(uid, cid)
	policy := aliyun.GetPolicy(readPaths, infoPaths, writePaths)
	// 获取sts凭证
	stsExpireSec := utils.GetConfigInt64("conf", "app.conf", "sts", "expire_min") * 60
	sts, err := aliyun.GetSts(policy, userToken, stsExpireSec)
	if err != nil || sts == nil {
		return nil, err
	}
	// 解析sts凭证
	credentials := sts.Credentials
	ossInfo := &OssInfo{}
	ossInfo.AccessKeyId = credentials.AccessKeyId
	ossInfo.AccessKeySecret = credentials.AccessKeySecret
	ossInfo.SecurityToken = credentials.SecurityToken
	// 封装oss信息
	ossInfo.Region = utils.GetConfigStr("conf", "app.conf", "oss", "region")
	ossInfo.Domain = utils.GetConfigStr("conf", "app.conf", "oss", "domain")
	ossInfo.Bucket = utils.GetConfigStr("conf", "app.conf", "oss", "bucket")
	// expire过期时间
	local := time.Now().Local().Unix()
	utc := time.Now().UTC().Unix()
	between := local - utc
	expireUTC := credentials.Expiration
	timeUTC, _ := time.Parse(utils.TIME_UTC_FORMAT, expireUTC)
	expireLocal := timeUTC.Unix() + between
	ossInfo.StsExpireTime = expireLocal
	ossInfo.OssRefreshSec = utils.GetConfigInt64("conf", "app.conf", "oss", "refresh_min") * 60
	ossInfo.UrlExpireSec = utils.GetConfigInt64("conf", "app.conf", "oss", "url_expire_min") * 60
	// 解析用户的path
	uidStr := strconv.FormatInt(uid, 10)
	cidStr := strconv.FormatInt(cid, 10)
	ossInfo.PathLog = OSS_PATH_LOG
	ossInfo.PathSuggest = strings.Replace(OSS_PATH_SUGGEST_UID, "?", uidStr, -1)
	ossInfo.PathCoupleAvatar = strings.Replace(OSS_PATH_COUPLE_AVATAR, "?", cidStr, -1)
	ossInfo.PathCoupleWall = strings.Replace(OSS_PATH_COUPLE_WALL, "?", cidStr, -1)
	ossInfo.PathNoteVideo = strings.Replace(OSS_PATH_NOTE_VIDEO, "?", cidStr, -1)
	ossInfo.PathNoteVideoThumb = strings.Replace(OSS_PATH_NOTE_VIDEO_THUMB, "?", cidStr, -1)
	ossInfo.PathNoteAudio = strings.Replace(OSS_PATH_NOTE_AUDIO, "?", cidStr, -1)
	ossInfo.PathNoteAlbum = strings.Replace(OSS_PATH_NOTE_ALBUM, "?", cidStr, -1)
	ossInfo.PathNotePicture = strings.Replace(OSS_PATH_NOTE_PICTURE, "?", cidStr, -1)
	ossInfo.PathNoteWhisper = strings.Replace(OSS_PATH_NOTE_WHISPER, "?", cidStr, -1)
	ossInfo.PathNoteDiary = strings.Replace(OSS_PATH_NOTE_DIARY, "?", cidStr, -1)
	ossInfo.PathNoteGift = strings.Replace(OSS_PATH_NOTE_GIFT, "?", cidStr, -1)
	ossInfo.PathNoteFood = strings.Replace(OSS_PATH_NOTE_FOOD, "?", cidStr, -1)
	ossInfo.PathNoteMovie = strings.Replace(OSS_PATH_NOTE_MOVIE, "?", cidStr, -1)
	ossInfo.PathTopicPost = strings.Replace(OSS_PATH_TOPIC_POST, "?", cidStr, -1)
	ossInfo.PathMoreMatch = strings.Replace(OSS_PATH_MORE_MATCH, "?", cidStr, -1)
	return ossInfo, nil
}

// GetOssInfoByAdmin
func GetOssInfoByAdmin(uid, cid int64, userToken string) (*OssInfo, error) {
	// 获取policy策略信息，不要上传权限
	bucket := utils.GetConfigStr("conf", "app.conf", "oss", "bucket")
	paths := make([]string, 0)
	paths = append(paths, OSS_PREDIX+bucket+"/*")
	policy := aliyun.GetPolicy(paths, paths, nil)
	// 获取sts凭证
	stsExpireSec := 60 * 60
	sts, err := aliyun.GetSts(policy, userToken, int64(stsExpireSec))
	if err != nil || sts == nil {
		return nil, err
	}
	// 解析sts凭证
	credentials := sts.Credentials
	ossInfo := &OssInfo{}
	ossInfo.AccessKeyId = credentials.AccessKeyId
	ossInfo.AccessKeySecret = credentials.AccessKeySecret
	ossInfo.SecurityToken = credentials.SecurityToken
	// 封装oss信息
	ossInfo.Region = utils.GetConfigStr("conf", "app.conf", "oss", "region")
	ossInfo.Domain = utils.GetConfigStr("conf", "app.conf", "oss", "domain")
	ossInfo.Bucket = bucket
	// expire过期时间
	local := time.Now().Local().Unix()
	utc := time.Now().UTC().Unix()
	between := local - utc
	expireUTC := credentials.Expiration
	timeUTC, _ := time.Parse(utils.TIME_UTC_FORMAT, expireUTC)
	expireLocal := timeUTC.Unix() + between
	ossInfo.StsExpireTime = expireLocal
	ossInfo.OssRefreshSec = utils.GetConfigInt64("conf", "app.conf", "oss", "refresh_min") * 60
	ossInfo.UrlExpireSec = utils.GetConfigInt64("conf", "app.conf", "oss", "url_expire_min") * 60
	// 返回path
	ossInfo.PathVersion = OSS_PATH_VERSION
	ossInfo.PathNotice = OSS_PATH_NOTICE
	ossInfo.PathBroadcast = OSS_PATH_BROADCAST
	return ossInfo, nil
}
