package controller

import (
	"database/sql"
	"fmt"
	"strconv"

	"gitee.com/goshark/dhs/model"

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

//执行配置
func (home *ControllerHome) Action() {
	var resp Response
	params := GetParams(home)
	//mos -> master or slave or master-master

	mos := params.GetString("type")
	fmt.Println("当前所请求的类型为:", mos)
	if mos == "" {
		home.Response.WriteJson(&Response{Success: 2, Msg: "未获取到执行对象"})
		return
	} else if mos == "master" { //主
		if master_ip := GetMasterIp(); master_ip != "" {
			home.Response.WriteJson(&Response{Success: 2, Msg: "master机器已存在!ip为:" + master_ip})
			return
		}
		params_b, _ := params.ToJson()
		r, e := SendAjax("/"+mos+"/exec", params_b, home)
		if e != nil {
			home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
			return
		}

		gjson.DecodeTo(r.ReadAll(), &resp)
		if resp.Success != 0 {
			home.Response.WriteJson(&Response{Success: resp.Success, Msg: resp.Msg})
			return

		}
	} else if mos == "slave" { //从
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
		if resp.Success != 0 {
			home.Response.WriteJson(&Response{Success: 2, Msg: resp.Msg})
			return
		}

	} else if mos == "master-master" { //主主

		var execchan = make(chan int)
		for k, v := range params.GetJson("list").ToArray() {
			fmt.Println("第", k+1, "次执行")
			go func(mv interface{}) {
				master_params, _ := gjson.New(mv).ToJson()
				r, e := ghttp.NewClient().Post("http://"+gjson.New(mv).GetString("ip")+":8888/master/exec", ghttp.BuildParams(map[string]string{
					"jsoninfo": string(master_params),
				}))
				if e != nil {
					execchan <- 2
					fmt.Println("line-133", e.Error())
					return
				}
				gjson.DecodeTo(r.ReadAll(), &resp)
				if resp.Success != 0 {
					execchan <- resp.Success
					fmt.Println("line-139", resp.Msg)
					return
				} else {
					r, e = ghttp.NewClient().Post("http://"+gjson.New(mv).GetString("ip")+":8888/master/grant-user", ghttp.BuildParams(map[string]string{
						"jsoninfo": string(master_params),
					}))
					if e != nil {
						execchan <- 2
						fmt.Println("line-147", e.Error())
						return
					}
					gjson.DecodeTo(r.ReadAll(), &resp)
					if resp.Success != 0 {
						execchan <- 2
						fmt.Println("line-153", e.Error())
						return
					}

					execchan <- resp.Success
					return
				}

			}(v)

		}

		var execsuccess int
		var i int = 1
		for {
			select {
			case c1 := <-execchan:
				i++
				execsuccess = c1

			default:
			}
			if execsuccess != 0 || i == 2 {
				break
			}
		}
		if execsuccess != 0 {
			home.Response.WriteJson(&Response{Success: execsuccess, Msg: "执行时出错!"})
			return
		}
		fmt.Println("双主设置完成！进行互从配置")
		//双主设置完成！进行互从配置
		for k, v := range params.GetJson("list").ToArray() {
			fmt.Println("第", k+1, "次执行互从配置")
			//获取master_status里两个主要的参数
			slave_params := gjson.New(v)
			slave_params.Set("masterip", slave_params.GetString("slaveip")) //重新构造参数，互为主从
			slave_params.Set("type", "slave")                               //模拟执行action；type:slave接口
			fmt.Println("获取masterstatus参数:", slave_params.GetString("masterip"))
			master_status := GetMasterStatus(slave_params.GetString("masterip"))
			slave_params.Set("file", master_status["file"])
			slave_params.Set("position", master_status["pos"])
			slave_params_b, _ := slave_params.ToJson()
			r, e := ghttp.NewClient().Post("http://"+slave_params.GetString("ip")+":8888/slave/exec", ghttp.BuildParams(map[string]string{
				"jsoninfo": string(slave_params_b),
			}))
			fmt.Println("routinezhixing参数", string(slave_params_b))
			if e != nil {
				fmt.Println("line-202", e.Error())
				break
			}

			gjson.DecodeTo(r.ReadAll(), &resp)
			if resp.Success != 0 {
				fmt.Println("209", resp)
				break
			}

		}

		//构造初始参数进行数据入库
		for _, v := range params.GetJson("list").ToArray() {
			params_n := gjson.New(v)
			params_n.Set("type", "master-master")
			params_n.Remove("file")
			params_n.Remove("position")
			params_n.Remove("masterip")
			if res_bool, err := home.updatedb(params_n); res_bool {
				resp.Success = 0
				resp.Msg = "设置成功！"
			} else {
				resp.Success = 2
				resp.Msg = "服务器(" + gjson.New(v).GetString("ip") + ")设置失败！" + err.Error()
				fmt.Println("line-243", "设置失败！")
				break
			}

		}
		home.Response.WriteJson(resp)
		return
	}

	//构造初始参数进行数据入库
	if res_bool, err := home.updatedb(params); res_bool {
		resp.Success = 0
		resp.Msg = "设置成功！"
	} else {
		resp.Success = 2
		resp.Msg = "服务器(" + params.GetString("ip") + ")设置失败！" + err.Error()
		fmt.Println("line-260", "设置失败！")
	}
	home.Response.WriteJson(resp)

}

//更新数据
func (home *ControllerHome) updatedb(params *gjson.Json) (bool, error) {
	fmt.Println("构造参数进行数据入库")
	fmt.Println(params.ToMap())
	mos := params.GetString("type")
	if err := params.Set("type", mos); err != nil {
		fmt.Println("设置type错误")
		return false, err

	}
	conf := new(Config)
	result := home.M("config").Where(gjson.New([]map[string]interface{}{{"name": "ip", "value": params.GetString("ip"), "op": "="}}).ToArray()).Select()
	var conf_key string
	for k, v := range result {
		gjson.DecodeTo(v, &conf)
		conf_key = k
		break
	}

	conf.Type = mos
	conf_b, _ := gjson.Encode(&conf)
	if home.M("config").Data(conf_b).WhereKey(conf_key).Update() {
		return true, nil
	} else {
		return false, nil
	}
}

//添加服务器到列表
func (home *ControllerHome) AppendList() {
	var res Response
	params := GetParams(home)
	params.Set("uuid", NewUuid("server"))
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

//删除列表中的服务器
func (home *ControllerHome) RemoveList() {
	params := GetParams(home)
	var conf_key string
	confinfo := home.M("config").Where(gjson.New([]map[string]interface{}{{"name": "ip", "value": params.GetString("ip"), "op": "="}}).ToArray()).Select()
	if len(confinfo) > 0 {
		for k, _ := range confinfo {
			conf_key = k
			break
		}
		if home.M("config").WhereKey(conf_key).Delete() {
			home.Response.WriteJson(&Response{Success: 0, Msg: "删除成功！"})
		} else {
			home.Response.WriteJson(&Response{Success: 2, Msg: "删除失败！"})
		}
	} else {
		home.Response.WriteJson(&Response{Success: 2, Msg: "未匹配到该信息！"})
	}
}

/***
statement:验证mysql连通性
params:{jsoninfo:ip:,dbuser:,dbpass,dbport:,}
***/

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

/***
statement:获取主从配置列表
params:{jsoninfo:condition:[],}
***/
//
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

/***
statement:获取主数据库列表
params:{jsoninfo:ip:,dbuser:,dbpass,dbport:,}
***/
//get master dblist
func (home *ControllerHome) GetDbList() {
	var res Response
	params_b, _ := GetParams(home).ToJson()
	r, e := SendAjax("/master/get-db-list", params_b, home)

	if e != nil {
		home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
		return
	}
	gjson.DecodeTo(r.ReadAll(), &res)
	home.Response.WriteJson(res)

}

/***
statement:忽略记录日志文件的数据库
params:{jsoninfo:ip:,dbuser:,dbpass,dbport:,dblist:[],}
***/
//
func (home *ControllerHome) IgnoreDb() {
	var res Response
	params_b, _ := GetParams(home).ToJson()
	r, e := SendAjax("/master/ignore-db", params_b, home)

	if e != nil {
		home.Response.WriteJson(&Response{Success: 2, Msg: e.Error()})
		return
	}
	gjson.DecodeTo(r.ReadAll(), &res)
	home.Response.WriteJson(res)
}
