package common

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ConnectMongoTable(databaseName, tableName string) (*mongo.Collection, *mongo.Client, error) {
	// 设置数据库连接选项
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// 检查连接是否成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	return client.Database(databaseName).Collection(tableName), client, nil
}

func CheckDataExists(collection *mongo.Collection, filter bson.M) (bool, error) {
	// 指定查询选项，这里我们只需要判断数据是否存在，不需要返回具体数据，所以 limit 设为 1 即可
	options := options.Count().SetLimit(1)

	// 使用传入的 filter 参数进行数据查询
	count, err := collection.CountDocuments(context.Background(), filter, options)
	if err != nil {
		return false, err
	}

	// 如果 count > 0，则表示数据存在
	return count > 0, nil
}

func IsExistsUserName(collection *mongo.Collection, userName string) bool {
	filter := bson.M{"userName": userName}
	exists, err := CheckDataExists(collection, filter)

	// 查询失败都算存在，覆盖逻辑
	if err != nil {
		return true
	}

	return exists
}

func IsExistsEamil(collection *mongo.Collection, userEmail string) bool {
	filter := bson.M{"userEmail": userEmail}
	exists, err := CheckDataExists(collection, filter)

	// 查询失败都算存在，覆盖逻辑
	if err != nil {
		return true
	}

	return exists
}

func IsNowBeforeTargetTime(targetTime int64) bool {
	targetTimestamp := int64(targetTime)
	target := time.Unix(targetTimestamp, 0)

	currentTime := time.Now()
	return currentTime.Before(target)
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// return privateKey, publicKey
func GenerateRSAKeyPair(bitSize int) (string, string, error) {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return "", "", fmt.Errorf("无法生成RSA密钥对：%v", err)
	}

	// 将私钥编码为PEM格式
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	// 将公钥编码为PEM格式
	publicKey := &privateKey.PublicKey
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("无法编码公钥为PEM格式：%v", err)
	}
	publicKeyStr := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))

	return privateKeyStr, publicKeyStr, nil
}

func ProtoToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true, // 包含未赋值的字段
		Indent:          "  ", // 缩进两个空格
	}

	jsonBytes, err := marshaler.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("无法将protobuf转换为JSON：%v", err)
	}

	return string(jsonBytes), nil
}

func RSAEncrypt(publicKeyPEM, plaintext string) (string, error) {
	// 解码公钥 PEM 格式字符串
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode public key")
	}

	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 将公钥转换为 *rsa.PublicKey 类型
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("Invalid public key type")
	}

	// 加密明文
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(plaintext))
	if err != nil {
		return "", err
	}

	// 返回加密后的密文字符串
	return string(ciphertext), nil
}

func RSADecrypt(privateKeyPEM, ciphertext string) (string, error) {
	// 解码私钥 PEM 格式字符串
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 解密密文
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(ciphertext))
	if err != nil {
		return "", err
	}

	// 返回解密后的明文字符串
	return string(plaintext), nil
}
