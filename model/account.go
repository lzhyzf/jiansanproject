package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ACCOUNT struct {
	AccountID, TransactionPrice, DisplayPrice  int       //账号id唯一索引， 交易价， 展示价
	HairStyle, Garments, Cloak, Mount, HangPet []string  //发型，衣服，披风，坐骑，挂宠
	Sect, Shape, Source, Faction, UID          string    //门派，体型，来源，阵营, 万宝楼id
	TimeUpdated                                time.Time //更新时间
}

// GetPrice 获取当前账号价格
func (a ACCOUNT) GetPrice() {
	fmt.Println(a.TransactionPrice)
}

// UpdateSumPrice 更新账号价格，把该账号下所有物品价格累加
func (a ACCOUNT) UpdateSumPrice() {
	fmt.Println(a.TransactionPrice)
}

// Sort 把账号下所有物品依照价格降序排序
func (a ACCOUNT) Sort() {
	fmt.Println("发型排序")
	//sort.Slice(a.HairStyle, func(i, j int) bool {
	//	return a[i] > a[j]
	//})
	fmt.Println("衣柜排序")

	fmt.Println("披风排序")
}

// Display 把账号下所有物品展示，按照发型，衣服，披风
func (a ACCOUNT) Display() {
	fmt.Println("----------------------------")

	fmt.Println("----------------------------")
}

// AddAccountInfo 往数据库中新增一个账号信息
func (a ACCOUNT) AddAccountInfo(client *mongo.Client) {
	//TODO 根据传入的万宝楼id，爬取该账号所有信息，再调用下面方法入库

	// 连接数据库并更新改物品的信息
	fmt.Println("新增一个账号")
	accountsCollection := client.Database("testing").Collection("accounts")
	insertOneIns := bson.M{
		"HairStyle": a.HairStyle, "Garments": a.Garments, "Cloak": a.Cloak,
		"Mount": a.Mount, "HangPet": a.HangPet, "DisplayPrice": a.DisplayPrice,
		"Sect": a.Sect, "Faction": a.Faction, "Shape": a.Shape, "TimeUpdated": a.TimeUpdated,
		"UID": a.UID,
	}
	// InsertOne
	insertOneResult, err := accountsCollection.InsertOne(context.TODO(), insertOneIns)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("_id:", insertOneResult.InsertedID)
}
