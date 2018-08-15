package controller

import (
	"fmt"
	"io/ioutil"
	"os"

	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
	_ "github.com/Unknwon/goconfig"
)

type ControllerSquid struct {
	gmvc.Controller
}

const (
	hosts  string = "/etc/hosts"
	config string = "/usr/local/squid/etc/squid.conf"
)

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/squid", &ControllerSquid{})
}

//查询
func (s *ControllerSquid) QueryHosts() {
	var path string
	params := GetParams(s)
	if params == nil {
		s.Response.WriteJson(&Response{Success: 2, Msg: "未获取到请求参数！"})
		return
	}
	name := params.GetString("name")
	if name == "" {
		s.Response.WriteJson(&Response{Success: 2, Msg: ""})
		return
	}
	switch name {
	case "hosts":
		path = hosts
	case "config":
		path = config
	}
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("出错")
		s.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	}
	arr, _ := ioutil.ReadAll(f)
	s.Response.WriteJson(&Response{Success: 0, Msg: string(arr)})
}

//修改
func (s *ControllerSquid) Userform() {
	var path string
	params := GetParams(s)
	name := params.GetString("name")
	txt := params.GetString("txt")
	if name == "" {
		s.Response.WriteJson(&Response{Success: 2, Msg: ""})
		return
	}
	if txt == "" {
		s.Response.WriteJson(&Response{Success: 2, Msg: "文件不能修改为空！"})
		return
	}
	switch name {
	case "hosts":
		path = hosts
	case "config":
		path = config
	}
	if path == "" {
		s.Response.WriteJson(&Response{Success: 2, Msg: "修改失败!"})
		return
	}

	d1 := []byte(txt)
	err := ioutil.WriteFile(path, d1, 0644)
	s.check(err)

	s.Response.WriteJson(&Response{Success: 0, Msg: "修改成功!"})
}

func (s *ControllerSquid) check(e error) {
	if e != nil {
		s.Response.WriteJson(&Response{Success: 2, Msg: "修改失败!"})
		s.Exit()
		return
	}
}
