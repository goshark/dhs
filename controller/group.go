package controller

import (
	_ "database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gitee.com/goshark/dhs/model"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
)

type ControllerGroup struct {
	gmvc.Controller
}

type Group struct {
	Intime     int64    `json:"intime"`               //录入时间
	Attr       string   `json:"attr,omitempty"`       //组的属性
	Configlist []Config `json:"configlist,omitempty"` //组内服务器列表
	Uuid       string   `json:"uuid,omitempty"`       //组uuid
	Deleted    bool     `json:"deleted"`              //是否删除(软删除)
	Name       string   `json:"name,omitempty"`       //组名
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/group", new(ControllerGroup))
}

//重写父类构造函数
func (group *ControllerGroup) Init(r *ghttp.Request) {

	defer group.Controller.Init(r)

	//身份验证层
	//获取action
	action := GetAction(r, "/group/")
	fmt.Println("当前请求:group", action)
	//	if r.Session.Get("user") == nil {
	//		r.Response.RedirectTo("login-index")
	//		r.Exit()
	//	}

}

//实例化表
func (g *ControllerGroup) M(tables string) *model.Models {
	return model.NewDB(tables)
}

//创建Group
func (g *ControllerGroup) CreateGroup() {
	params, _ := GetParams(g).ToJson()
	group := new(Group)
	gjson.DecodeTo(params, &group)
	group.Uuid = NewUuid("group")
	group.Intime = time.Now().Unix()
	group_b, _ := gjson.New(group).ToJson()
	//在此解析configlist执行home/action相关操作。。
	group_key := "group_" + strconv.Itoa(len(g.M("group").GetKey())+1)
	if ok := g.M("group").Data(group_b).WhereKey(group_key).Insert(); ok {
		g.Response.WriteJson(&Response{Success: 0, Msg: "创建成功！"})
	} else {
		g.Response.WriteJson(&Response{Success: 2, Msg: "创建失败！"})
	}
}

//查询组
func (g *ControllerGroup) SelectGroup() {
	params := GetParams(g)
	var result []interface{}
	res := g.M("group").Where(params.GetArray("condition")).Select()
	if len(res) <= 0 {
		g.Response.WriteJson(&Response{Success: 0, Msg: ""})
		return
	}
	for _, v := range res {
		b, _ := gjson.DecodeToJson(v)
		result = append(result, b.ToMap())
	}
	g.Response.WriteJson(&Response{Success: 0, Msg: result})
}

//修改组
func (g *ControllerGroup) UpdateGroup() {
	var g_key string
	var group Group
	var group_old Group

	params_b, _ := GetParams(g).ToJson()
	gjson.DecodeTo(params_b, &group)
	//v_group := reflect.ValueOf(group)
	g_value := g.M("group").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": group.Uuid, "op": "="}}).ToArray()).Select()

	for k, v := range g_value {
		g_key = k
		gjson.DecodeTo(v, &group_old)
		break
	}

	//v_group_old := reflect.ValueOf(group_old)

	group_old_elem := reflect.ValueOf(&group_old).Elem() //返回 group_old 指针保管的值
	//	for i := 0; i < v_group.NumField(); i++ {
	//		field_group := v_group.Field(i)
	//		field_group_old := v_group_old.Field(i)
	//		if !reflect.DeepEqual(field_group.Interface(), reflect.Zero(field_group.Type()).Interface()) {
	//			if !reflect.DeepEqual(field_group.Interface(), field_group_old.Interface()) { //替换原数据
	//				if group_old_elem.Field(i).CanSet() != true {
	//					fmt.Println("can not set value", group_old_elem.Field(i).String())
	//					continue
	//				}
	//				group_old_elem.Field(i).Set(field_group) //set value
	//			}
	//		}
	//	}
	Merge(group_old, group, group_old_elem)
	group_b, err := gjson.New(group_old).ToJson()

	if err != nil {
		fmt.Println(err.Error())
	}

	if g.M("group").Data(group_b).WhereKey(g_key).Update() {
		g.Response.WriteJson(&Response{Success: 0, Msg: "修改成功"})
	} else {
		g.Response.WriteJson(&Response{Success: 2, Msg: "修改失败"})
	}
}

func Merge(v1, v2 interface{}, v1_e reflect.Value) {
	v1_v := reflect.ValueOf(v1) //old
	v2_v := reflect.ValueOf(v2) //new

	for i := 0; i < v2_v.NumField(); i++ {
		v2_f := v2_v.Field(i)
		v1_f := v1_v.Field(i)

		if !reflect.DeepEqual(v2_f.Interface(), reflect.Zero(v2_f.Type()).Interface()) || v1_v.Type().Field(i).Name == "Deleted" {

			if !reflect.DeepEqual(v2_f.Interface(), v1_f.Interface()) { //替换原数据
				if v1_e.Field(i).CanSet() != true {
					fmt.Println("can not set value", v1_e.Field(i).String())
					continue
				}
				v1_e.Field(i).Set(v2_f) //set value
			}
		}
	}
}

//删除组（使用软删除，到回收站）
func (g *ControllerGroup) DeleteGroup() {
	var g_key string
	params := GetParams(g)
	group := new(Group)
	g_value := g.M("group").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": params.GetString("uuid"), "op": "="}}).ToArray()).Select()

	for k, v := range g_value {
		g_key = k
		gjson.DecodeTo(v, &group)
		break
	}
	group.Deleted = true
	group_b, _ := gjson.New(group).ToJson()
	if g.M("group").Data(group_b).WhereKey(g_key).Update() {
		g.Response.WriteJson(&Response{Success: 0, Msg: "删除成功"})
	} else {
		g.Response.WriteJson(&Response{Success: 2, Msg: "删除失败"})
	}
}

//从回收站恢复
func (g *ControllerGroup) RecoverGroup() {
	var g_key string
	params := GetParams(g)
	group := new(Group)
	g_value := g.M("group").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": params.GetString("uuid"), "op": "="}}).ToArray()).Select()

	for k, v := range g_value {
		g_key = k
		gjson.DecodeTo(v, &group)
		break
	}
	group.Deleted = false
	group_b, _ := gjson.New(group).ToJson()
	if g.M("group").Data(group_b).WhereKey(g_key).Update() {
		g.Response.WriteJson(&Response{Success: 0, Msg: "恢复成功"})
	} else {
		g.Response.WriteJson(&Response{Success: 2, Msg: "恢复失败"})
	}
}

//彻底删除(直接从数据库中删除组)
func (g *ControllerGroup) RemoveGroup() {
	var g_key string
	params := GetParams(g)
	group := new(Group)
	g_value := g.M("group").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": params.GetString("uuid"), "op": "="}}).ToArray()).Select()

	for k, v := range g_value {
		g_key = k
		gjson.DecodeTo(v, &group)
		break
	}

	if !group.Deleted {
		g.Response.WriteJson(&Response{Success: 2, Msg: "移除失败,请先先添加到回收站"})
		return
	}
	if g.M("group").WhereKey(g_key).Delete() {
		g.Response.WriteJson(&Response{Success: 0, Msg: "移除成功"})
	} else {
		g.Response.WriteJson(&Response{Success: 2, Msg: "移除失败"})
	}
}
