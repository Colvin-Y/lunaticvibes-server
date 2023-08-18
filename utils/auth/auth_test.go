package auth_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Colvin-Y/lunaticvibes-server/utils/auth"
	"github.com/Colvin-Y/lunaticvibes-server/utils/common"
)

func TestFlow(t *testing.T) {
	user := "user1"
	dataStr := strconv.Itoa(int(common.GetCurrentTimestamp()))
	auth.GenRSAKeyByUser(user)
	privKey, pubKey := auth.GetRSAKeyByUser(user)
	encryptStr, err := common.RSAEncrypt(pubKey, dataStr)
	fmt.Printf("err: %v\n", err)
	decryptStr, err := common.RSADecrypt(privKey, encryptStr)
	fmt.Printf("err: %v\n", err)

	fmt.Printf("pubKey: %v\n", pubKey)
	fmt.Printf("privKey: %v\n", privKey)
	fmt.Printf("encryptStr: %v\n", encryptStr)
	fmt.Printf("decryptStr: %v\n", decryptStr)
	fmt.Printf("dataStr: %v\n", dataStr)

	if decryptStr != dataStr {
		t.Errorf("解析出错！")
	}

	isLegal := auth.IsUserLegal(user, encryptStr)
	if !isLegal {
		t.Errorf("解析出错！")
	}
}
