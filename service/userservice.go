package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"jiansan_go_project/database"
	"jiansan_go_project/model"
	"jiansan_go_project/utils"
	"log"
	"time"
)

// GetUserFromDbByName 根据用户名称从数据库获取用户信息
func GetUserFromDbByName(name string) (*model.USER, error) {
	// 连接数据库并获取物品的信息
	fmt.Printf("获取用户%s信息中", name)
	itemsCollection := database.DefaultDBClient.Database("testing").Collection("users")
	singleResult := itemsCollection.FindOne(context.TODO(), bson.M{"Name": name})
	var result bson.M
	if err := singleResult.Decode(&result); err != nil {
		log.Fatal(result)
	}
	i := &model.USER{}
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

// UpdateUserToDB 更新用户信息
func UpdateUserToDB(username string, updates map[string]interface{}) error {
	fmt.Printf("更新用户%s信息", username)
	updates["UpdatedUnix"] = time.Now().Unix()

	usersCollection := database.DefaultDBClient.Database("testing").Collection("users")
	//UpdateOne
	updateOneFilter := bson.M{"Name": username}
	updateOneSet := bson.M{"$set": utils.ToBson(updates)}
	updateResult, err := usersCollection.UpdateOne(context.TODO(), updateOneFilter, updateOneSet)
	if err != nil {
		return err
	}
	fmt.Printf("updateResult:%v", updateResult)
	return nil
}

// UserExists 判断用户是否存在
func UserExists(field, val string) bool {
	fmt.Printf("根据[%s]:[%s]判断用户是否存在", field, val)
	user := &model.USER{}
	itemsCollection := database.DefaultDBClient.Database("testing").Collection("users")
	singleResult := itemsCollection.FindOne(context.TODO(), bson.M{field: val})
	if singleResult.Err() != nil || user.Uid == 0 {
		fmt.Println("user logic UserExists error:", singleResult.Err())
		return false
	}
	return true
}

// AddUserToDB 新建一个用户并插入数据库
func AddUserToDB(username, password, email string) error {
	newUser := &model.USER{}
	newUser.Email = email
	newUser.Name = username
	newUser.Salt = utils.GenRandomSalt()
	newUser.Passwd = utils.GenMD5WithSalt(password, newUser.Salt)
	newUser.CreatedUnix = time.Now().Unix()
	newUser.UpdatedUnix = time.Now().Unix()
	newUser.IsAdmin = false
	newUser.ProhibitLogin = false
	return newUser.ADDToDB(database.DefaultDBClient)
}
