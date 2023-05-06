package common

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
