package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type ITEM struct { //名字是物品唯一索引
	Name                                                           string
	Alias                                                          string
	Type                                                           string
	MinimumPublishedPrice, MinimumOnSalePrice, TheHighestPriceEver int
	AppearanceSaleTime, TimeUpdated                                time.Time
}

// UpdateInfoFromDbByName 根据物品名称从数据库更新物品信息，如不存在该物品，则插入
func (i ITEM) UpdateInfoFromDbByName(name string, client *mongo.Client) {
	// 连接数据库并更新改物品的信息
	fmt.Printf("更新物品%s信息中", name)
	itemsCollection := client.Database("testing").Collection("items")
	//UpdateOne
	updateOneOpts := options.Update().SetUpsert(true)
	updateOneFilter := bson.M{"Name": name}
	updateOneSet := bson.M{"$set": bson.M{
		"Name": i.Name, "Alias": i.Alias, "Type": i.Type,
		"MinimumPublishedPrice": i.MinimumPublishedPrice, "MinimumOnSalePrice": i.MinimumOnSalePrice, "TheHighestPriceEver": i.TheHighestPriceEver,
		"AppearanceSaleTime": i.AppearanceSaleTime, "TimeUpdated": i.TimeUpdated,
	}}
	updateResult, err := itemsCollection.UpdateOne(context.TODO(), updateOneFilter, updateOneSet, updateOneOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"matched: %d  modified: %d  upserted: %d  upsertedID: %v\n",
		updateResult.MatchedCount,
		updateResult.ModifiedCount,
		updateResult.UpsertedCount,
		updateResult.UpsertedID,
	)
}
