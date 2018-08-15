package controller

import (
	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/net/ghttp"
	//"reflect"
	//"fmt"
	"os"
	"strings"
)

type Response struct {
	Success int         `json:success,omitempty`
	Msg     interface{} `json:msg,omitempty`
}

type Config struct {
	Ip       string      `json:"ip,omitempry"`
	Dbuser   string      `json:"dbuser,omitempry"`
	Dbpass   string      `json:"dbpass,omitempry"`
	Dbport   interface{} `json:"dbport,omitempry"`
	Slaveips []string    `json:"slaveips,omitempty"`
	Type     string      `json:"type,omitempty"`
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
		data = GetParam(controller.Request, controller.Response)
	case *ControllerMaster:
		data = GetParam(controller.Request, controller.Response)
	case *ControllerSlave:
		data = GetParam(controller.Request, controller.Response)
	case *ControllerHome:
		data = GetParam(controller.Request, controller.Response)
	case *ControllerSquid:
		data = GetParam(controller.Request, controller.Response)
	}
	return
}

func GetParam(r *ghttp.Request, w *ghttp.Response) *gjson.Json {
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
