package controller

import (
	"fmt"
	"newproject/model"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
)

type ControllerUser struct {
	gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/user", new(ControllerUser))
}

//重写父类构造函数
func (user *ControllerUser) Init(r *ghttp.Request) {

	defer user.Controller.Init(r)

	//身份验证层
	//获取action
	action := GetAction(r, "/user/")
	fmt.Println("当前控制器方法:/user/", action)
	//	if action != "login-index" && action != "login" && action != "register" && action != "register-index" { //登录和注册不存在token验证
	//		if r.Session.Get("user") == nil {
	//			r.Response.RedirectTo("login-index")
	//			r.Exit()
	//		}

	//	}
}

//实例化表
func (user *ControllerUser) M(tables string) *model.Models {
	return model.NewDB(tables)
}

//用户注册
func (user *ControllerUser) Register() {
	var user_exist bool = false //用户是否存在变量初始化,默认不存在该用户
	tb_user := user.M("user")
	data := GetParams(user)          //获取请求数据已*gjson.json数据格式接收
	username := data.Get("username") //返回interface{}类型,获取某一json字段的值，取值采用“.”递进层级的方式,遇到源数据结构字段带"."，设置 json.SetViolenceCheck(true)
	data_json := data.GetString("")  //获取请求的完整json数据转为[]byte类型
	if username == nil {             //非空判断
		user.Response.WriteJson(Response{Success: 2, Msg: "用户名不能为空！"})
		return
	}

	//判断用户名是否存在
	for _, v := range tb_user.GetKey() {
		if username.(string) == v {
			user_exist = true
			break
		}
	}
	if user_exist {
		user.Response.WriteJson(Response{Success: 2, Msg: "该用户已经存在!"})
		return
	}

	//写入数据
	ok := tb_user.WhereKey(username.(string)).Data([]byte(data_json)).Update()
	if ok {
		user.Response.WriteJson(Response{Success: 0, Msg: "注册成功!"})
	} else {
		user.Response.WriteJson(Response{Success: 2, Msg: "注册失败!"})
	}

}

func (user *ControllerUser) LoginIndex() {
	user.View.Display("user/login.html")
}
func (user *ControllerUser) RegisterIndex() {
	user.View.Display("user/register.html")
}

func (user *ControllerUser) Login() {
	tb_user := user.M("user")
	data := GetParams(user)          //获取请求数据已*gjson.json数据格式接收
	username := data.Get("username") //返回 interface{}类型，参数支持点递进寻层

	if username == nil || data.Get("password") == nil { //非空判断
		user.Response.WriteJson(Response{Success: 2, Msg: "用户名或密码不能为空！"})
		return
	}
	users := tb_user.WhereKey(username.(string)).Find()
	if users == nil {
		user.Response.WriteJson(Response{Success: 2, Msg: "没有这个用户"})
		return
	}
	userinfo, err := gjson.DecodeToJson(users)
	if err != nil {
		fmt.Println(err)
		return
	}

	if userinfo.Get("username") == username && userinfo.Get("password") == data.Get("password") {
		user.Session.Set("user", username.(string))
		user.Response.WriteJson(Response{Success: 0, Msg: "登录成功!"})
	} else {
		user.Response.WriteJson(Response{Success: 2, Msg: "用户名或密码错误！"})
	}

}

func (user *ControllerUser) Resouce() {
	data := GetParams(user)
	fmt.Println(data.Get())
	user.Response.WriteJsonP(Response{Success: 0, Msg: "give you some color to see see"})
}
func (user *ControllerUser) Server() {

	user.View.Display("listpage/server.html")
}
func (user *ControllerUser) Logout() {

	user.Session.Set("user", nil)
	user.Response.WriteJson(&Response{Success: 0, Msg: "注销成功！"})
}
func (user *ControllerUser) QueryUser() {
	/**********条件查询开始*****************/
	tb_user := user.M("user")
	var result map[string]interface{} = make(map[string]interface{})
	data := GetParams(user)
	if data == nil {
		user.Response.WriteJson(Response{Success: 2, Msg: "没有获取到请求数据"})
		return
	}

	fmt.Println(data.GetArray("condition"))
	res := tb_user.Where(data.GetArray("condition")).Select() // data.GetArray("condition") 返回 []interface{}
	for k, v := range res {
		b, _ := gjson.DecodeToJson(v)
		result[k] = b.ToMap()
	}
	user.Response.WriteJson(Response{Success: 0, Msg: result})

	/**********条件查询结束*****************/
	// if ok := tb_user.WhereKey("asd3").Delete(); ok {
	// 	fmt.Println("删除成功!")
	// 	for k, v := range tb_user.SelectAll() {
	// 		fmt.Println(k, string(v))
	// 	}
	// } else {
	// 	fmt.Println("删除失败!")
	// }
	/**********全部查询*****************/
	//下面是查所有，未使用条件筛选
	//res := user.M("user").SelectAll()
	// v1, _ := gjson.Encode(res)
	// fmt.Println(string(v1))
	// for k, v := range res {

	// 	fmt.Println(k, string(v))
	// }

	/**********插入数据*****************/

}
