package main

import (
	"fmt"
	"jiansan_go_project/database"
	"jiansan_go_project/routes"
	"log"
	"time"
)

// https://cloud.tencent.com/developer/article/1856734 之后要重构成改模，web应用mvc经典模式，为了方便以后扩展

// https://blog.csdn.net/qq_50737715/article/details/124335666 加入 中间件以及token 来进行登录态权限控制

// https://zhuanlan.zhihu.com/p/617512279 更高级的登录注册服务

//

func main() {
	bT := time.Now() // 开始时间
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.DBClose()
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r := routes.SetRouter()
	r.Run()
	//创建一个io对象
	//filename := "C:\\Users\\Administrator\\Desktop\\代售账号测试.csv"
	//ReadCsvAndUpdateAccountsInfo(filename, database.DefaultDBClient)

	eT := time.Since(bT) // 从开始到当前所消耗的时间
	fmt.Println("Run time: ", eT)

	//retrieve single and multiple documents with a specified filter using FindOne() and Find()
	//create a search filer
	//filter := bson.D{
	//	{"$and",
	//		bson.A{
	//			bson.D{
	//				{"age", bson.D{{"$gt", 25}}},
	//			},
	//		},
	//	},
	//}

	// retrieving the first document that match the filter
	//var result bson.M
	//// check for errors in the finding
	//if err = usersCollection.FindOne(context.TODO(), bson.D{}).Decode(&result); err != nil {
	//	panic(err)
	//}
	//
	//// display the document retrieved
	//fmt.Println("displaying the first result from the search filter")
	//fmt.Println(result)
}
