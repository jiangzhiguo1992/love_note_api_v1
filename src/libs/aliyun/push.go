package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"libs/utils"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// PushReply 推送返回
type PushReply struct {
	RequestId string `json:"RequestId"` // 8906582E-6722
	MessageId string `json:"MessageId"` // 响应参数
	Code      string `json:"Code"`      // OK
	Message   string `json:"Message"`   // 请求成功
}

const (
	PUSH_URL = "https://cloudpush.aliyuncs.com/"
)

func SendPushNotice2User(uid int64, platform, title, content, extraJson string) error {
	if uid <= 0 {
		return errors.New("nil_user")
	}
	var aliAppKey string
	p := strings.ToLower(strings.TrimSpace(platform))
	platformAndroid := utils.GetConfigStr("conf", "app.conf", "phone", "platform_android")
	platformIos := utils.GetConfigStr("conf", "app.conf", "phone", "platform_ios")
	if p == platformAndroid {
		aliAppKey = utils.GetConfigStr("conf", "third.conf", "push", "ali_android_app_key")
	} else if p == platformIos {
		aliAppKey = utils.GetConfigStr("conf", "third.conf", "push", "ali_ios_app_key")
	}
	return sendPush(aliAppKey, "NOTICE", "ACCOUNT", strings.TrimSpace(strconv.FormatInt(uid, 10)), title, content, extraJson)
}

// sendPush
func sendPush(aliAppKey, pushType, target, targetValue, title, content, extraJson string) error {
	if len(aliAppKey) <= 0 {
		return nil
	} else if len(title) <= 0 || len(content) <= 0 {
		return errors.New("limit_content_text_nil")
	}
	accessKeyID := utils.GetConfigStr("conf", "third.conf", "push", "user_key_id")
	accessKeySecret := utils.GetConfigStr("conf", "third.conf", "push", "user_key_secret")
	channelId := utils.GetConfigStr("conf", "third.conf", "push", "channel_id")
	var notifyType string
	sound := utils.GetConfigBool("conf", "third.conf", "push", "notice_sound")
	vibrate := utils.GetConfigBool("conf", "third.conf", "push", "notice_vibrate")
	if sound && vibrate {
		notifyType = "BOTH"
	} else if sound {
		notifyType = "SOUND"
	} else if vibrate {
		notifyType = "VIBRATE"
	} else {
		notifyType = "NONE"
	}
	// 参数封装
	params := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   createPushSignatureNonce(),
		"AccessKeyId":      accessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format(utils.TIME_UTC_FORMAT),
		"Format":           "JSON",
		"Version":          "2016-08-01", // openAPI 2.0概览 中查看
		"RegionId":         "cn-hangzhou",
		// 私有参数
		"AppKey":        aliAppKey,
		"Action":        "Push",
		"DeviceType":    "ALL",
		"PushType":      pushType,
		"Target":        target,
		"TargetValue":   targetValue,
		"Title":         title,
		"Body":          content,
		"ExtParameters": extraJson,
		// android
		"AndroidNotificationChannel":     channelId,
		"AndroidOpenType":                "APPLICATION",
		"AndroidNotifyType":              notifyType,
		"AndroidNotificationBarPriority": "2",
		"AndroidExtParameters":           extraJson,
		"StoreOffline":                   "true", // 默认不离线
		"AndroidPopupTitle":              title,
		"AndroidPopupBody":               content,
		//"ExpireTime": "", // 默认72小时过期
		//"PushTime":"", // 默认立即发送
		// ios
		"ApnsEnv":    "PRODUCT", // DEV/PRODUCT
		"iOSApnsEnv": "PRODUCT", // DEV/PRODUCT
		"iOSMusic":   "default",
		//"iOSBadge":              "1",
		"iOSBadgeAutoIncrement": "true", // 角标自增
		"iOSExtParameters":      extraJson,
		"iOSRemind":             "true",  // 离线消息转通知仅适用于生产环境
		"iOSRemindBody":         content, // 仅当iOSApnsEnv=PRODUCT && iOSRemind为true时有效
		//"iOSNotificationCategory":, // 指定iOS通知Category（iOS 10+）
		//"iOSMutableContent":, // 是否使能iOS通知扩展处理（iOS 10+）
	}
	// 先排序params
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 再连接params(包括连接key-value)，获得url除sign之外的params
	var sortQueryString string
	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, replace(v), replace(params[v]))
	}
	// 再替换/
	stringToSign := fmt.Sprintf("GET&%s&%s", replace("/"), replace(sortQueryString[1:]))
	// 再加密
	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessKeySecret)))
	mac.Write([]byte(stringToSign))
	sign := replace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	// 记得最后把sign加入params的后面
	str := fmt.Sprintf(PUSH_URL+"?Signature=%s%s", sign, sortQueryString)
	// 开始请求
	resp, err := http.Get(str)
	if err != nil {
		utils.LogErr("push-1", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogErr("push-2", err)
		return err
	}
	if body == nil || len(body) <= 0 {
		return nil
	}
	// 格式化返回信息
	result := &PushReply{}
	if err := json.Unmarshal(body, result); err != nil || result == nil {
		utils.LogErr("push-3", err)
		return err
	}
	// 检查处理
	if strings.Contains(result.Code, "SignatureNonceUsed") {
		utils.LogDebug("push-4", "SignatureNonceUsed")
		return sendPush(aliAppKey, pushType, target, targetValue, title, content, extraJson)
	} else if len(result.Code) > 0 {
		utils.LogWarn("push-5", fmt.Sprintf("%+v", result))
		return errors.New(result.Code)
	}
	return nil
}

// createPushSignatureNonce
func createPushSignatureNonce() string {
	// 前半截时间戳
	unixNa := time.Now().UnixNano()
	unix := strconv.FormatInt(unixNa, 16)
	// 后半截随机数
	length := 16
	max := math.Pow10(length) - 1
	min := math.Pow10(length - 1)
	rand16 := utils.GetRandRange(int(max), int(min))
	return unix + "_" + strconv.Itoa(rand16)
}
