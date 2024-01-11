package utils

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	model2 "jiansan_go_project/model"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ConfMaxUsernameLen    = 30
	ConfMinPasswordLength = 8
)

type Response struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Data interface{} `json:"data"`

	// Reference returns the reference document which maybe useful to solve this error.
	Msg string `json:"msg"`
}

// tran 将字符串转化为go的时间格式
func tran(s string) (time.Time, error) {
	tms1 := strings.Split(s, " ")
	tms2 := strings.Split(tms1[0], "/")
	if len(tms2[1]) < 2 {
		tms2[1] = "0" + tms2[1]
	}
	if len(tms2[2]) < 2 {
		tms2[2] = "0" + tms2[2]
	}
	newstr := tms2[0] + "-" + tms2[1] + "-" + tms2[2] + " " + tms1[1]
	t, err := time.ParseInLocation("2006-01-02 15:04:05", newstr, time.Local)
	if err != nil {
		return t, err
	}
	return t, nil
}

// ReadCsvAndUpdateItemsInfo csv文件读取并更新物品表 https://blog.csdn.net/finishy/article/details/122391592
func ReadCsvAndUpdateItemsInfo(filepath string, client *mongo.Client) {
	var item model2.ITEM
	//打开文件(只读模式)，创建io.read接口实例
	opencast, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer opencast.Close()

	//创建csv读取接口实例
	r := csv.NewReader(opencast)

	for {
		// Read each record from csv
		record, readerr := r.Read()
		if readerr == io.EOF {
			break
		} else if readerr != nil {
			log.Fatal(readerr)
		}
		if record[3] != "" && record[3] != "出售时间" {
			item.TimeUpdated, err = tran(record[3])
			if err != nil {
				fmt.Printf("transfer err: %v, csv content: %v, %v, %v, %v, %v",
					err, record[0], record[1], record[2], record[3], record[4])
				continue
			}
			item.Name = record[0]
			item.Type = record[1]
			p, _ := strconv.Atoi(record[4])
			item.MinimumPublishedPrice = p
			item.UpdateInfoFromDbByName(item.Name, client)
		}
	}
}

// ReadCsvAndUpdateAccountsInfo csv文件读取并更新账号表
func ReadCsvAndUpdateAccountsInfo(filepath string, client *mongo.Client) {
	var account model2.ACCOUNT
	//打开文件(只读模式)，创建io.read接口实例
	opencast, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer opencast.Close()

	//创建csv读取接口实例
	r := csv.NewReader(opencast)

	colname := []string{
		"账号类型", "门派", "体型", "价格", "阵营", "战阶等级", "三内", "JJC段位22", "JJC段位33", "JJC段位55", "资历",
		"宠物分", "五甲数量", "积分", "充销", "点卡", "装分一", "装分二", "装分三", "装分四", "面挂数", "背挂数", "腰挂数",
		"左肩饰数量", "右肩饰数量", "成衣数", "披风数", "盒子数", "奇遇数", "奇珍数", "奇趣数量", "拓印数量", "武器数量",
		"易容捏脸数量", "金发数量", "红发数量", "黑发数量", "白发数量", "发型总数量", "佩囊数量", "挂宠数量", "称号名称",
		"挂件名称", "肩饰名称", "橙武玄晶名称", "挂宠佩囊名称", "坐骑奇趣名称", "奇遇奇珍名称", "发型名称", "成衣名称",
		"披风名称", "盒子名称", "拓印名",
	}
	i := 0
	for {
		// Read each record from csv
		record, readerr := r.Read()
		if readerr == io.EOF {
			break
		} else if readerr != nil {
			log.Fatal(readerr)
		}
		i++
		if i == 1 {
			continue
		}
		err = json.Unmarshal([]byte(record[returnIndex(colname, "发型名称")]), &account.HairStyle)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(record[returnIndex(colname, "成衣名称")]), &account.Garments)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(record[returnIndex(colname, "披风名称")]), &account.Cloak)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(record[returnIndex(colname, "坐骑奇趣名称")]), &account.Mount)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(record[returnIndex(colname, "挂宠佩囊名称")]), &account.HangPet)
		if err != nil {
			log.Fatal(err)
		}
		account.DisplayPrice = myAtoI(record[returnIndex(colname, "价格")])
		account.Sect = record[returnIndex(colname, "门派")]
		account.Faction = record[returnIndex(colname, "阵营")]
		account.Shape = record[returnIndex(colname, "体型")]
		// 先随机生成一个全局唯一的id，模拟万宝楼id
		account.UID = genXid()
		account.TimeUpdated = time.Now()

		account.AddAccountInfo(client)

	}
}

