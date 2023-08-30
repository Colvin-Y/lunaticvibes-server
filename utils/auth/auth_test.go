package auth_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	constant "github.com/Colvin-Y/lunaticvibes-server/const"
	"github.com/Colvin-Y/lunaticvibes-server/utils/auth"
	"github.com/Colvin-Y/lunaticvibes-server/utils/common"
)

/*
用例逻辑：
当前时间   用户操作时间			 			操作				过期时间										期待表现
T		   	T				  			GenRSAKeyByUser	  T+constant.Timeout							代表登录成功，生成对应的 key
T+3        	T+3				  			IsUserLegal       T+constant.Timeout => T+constant.Timeout+3	未过期，处于登录态且刷新过期时间
T+3			T+3+constant.Timeout-1		IsUserLegal		  T+3+constant.Timeout							未过期，处于登录态且刷新过期时间
T+3			T+3+constant.Timeout		IsUserLegal		  T+3+constant.Timeout							过期，删除 key
*/
func TestFlow(t *testing.T) {
	user := "user1" + strconv.Itoa(int(common.GetCurrentTimestamp()))
	dataStr := strconv.Itoa(int(common.GetCurrentTimestamp()))
	auth.GenRSAKeyByUser(user)
	privKey, pubKey := auth.GetRSAKeyByUser(user)
	encryptStr, err := common.RSAEncrypt(pubKey, dataStr)
	fmt.Printf("err: %v\n", err)
	decryptStr, err := common.RSADecrypt(privKey, encryptStr)
	fmt.Printf("err: %v\n", err)

	//fmt.Printf("pubKey: %v\n", pubKey)
	//fmt.Printf("privKey: %v\n", privKey)
	fmt.Printf("encryptStr: %v\n", encryptStr)
	fmt.Printf("decryptStr: %v\n", decryptStr)
	fmt.Printf("dataStr: %v\n", dataStr)

	if decryptStr != dataStr {
		t.Errorf("解析出错！")
	}

	time.Sleep(3 * time.Second)
	isLegal, reason := auth.IsUserLegal(user, encryptStr)
	if !isLegal {
		t.Errorf("解析出错！ reason:%v", reason)
	}

	timeinDataStr := strconv.Itoa(int(common.GetCurrentTimestamp() + constant.RSAKeyTimeout - 1))
	timeinStr, _ := common.RSAEncrypt(pubKey, timeinDataStr)
	isLegal, reason = auth.IsUserLegal(user, timeinStr)
	if !isLegal {
		t.Errorf("解析出错！ reason:%v", reason)
	}

	timeoutDataStr := strconv.Itoa(int(common.GetCurrentTimestamp() + constant.RSAKeyTimeout))
	timeinoutStr, _ := common.RSAEncrypt(pubKey, timeoutDataStr)
	isLegal, reason = auth.IsUserLegal(user, timeinoutStr)
	if isLegal {
		t.Errorf("解析出错！ reason:%v", reason)
	}
	if reason != constant.AUTH_KEY_TIMEOUT {
		t.Errorf("应当是超时导致的错误! reason:%v", reason)
	}

	isLegal, reason = auth.IsUserLegal(user, timeinoutStr)
	if isLegal || reason != constant.AUTH_RSAKEY_NOT_FOUND {
		t.Errorf("应当不存在该 key, reason:%v", reason)
	}
}
