package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"jiansan_go_project/database"
	"jiansan_go_project/model"
	"log"
)

// GetAccountInfoFromDbByID 根据账号id从数据库获取账号信息
func GetAccountInfoFromDbByID(id string) (*model.ACCOUNT, error) {
	// 连接数据库并获取账号的信息
	fmt.Printf("获取账号%s信息中", id)
	itemsCollection := database.DefaultDBClient.Database("testing").Collection("accounts")
	singleResult := itemsCollection.FindOne(context.TODO(), bson.M{"UID": id})
	var result bson.M
	if err := singleResult.Decode(&result); err != nil {
		log.Fatal(result)
	}
	i := &model.ACCOUNT{}
	data, err := bson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	err = bson.Unmarshal(data, i)
	if err != nil {
		log.Fatal(err)
	}
	return i, nil
}