// returnIndex 返回给定字符串数组中某个字符串的索引位置
func returnIndex(sl []string, s string) int {
	for k, v := range sl {
		if v == s {
			return k
		}
	}
	return -1
}

// myAtoI asc to int
func myAtoI(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

// genXid 生成16位唯一字符串，用于account中的uid
func genXid() string {
	id := xid.New()
	return id.String()
}

// WriteResponseWithCode write an error or the response data into http response body.
func WriteResponseWithCode(c *gin.Context, msg string, data interface{}, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// GenMD5WithSalt 通过salt生成md5值
func GenMD5WithSalt(passwd, salt string) string {
	s := passwd + "::" + salt
	md5Hash := md5.New()
	md5Hash.Write([]byte(s))
	// 转16进制
	return hex.EncodeToString(md5Hash.Sum(nil))
}

// ToBson 将任意数据转化为bsonM的格式
func ToBson(r interface{}) bson.M {
	result := make(bson.M)
	v := reflect.ValueOf(r)
	t := reflect.TypeOf(r)

	for i := 0; i < v.NumField(); i++ {
		filed := v.Field(i)
		tag := t.Field(i).Tag
		key := tag.Get("bson")
		if key == "" || key == "-" || key == "_id" {
			continue
		}
		keys := strings.Split(key, ",")
		if len(keys) > 0 {
			key = keys[0]
		}
		// TODO: 处理字段嵌套问题
		switch filed.Kind() {
		case reflect.Int, reflect.Int64:
			v := filed.Int()
			if v != 0 {
				result[key] = v
			}
		case reflect.String:
			v := filed.String()
			if v != "" {
				result[key] = v
			}
		case reflect.Bool:
			result[key] = filed.Bool()
		case reflect.Ptr:

		case reflect.Float64:
			v := filed.Float()
			if v != 0 {
				result[key] = v
			}
		case reflect.Float32:
			v := filed.Float()
			if v != 0 {
				result[key] = v
			}
		default:
		}
	}
	return result
}

var nameMatch = regexp.MustCompile(`\A((@[^\s\/~'!\(\)\*]+?)[\/])?([^_.][^\s\/~'!\(\)\*]+)\z`)

// IsValidName 判断注册的名字是否合法
func IsValidName(name string) error {
	if strings.TrimSpace(name) != name {
		return errors.New("username contains space")
	}
	if len(name) == 0 || len(name) > ConfMaxUsernameLen {
		return errors.New("username invalid len")
	}
	if !nameMatch.MatchString(name) {
		return errors.New("username invalid pattern")
	}
	return nil
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+-/=?^_`{|}~]*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsValidEmail 判断注册的email是否合法
func IsValidEmail(email string) error {
	if len(email) == 0 {
		return errors.New("invalid email len")
	}

	if email[0] == '-' {
		return errors.New("invalid email")
	}

	n := strings.LastIndex(email, "@")
	if n <= 0 {
		return errors.New("invalid email address")
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email pattern")
	}

	return nil
}

// IsValidPasswd 判断密码是否合法
func IsValidPasswd(passwd string) error {
	if len(passwd) < ConfMinPasswordLength {
		return errors.New("password too short, should be more than 7")
	}
	if !IsComplexEnough(passwd) {
		return errors.New("password too simple")
	}

	return nil
}
