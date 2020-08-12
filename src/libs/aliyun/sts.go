package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"libs/utils"
	"math"
)

type (
	// url参数
	UrlParams struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	// 策略集合
	Policy struct {
		Version   string       `json:"Version"` // 1
		Statement []*Statement `json:"Statement"`
	}
	Statement struct {
		Effect   string   `json:"Effect"`   // Allow/Deny
		Action   []string `json:"Action"`   // "oss:ListObjects", "oss:GetObject"
		Resource []string `json:"Resource"` // "acs:oss:*:*:ram-test-dev", "acs:oss:*:*:ram-test-dev/*"
	}
	// STS
	AliSts struct {
		RequestId       string           `json:"requestId"` // "6894B13B-6D71-4EF5-88FA-F32781734A7F"
		HostId          string           `json:"hostId"`    // "sts.aliyuncs.com"
		Code            string           `json:"code"`      // "InvalidParameter"
		Message         string           `json:"message"`   // "The specified parameter \"Action or Version\" is not valid."
		Credentials     *Credentials     `json:"credentials"`
		AssumedRoleUser *AssumedRoleUser `json:"assumedRoleUser"`
	}
	Credentials struct {
		AccessKeyId     string `json:"accessKeyId"`     // "STS.L4aBSCSJVMuKg5U1vFDw"
		AccessKeySecret string `json:"accessKeySecret"` // "wyLTSmsyPGP1ohvvw8xYgB29dlGI8KMiH2pKCNZ9",
		Expiration      string `json:"expiration"`      // "2015-04-09T11:52:19Z",
		SecurityToken   string `json:"securityToken"`   // "CAESrAIIARKAAShQquMnLIlbvEcIxO6wCoqJufs8sWwieUxu45hS9AvKNEte8KRUWiJWJ6Y+YHAPgNwi7yfRecMFydL2uPOgBI7LDio0RkbYLmJfIxHM2nGBPdml7kYEOXmJp2aDhbvvwVYIyt/8iES/R6N208wQh0Pk2bu+/9dvalp6wOHF4gkFGhhTVFMuTDRhQlNDU0pWTXVLZzVVMXZGRHciBTQzMjc0KgVhbGljZTCpnJjwySk6BlJzYU1ENUJuCgExGmkKBUFsbG93Eh8KDEFjdGlvbkVxdWFscxIGQWN0aW9uGgcKBW9zczoqEj8KDlJlc291cmNlRXF1YWxzEghSZXNvdXJjZRojCiFhY3M6b3NzOio6NDMyNzQ6c2FtcGxlYm94L2FsaWNlLyo="
	}
	AssumedRoleUser struct {
		Arn               string `json:"arn"`               // "acs:sts::1234567890123456:assumed-role/AdminRole/alice",
		AssumedRoleUserId string `json:"assumedRoleUserId"` // "344584339364951186:alice"
	}
)

// 策略子集
func getStatementRead(paths []string) *Statement {
	statement := &Statement{
		Effect:   "Allow",
		Action:   []string{"oss:Get*"},
		Resource: paths,
	}
	return statement
}

// 策略子集
func getStatementInfo(paths []string) *Statement {
	statement := &Statement{
		Effect:   "Allow",
		Action:   []string{"oss:List*"},
		Resource: paths,
	}
	return statement
}

// 策略子集
func getStatementWrite(paths []string) *Statement {
	statement := &Statement{
		Effect:   "Allow",
		Action:   []string{"oss:Put*", "oss:Abort*"},
		Resource: paths,
	}
	return statement
}

// 获取所有策略(最终策略取和后台角色的交集)
func GetPolicy(readPaths, infoPaths, writePaths []string) *Policy {
	policy := &Policy{}
	policy.Version = "1"
	policy.Statement = make([]*Statement, 0)
	if readPaths != nil && len(readPaths) > 0 {
		read := getStatementRead(readPaths)
		policy.Statement = append(policy.Statement, read)
	}
	if infoPaths != nil && len(infoPaths) > 0 {
		info := getStatementInfo(infoPaths)
		policy.Statement = append(policy.Statement, info)
	}
	if writePaths != nil && len(writePaths) > 0 {
		write := getStatementWrite(writePaths)
		policy.Statement = append(policy.Statement, write)
	}
	return policy
}

// getSts
// 1.先利用keyId和keySecret，登录用户 (用户可以没有权限)
// 2.再利用roleName获取角色及其权限 (只有登录了用户，才能扮演角色)
// 3.再返回policy与角色交集的权限 (没有policy则返回角色的权限)
func GetSts(policy *Policy, userToken string, stsExpireSec int64) (*AliSts, error) {
	url := getStsUrl(policy, userToken, stsExpireSec)
	// request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.LogErr("sts", err)
		return nil, errors.New("res_access_perm_make_fail")
	}
	// http
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.LogErr("sts", err)
		return nil, errors.New("res_access_perm_request_fail")
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.LogErr("sts", err)
		return nil, errors.New("res_access_perm_read_fail")
	}
	result := &AliSts{}
	err = json.Unmarshal(bytes, result)
	if err != nil || result == nil {
		utils.LogErr("sts", err)
		return nil, errors.New("res_access_perm_decode_fail")
	}
	if strings.Contains(result.Code, "SignatureNonceUsed") {
		return GetSts(policy, userToken, stsExpireSec)
	} else if len(strings.TrimSpace(result.Message)) > 0 && (result.Credentials == nil || len(result.Credentials.AccessKeyId) <= 0) {
		err = errors.New(result.Message)
		utils.LogErr("sts", err)
		return nil, err
	}
	return result, nil
}

