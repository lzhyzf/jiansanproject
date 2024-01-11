package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"jiansan_go_project/database"
	"jiansan_go_project/model"
	"log"
)

// GetItemInfoFromDbByName 根据物品名称从数据库获取物品信息
func GetItemInfoFromDbByName(name string) (*model.ITEM, error) {
	// 连接数据库并获取物品的信息
	fmt.Printf("获取物品%s信息中", name)
	itemsCollection := database.DefaultDBClient.Database("testing").Collection("items")
	singleResult := itemsCollection.FindOne(context.TODO(), bson.M{"Name": name})
	var result bson.M
	if err := singleResult.Decode(&result); err != nil {
		log.Fatal(result)
	}
	i := &model.ITEM{}
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

//func GetAllUser() (userList []*entity.User, err error) {
//	if err := dao.SqlSession.Find(&userList).Error; err != nil {
//		return nil, err
//	}
//	return
//}
//
//func DeleteUserById(id string) (err error) {
//	err = dao.SqlSession.Where("id=?", id).Delete(&entity.User{}).Error
//	return
//}
//
//func GetUserById(id string) (user *entity.User, err error) {
//	if err = dao.SqlSession.Where("id=?", id).First(user).Error; err != nil {
//		return nil, err
//	}
//	return
//}
//
//func UpdateUser(user *entity.User) (err error) {
//	err = dao.SqlSession.Save(user).Error
//	return
//}
