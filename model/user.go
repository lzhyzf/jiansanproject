package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type USER struct { //名字是物品唯一索引
	Uid    int64  `bson:"Uid,omitempty"`    //主键，UID
	Name   string `bson:"Name,omitempty"`   // 用户名，唯一键
	Email  string `bson:"Email,omitempty"`  // 邮箱，唯一键
	Passwd string `bson:"Passwd,omitempty"` //已使用Salt进行加密的密码串MD5(原始password+Salt)

	Salt string `bson:"Salt,omitempty"` // 用于密码加盐哈希，注册时随机生成

	CreatedUnix   int64 `bson:"CreatedUnix,omitempty"`   // 账号创建时间
	UpdatedUnix   int64 `bson:"UpdatedUnix,omitempty"`   // 账号更新时间
	LastLoginUnix int64 `bson:"LastLoginUnix,omitempty"` // 账号上次登录时间

	IsAdmin bool `bson:"IsAdmin,omitempty"` // 管理员标记

	ProhibitLogin bool `bson:"ProhibitLogin,omitempty"` // 禁止登录标记

	LastLoginIp string `bson:"LastLoginIp,omitempty"` // 登录Ip
}

func (u *USER) ADDToDB(client *mongo.Client) error {
	// 连接数据库并更新改物品的信息
	fmt.Println("新增一个用户")
	usersCollection := client.Database("testing").Collection("users")
	insertOneIns := bson.M{
		"Uid": u.Uid, "Name": u.Name, "Email": u.Email, "Passwd": u.Passwd,
		"Salt": u.Salt, "CreatedUnix": u.CreatedUnix, "UpdatedUnix": u.UpdatedUnix, "LastLoginUnix": u.LastLoginUnix,
		"IsAdmin": u.IsAdmin, "ProhibitLogin": u.ProhibitLogin, "LastLoginIp": u.LastLoginIp,
	}
	insertOneResult, err := usersCollection.InsertOne(context.TODO(), insertOneIns)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("_id:", insertOneResult.InsertedID)
	return nil
}