// 获取sts的请求路径(没有sdk，只能自己拼)
func getStsUrl(policy *Policy, session string, stsExpireSec int64) string {
	stsAliId := utils.GetConfigStr("conf", "app.conf", "sts", "ali_id")
	stsKeyId := utils.GetConfigStr("conf", "app.conf", "sts", "user_key_id")
	stsKeySecret := utils.GetConfigStr("conf", "app.conf", "sts", "user_key_secret")
	stsRoleName := utils.GetConfigStr("conf", "app.conf", "sts", "role_name")

	// session 2-32个字符
	if len(session) < 2 { // 随机数
		int63 := rand.Int63()
		session = strconv.FormatInt(int63, 10)
	} else if len(session) > 32 {
		rs := []rune(session)
		session = string(rs[0:31])
	}
	bytes, _ := json.Marshal(policy)
	policySrt := string(bytes)
	// url
	urlQuerys := make([]*UrlParams, 0)
	// url参数 首字母排序 utf8转asc码 ( /->47->%2F，&->38->%26，=->61->%3D，:->58->%3A，%->37->%25)
	urlQuerys = append(urlQuerys, &UrlParams{"AccessKeyId", stsKeyId})
	urlQuerys = append(urlQuerys, &UrlParams{"Action", "AssumeRole"})
	urlQuerys = append(urlQuerys, &UrlParams{"DurationSeconds", strconv.FormatInt(stsExpireSec, 10)})
	urlQuerys = append(urlQuerys, &UrlParams{"Format", "JSON"})
	urlQuerys = append(urlQuerys, &UrlParams{"Policy", policySrt})
	urlQuerys = append(urlQuerys, &UrlParams{"RoleArn", "acs:ram::" + stsAliId + ":role/" + stsRoleName})
	urlQuerys = append(urlQuerys, &UrlParams{"RoleSessionName", session})
	urlQuerys = append(urlQuerys, &UrlParams{"SignatureMethod", "HMAC-SHA1"})
	urlQuerys = append(urlQuerys, &UrlParams{"SignatureNonce", createStsSignatureNonce()})
	urlQuerys = append(urlQuerys, &UrlParams{"SignatureVersion", "1.0"})
	urlQuerys = append(urlQuerys, &UrlParams{"Timestamp", time.Now().UTC().Format(utils.TIME_UTC_FORMAT)})
	urlQuerys = append(urlQuerys, &UrlParams{"Version", "2015-04-01"})
	// 转ascii
	stringToSign := "GET&%2F&"
	for i := 0; i < len(urlQuerys); i++ {
		query := urlQuerys[i]
		// 第一次为url形式
		key := utf8ConvertASCII(query.Key)
		value := utf8ConvertASCII(query.Value)
		// 第二次为signature形式
		s := utf8ConvertASCII(key) + utf8ConvertASCII("=") + utf8ConvertASCII(value)
		if i < len(urlQuerys)-1 {
			s = s + utf8ConvertASCII("&")
		}
		stringToSign = stringToSign + s
	}
	utils.LogDebug("sts", "STS-ToSignature:"+stringToSign)
	// sign的获取: 先sha1的hmac，再base64
	hash := hmac.New(sha1.New, []byte(stsKeySecret+"&"))
	hash.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	urlQuerys = append(urlQuerys, &UrlParams{"Signature", signature})
	urlPath := "https://sts.aliyuncs.com/?"
	for _, v := range urlQuerys {
		key := utf8ConvertASCII(v.Key)
		value := utf8ConvertASCII(v.Value)
		urlPath = urlPath + key + "=" + value + "&"
	}
	strings.TrimSuffix(urlPath, "&")
	utils.LogDebug("sts", "STS-SignatureURL:"+urlPath)
	return urlPath
}

// 路径特殊字符串的转化
func utf8ConvertASCII(utf8 string) string {
	regex2AscII := "^[^0-9a-zA-Z-_.~]$"
	ascII := ""
	rs := []rune(utf8)
	for i := 0; i < len(rs); i++ {
		char := string(rs[i])
		if regexp.MustCompile(regex2AscII).MatchString(char) {
			//asc16 := int64(rs[i])
			// 十进制转16进制
			char = "%" + fmt.Sprintf("%02X", rs[i])
		}
		ascII = ascII + char
	}
	return ascII
}

// createStsSignatureNonce
func createStsSignatureNonce() string {
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
