/*
权限校验相关模块，包括登录态校验以及密钥的发放及维护
使用方式:

	Client															Server
	login															GenRSAKeyByUser
	token=common.RSAEncrypt(pubKey, current_timestamp_str)
	do_action(data, userUUID, token)								IsUserLegal
	logout															DeleteRSAKeyByUser
*/
package auth

import (
	"strconv"

	constant "github.com/Colvin-Y/lunaticvibes-server/const"

	"github.com/Colvin-Y/lunaticvibes-server/utils/common"
)

type RSAKeyInfo struct {
	PrivateKey        string
	PublicKey         string
	LoginTimestamp    int64
	DeadlineTimestamp int64
}

var rsaKeyInfo map[string]RSAKeyInfo

func init() {
	rsaKeyInfo = make(map[string]RSAKeyInfo)
}

// return (privKey, pubKey)
func GetRSAKeyByUser(userUUID string) (priv string, pub string) {
	if value, ok := rsaKeyInfo[userUUID]; ok {
		return value.PrivateKey, value.PublicKey
	}

	return "", ""
}

// logout 或者判断用户登录态失效的时候使用
func DeleteRSAKeyByUser(userUUID string) {
	delete(rsaKeyInfo, userUUID)
}

// 续费租约, 须确保 userUUID 存在
func leaseRenewal(userUUID string) {
	if info, ok := rsaKeyInfo[userUUID]; ok {
		info.DeadlineTimestamp = common.GetCurrentTimestamp() + constant.RSAKeyTimeout
		rsaKeyInfo[userUUID] = info
	}
}

// 强制覆盖 RSAKey，这样对于同一用户如果在不同地点登录就可以实现抢占
func GenRSAKeyByUser(userUUID string) error {
	privateKey, pubKey, err := common.GenerateRSAKeyPair(constant.RSAKeyBitsize)
	if err != nil {
		return err
	}
	rsaKeyInfo[userUUID] = RSAKeyInfo{privateKey, pubKey, common.GetCurrentTimestamp(), common.GetCurrentTimestamp() + constant.RSAKeyTimeout}
	return nil
}

// 登录态校验 + 数据校验
// token 是用户发消息时候 Timestamp 的加密值
// 如果用户成功登录，则续费过期时间
func IsUserLegal(userUUID, token string) (isLegal bool, reason int) {
	privKey, _ := GetRSAKeyByUser(userUUID)
	if privKey == "" {
		return false, constant.AUTH_RSAKEY_NOT_FOUND
	}

	timestampStr, err := common.RSADecrypt(privKey, token)
	if err != nil {
		// 解析不了代表有人冒充 user 发消息，直接拒绝
		return false, constant.AUTH_PARSE_ERR
	}

	// 获取消息编码的时间戳
	msgTimestamp, err := strconv.Atoi(timestampStr)
	if err != nil {
		// 解析不了代表有人冒充 user 发消息，直接拒绝
		return false, constant.AUTH_PARSE_ERR
	}

	// 验证登录态是否过期
	if common.IsNowBeforeTargetTime(rsaKeyInfo[userUUID].DeadlineTimestamp) && int64(msgTimestamp) < rsaKeyInfo[userUUID].DeadlineTimestamp && int64(msgTimestamp) >= rsaKeyInfo[userUUID].LoginTimestamp {
		leaseRenewal(userUUID)
		return true, constant.SUCCESS
	} else {
		// 解析对了，说明用户侧有正确的 pubkey，但是过期了
		DeleteRSAKeyByUser(userUUID)
		return false, constant.AUTH_KEY_TIMEOUT
	}
}
