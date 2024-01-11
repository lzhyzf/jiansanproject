package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var DefaultDBClient *mongo.Client

// ClientOpts mongoClient 连接客户端参数
var ClientOpts = options.Client().
	//SetAuth(options.Credential{
	//	AuthMechanism: "",
	//	AuthSource:    "",
	//	Username:      "",
	//	Password:      "",
	//}).
	SetConnectTimeout(10 * time.Second).
	SetHosts([]string{"localhost:27017"}).
	SetMaxPoolSize(20).
	SetMinPoolSize(5).
	SetReadPreference(readpref.Primary()).
	SetReplicaSet("")

// https://juejin.cn/post/6908063164726771719 参考这篇文章即可

func InitDB() (err error) {
	DefaultDBClient, err = mongo.Connect(context.TODO(), ClientOpts)
	if err != nil {
		log.Fatal(err)
	}

	//验证数据库连接是否成功，若成功，则无异常
	return DefaultDBClient.Ping(context.TODO(), readpref.Primary())
}

func DBClose() {
	DefaultDBClient.Disconnect(context.TODO())
}
