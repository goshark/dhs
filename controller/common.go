package controller

import (
	"gitee.com/goshark/dhs/model"
	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/net/ghttp"
	//"reflect"
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Success int         `json:success,omitempty`
	Msg     interface{} `json:msg,omitempty`
}

type Config struct {
	Ip     string      `json:"ip,omitempry"`     //服务器ip
	Dbuser string      `json:"dbuser,omitempry"` //数据库用户名
	Dbpass string      `json:"dbpass,omitempry"` //数据库密码
	Dbport interface{} `json:"dbport,omitempry"` //数据库连接端口
	Type   string      `json:"type,omitempty"`   //使用类型
	Uuid   string      `json:"uuid,omitempty"`   //服务器uuid
}

type Source struct {
	Category string //资源类别group,servers,user
	Uuid     string //资源uuid
	Intime   int64  //资源录入时间
}

//生成uuid
func generate_uuid() string { //uuid由 当前时间戳+random组成
	rand.Seed(time.Now().Unix()) //随机种子
	md := strconv.Itoa(rand.Intn(99)) + strconv.FormatInt(time.Now().Unix(), 10)
	data := []byte(md)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

//生成资源入库并返回uuid
func NewUuid(category string) string {
	uuid := generate_uuid()
	if ok := NewSource(category, uuid); ok { //资源入库
		return uuid
	} else {
		return ""
	}
}

//创建新资源
func NewSource(category, uuid string) bool {
	source := new(Source)
	source.Intime = time.Now().Unix()
	source.Category = category
	source.Uuid = uuid
	source_b, _ := gjson.New(source).ToJson()
	ok := model.NewDB("source").Data(source_b).WhereKey("source_" + strconv.Itoa(len(model.NewDB("source").GetKey())+1)).Insert()
	return ok
}

//cname 控制器路由字符串格式
func GetAction(r *ghttp.Request, CNAME string) (action string) {
	//获取当前action
	action_withsuffix := strings.Trim(strings.Replace(r.RequestURI, CNAME, "", 1), " ")
	pint := strings.Index(action_withsuffix, "?")
	if pint > 0 {
		action = string([]rune(action_withsuffix)[:pint])

	} else {
		action = action_withsuffix
	}
	return
}

//请求参数获取
func GetParams(c interface{}) (data *gjson.Json) {

	switch controller := c.(type) {
	case *ControllerUser:
		data = getparam(controller.Request, controller.Response)
	case *ControllerMaster:
		data = getparam(controller.Request, controller.Response)
	case *ControllerSlave:
		data = getparam(controller.Request, controller.Response)
	case *ControllerHome:
		data = getparam(controller.Request, controller.Response)
	case *ControllerSquid:
		data = getparam(controller.Request, controller.Response)
	case *ControllerGroup:
		data = getparam(controller.Request, controller.Response)
	case *ControllerBackup:
		data = getparam(controller.Request, controller.Response)
	case *ControllerTask:
		data = getparam(controller.Request, controller.Response)
	}
	return
}

func getparam(r *ghttp.Request, w *ghttp.Response) *gjson.Json {
	//	fmt.Println(r.Get("jsoninfo"))
	data, err := gjson.DecodeToJson([]byte(r.Get("jsoninfo")))
	if err != nil {
		return nil
	}
	return data
}

//判断文件或者文件是否存在
func IsExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err

}

func SendAjax(url string, params []byte, controller interface{}) (r *ghttp.ClientResponse, e error) {
	params_b := ghttp.BuildParams(map[string]string{
		"jsoninfo": string(params),
	})
	r, e = ghttp.NewClient().Post("http://"+GetParams(controller).GetString("ip")+":8888"+url, params_b)
	return
}
