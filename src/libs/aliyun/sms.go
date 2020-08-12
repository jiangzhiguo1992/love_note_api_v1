package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"libs/utils"
	"math"
	"strconv"
)

// SmsReply 发送短信返回
type SmsReply struct {
	RequestId string `json:"RequestId"` // 8906582E-6722
	Code      string `json:"Code"`      // OK
	Message   string `json:"Message"`   // 请求成功
	BizId     string `json:"BizId"`     // 134523^4351232
}

const (
	SMS_URL = "http://dysmsapi.aliyuncs.com/"
)

func replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return rep.Replace(url.QueryEscape(in))
}

// SendSms 发送短信
func SendSms(phoneNumbers, templateCode, templateParam string) error {
	accessKeyID := utils.GetConfigStr("conf", "third.conf", "sms", "user_key_id")
	accessKeySecret := utils.GetConfigStr("conf", "third.conf", "sms", "user_key_secret")
	signName := utils.GetConfigStr("conf", "third.conf", "sms", "sign_name")
	// 参数封装
	params := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   createSmsSignatureNonce(),
		"AccessKeyId":      accessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format(utils.TIME_UTC_FORMAT),
		"Format":           "JSON",
		"Version":          "2017-05-25",
		"RegionId":         "cn-hangzhou",
		// 私有参数
		"Action":        "SendSms",
		"PhoneNumbers":  phoneNumbers,
		"SignName":      signName,
		"TemplateCode":  templateCode,
		"TemplateParam": templateParam,
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
	str := fmt.Sprintf(SMS_URL+"?Signature=%s%s", sign, sortQueryString)
	// 开始请求
	resp, err := http.Get(str)
	if err != nil {
		utils.LogErr("sms", err)
		return errors.New("sms_request_send_fail")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogErr("sms", err)
		return errors.New("sms_response_get_fail")
	}
	// 格式化返回信息
	result := &SmsReply{}
	if err := json.Unmarshal(body, result); err != nil || result == nil {
		utils.LogErr("sms", err)
		return errors.New("sms_response_decode_fail")
	}
	// 检查处理
	if strings.Contains(result.Code, "SignatureNonceUsed") {
		utils.LogDebug("sms", "SignatureNonceUsed")
		return SendSms(phoneNumbers, templateCode, templateParam)
	} else if strings.Contains(result.Code, "BUSINESS_LIMIT_CONTROL") {
		utils.LogDebug("sms", fmt.Sprintf("%+v", result))
		return errors.New("sms_send_frequent")
	} else if result.Code != "OK" {
		utils.LogDebug("sms", fmt.Sprintf("%+v", result))
		return errors.New(result.Code)
	}
	return nil
}

// createSmsSignatureNonce
func createSmsSignatureNonce() string {
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
