一个奇奇怪怪的剑三项目。
以下随缘乱写

明确需求 ->  拆解需求

大的模块

1. 拉钩模块
1.1 无登录态的基于万宝楼id的拉钩方式
1.1.1前端设计
    

1.1.2 后端设计
1.1.2.1数据库
数据库只需要一张联系方式表 （万宝楼id ，联系方式）
1.1.2.2 功能
·基于万宝楼id获取联系方式
Get TICK id=xxx
RES {
Status  200/400
Contact  135xxxxxx/email
}
基于万宝楼id更新（只增加？联系方式）
POST TICK id=xxx
RES {
Status 200/400
}

1.1.3 宣传/运营
前期可能需要爬贴吧先收集部分联系方式
然后向所有人宣传平台的拉钩模块，方便买卖双方快速找号
1.2 有登录态的基于万宝楼id的拉钩方式

1.3 有登录态的所有账号（挂在万宝楼账号 + 不挂在万宝楼）的拉钩方式

2. 估价模块
2.1 比价系统
2.1.1前端设计
    
2.1.2 后端设计
2.1.2.1数据库 Untitled - dbdiagram.io
账号表（门派，体型，发型，衣服，披风，价格，更新时间）
物品表（名字，别名，类型[发型/衣服/披风]，最小公示价，最小在售价，最高历史成交价，更新时间）
交易表（交易地点，类型[账号/物品/金币]，价格，更新时间）
 

2.1.2.2功能
	·通过万宝楼更新物品表数据
	24小时对所有物品进行一次更新

·基于万宝楼账号id返回自身账号物品价格查询信息
Get ACCOUNT id=xxx
RES {
Status  200/400
	Object {
	[]发型{   //返回所有价格（前端默认展示5件以及…）
		发型名: {
	最小在售价
	查询时间
} …
},
	[]衣服{  //返回所有价格（前端默认展示5件以及…）
		衣服名: {
	最小在售价
	查询时间
} …
},
[]披风{  //返回所有价格（前端默认展示5件以及…）
		披风名: {
	最小在售价
	查询时间
} …
},
}
}
			
2.1.3 宣传/运营
前期可能需要收集所有物品的当前价格
去剑三大群/贴吧宣传可以一键获取自己账号物品价值

2.2估价系统
2.2.1前端设计

2.2.2 后端设计
			2.2.2.1数据部分
				·数据样式
			
			2.2.2.3 功能
				·基于机器学习的剑三账号估价模型，可参考零基础入门数据挖掘 - 二手车交易价格预测_学习赛_相关的问题_天池大赛-阿里云天池的论坛 (aliyun.com)
                 输入万宝楼id，输出价格
Get VALUATION id=xxx
RES {
Status  200/400
Price  ￥8888
}

1.清洗重复数据
2.

2.2.3 宣传/运营
向所有人宣传免费估价









































