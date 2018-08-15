package controller

import (
	"database/sql"
	"fmt"
	"newproject/model"
	"strconv"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
)

type ControllerHome struct {
	gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/home", new(ControllerHome))
}

//重写父类构造函数
func (home *ControllerHome) Init(r *ghttp.Request) {

	defer home.Controller.Init(r)

	//身份验证层
	//获取action
	action := GetAction(r, "/home/")
	fmt.Println("当前请求:home", action)
	//	if r.Session.Get("user") == nil {
	//		r.Response.RedirectTo("login-index")
	//		r.Exit()
	//	}

}

//实例化表
func (home *ControllerHome) M(tables string) *model.Models {
	return model.NewDB(tables)
}

func (home *ControllerHome) Action() {
	var resp Response
	params := GetParams(home)
	//mos -> master or slave

	mos := params.GetString("type")
	fmt.Println("当前所请求的主从类型为:", mos)
	if mos == "" {
		home.Response.WriteJson(&Response{Success: 2, Msg: "请设置主从类别"})
		return
	} else if mos == "master" {
		if master_ip := GetMasterIp(); master_ip != "" {
			home.Response.WriteJson(&Response{Success: 2, Msg: "master机器已存在!ip为:" + master_ip})
			return
		}
	} else if mos == "slave" {
		masterip := params.GetString("masterip")

		condition, _ := gjson.Encode([]map[string]interface{}{{"name": "ip", "value": masterip, "op": "="}})
		cond, _ := gjson.DecodeToJson(condition)
		mres := home.M("config").Where(cond.ToArray()).Select()
		var masterinfo []byte
		for _, v := range mres {
			masterinfo = v
			break
		}
		masterinfo_j, _ := gjson.DecodeToJson(masterinfo)
		masterinfo_j.Set("slaveip", params.GetString("ip"))
		masterinfo_b, _ := masterinfo_j.ToJson()
		//先请求master服务器端,GrantUser，携带slaveip授权给slave端一个账号,
		r, e := ghttp.NewClient().Post("http://"+masterip+":8888/master/grant-user", ghttp.BuildParams(map[string]string{
			"jsoninfo": string(masterinfo_b),
		}))
		if e != nil {
			home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
			return
		}
		gjson.DecodeTo(r.ReadAll(), &resp)
		if resp.Success != 0 {
			home.Response.WriteJson(&Response{Success: resp.Success, Msg: resp.Msg})
			return
		}
		//获取master_status里两个主要的参数
		master_status := GetMasterStatus(masterip)

		//请求slave服务端
		params.Set("file", master_status["file"])
		params.Set("position", master_status["pos"])
		params_b, _ := params.ToJson()
		r, e = ghttp.NewClient().Post("http://"+params.GetString("ip")+":8888/slave/exec", ghttp.BuildParams(map[string]string{
			"jsoninfo": string(params_b),
		}))

		if e != nil {
			home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
			return
		}

		gjson.DecodeTo(r.ReadAll(), &resp)
		fmt.Println("103", resp)
		if resp.Success != 0 {
			home.Response.WriteJson(&Response{Success: 2, Msg: resp.Msg})
			return
		}

	}

	params_b, _ := params.ToJson()
	r, e := SendAjax("/"+mos+"/exec", params_b, home)
	if e != nil {
		home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
		return
	}

	gjson.DecodeTo(r.ReadAll(), &resp)
	fmt.Println("118", resp)
	if resp.Success != 0 {
		home.Response.WriteJson(&Response{Success: resp.Success, Msg: resp.Msg})
		return

	} else { //执行成功后

		if err := params.Set("type", mos); err != nil {
			fmt.Println("设置type错误")
			home.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		}
		conf := new(Config)

		condition, _ := gjson.Encode([]map[string]interface{}{{"name": "ip", "value": params.GetString("ip"), "op": "="}})
		cond, _ := gjson.DecodeToJson(condition)
		fmt.Println(cond.ToArray())
		result := home.M("config").Where(cond.ToArray()).Select()

		var conf_key string
		for k, v := range result {
			gjson.DecodeTo(v, &conf)
			conf_key = k
			break
		}

		conf.Type = mos
		conf_b, _ := gjson.Encode(&conf)
		if home.M("config").Data(conf_b).WhereKey(conf_key).Update() {
			home.Response.WriteJson(&Response{Success: 0, Msg: "设置成功！"})
		} else {
			home.Response.WriteJson(&Response{Success: 2, Msg: "设置失败！"})
		}
	}

}
func (home *ControllerHome) AppendList() {
	var res Response
	params := GetParams(home)
	params_b, err := params.ToJson()
	if err != nil {
		fmt.Println(err)
	}
	r, e := SendAjax("/home/check", params_b, home)
	if e != nil {
		home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
		return
	}

	gjson.DecodeTo(r.ReadAll(), &res)
	fmt.Println("line-168", res)
	if res.Success != 0 {
		home.Response.WriteJson(&Response{Success: res.Success, Msg: res.Msg})
		return
	}

	//检查是该机器是否存在列表中
	fmt.Println(params.GetString("ip"))
	result := home.M("config").Where(gjson.New([]map[string]interface{}{{"name": "ip", "value": params.GetString("ip"), "op": "="}}).ToArray())

	var conf_key string
	if len(result.Select()) > 0 {
		home.Response.WriteJson(&Response{Success: 2, Msg: "添加失败,列表中已存在该服务器信息!ip:" + params.GetString("ip")})
		return
	}
	conf_key = "conf_" + strconv.Itoa(len(home.M("config").GetKey())+1)

	ok := home.M("config").Data(params_b).WhereKey(conf_key).Insert()
	if ok {
		home.Response.WriteJson(&Response{Success: 0, Msg: "添加成功!"})
	} else {
		home.Response.WriteJson(&Response{Success: 2, Msg: "添加失败!"})
	}

}

func (home *ControllerHome) Check() {

	conf := new(Config)
	params, _ := GetParams(home).ToJson()
	gjson.DecodeTo(params, &conf)

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", conf.Dbuser, conf.Dbpass, DBNET, DBSERVER, conf.Dbport, DATABASE)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		home.Response.WriteJson(&Response{Success: 2, Msg: "添加失败！" + err.Error()})
		return
	}
	if err := DB.Ping(); err != nil {
		fmt.Println(err)
		home.Response.WriteJson(&Response{Success: 2, Msg: "添加失败！" + err.Error()})
		return
	}
	home.Response.WriteJson(&Response{Success: 0, Msg: "验证成功！"})
}

//获取主从配置列表
func (home *ControllerHome) Configlist() {
	tb_config := home.M("config")
	var result []interface{}
	data := GetParams(home)
	if data == nil {
		home.Response.WriteJson(Response{Success: 2, Msg: "没有获取到请求数据"})
		return
	}

	res := tb_config.Where(data.GetArray("condition")).Select() // data.GetArray("condition") 返回 []interface{}
	if len(res) <= 0 {
		home.Response.WriteJson(&Response{Success: 0, Msg: ""})
		return
	}
	for _, v := range res {
		b, _ := gjson.DecodeToJson(v)
		result = append(result, b.ToMap())
	}
	home.Response.WriteJson(&Response{Success: 0, Msg: result})

}
