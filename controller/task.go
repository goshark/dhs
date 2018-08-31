package controller

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gitee.com/goshark/dhs/model"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
)

type ControllerTask struct {
	gmvc.Controller
}

//任务周期
type Interval struct {
	Unit string `json:"unit,omitempty"` //单位。y,m,d,h,i
	Val  int    `json:"val,omitempty"`  //值
}

//任务内容
type Content struct {
	Serverip string                 `json:"server,omitempty"`   //目标服务器ip
	Name     string                 `json:"name,omitempty"`     //内容名称
	Category string                 `json:"category,omitempty"` //内容类别 "backup"备份数据库
	Value    map[string]interface{} `json:"value,omitempty"`    //任务内容的参数值

}

//任务
type Task struct {
	Uuid        string                      `json:"uuid,omitempty"`      //任务uuid
	Intime      int64                       `json:"intime,omitempty"`    //任务录入时间
	Starttime   int64                       `json:"starttime,omitempty"` //任务开始时间
	Interval    `json:"interval,omitempty"` //任务周期
	Deleted     bool                        `json:"deleted,omitempty"`     //是否删除(软删除)
	Status      int                         `json:"status,omitempty"`      //任务执行状态,0:未执行,1:执行中,2:执行完成
	Name        string                      `json:"name,omitempty"`        //任务名称
	Repeatcount int                         `json:"repeatcount,omitempty"` //执行次数.0:重复执行,n>0,执行n次
	Content     `json:"content,omitempty"`  //任务内容
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/task", new(ControllerTask))
}

//重写父类构造函数
func (t *ControllerTask) Init(r *ghttp.Request) {

	defer t.Controller.Init(r)

	//身份验证层
	//获取action
	action := GetAction(r, "/task/")
	fmt.Println("当前请求:task", action)
	//	if r.Session.Get("user") == nil {
	//		r.Response.RedirectTo("login-index")
	//		r.Exit()
	//	}

}

//实例化表
func (t *ControllerTask) M(tables string) *model.Models {
	return model.NewDB(tables)
}

//查询任务列表
func (t *ControllerTask) TaskList() {
	params := GetParams(t)
	var result []interface{}
	res := t.M("task").Where(params.GetArray("condition")).Select()
	if len(res) <= 0 {
		t.Response.WriteJson(&Response{Success: 0, Msg: ""})
		return
	}
	for _, v := range res {
		b, _ := gjson.DecodeToJson(v)
		result = append(result, b.ToMap())
	}
	t.Response.WriteJson(&Response{Success: 0, Msg: result})

}

//添加新任务
func (t *ControllerTask) AppendTask() {
	task := new(Task)
	params_b, _ := GetParams(t).ToJson()
	gjson.DecodeTo(params_b, &task)
	task.Uuid = NewUuid("task")
	task.Intime = time.Now().Unix()
	task_b, _ := gjson.New(task).ToJson()
	task_key := "task_" + strconv.Itoa(len(t.M("task").GetKey())+1)
	if t.M("task").Data(task_b).WhereKey(task_key).Insert() {
		t.Response.WriteJson(&Response{Success: 0, Msg: "添加成功!"})
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "添加失败!"})
	}

}

//更改任务
func (t *ControllerTask) UpdateTask() {
	params := GetParams(t)
	params_b, _ := params.ToJson()
	var task_key string
	var task_old Task
	var task Task
	gjson.DecodeTo(params_b, &task)
	task_oval := t.M("task").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": params.GetString("uuid"), "op": "="}}).ToArray()).Select()
	for k, v := range task_oval {
		task_key = k
		gjson.DecodeTo(v, &task_old)
		break
	}

	task_o_elem := reflect.ValueOf(&task_old).Elem()
	Merge(task_old, task, task_o_elem)
	task_b, err := gjson.New(task_old).ToJson()

	if err != nil {
		fmt.Println(err.Error())
	}

	if t.M("task").Data(task_b).WhereKey(task_key).Update() {
		t.Response.WriteJson(&Response{Success: 0, Msg: "修改成功"})
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "修改失败"})
	}
}

//删除任务(软删除到回收站)
func (t *ControllerTask) DeleteTask() {
	var task_key string
	uuid := GetParams(t).GetString("uuid")
	task := new(Task)
	taskm := t.M("task").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": uuid, "op": "="}}).ToArray()).Select()
	for k, v := range taskm {
		task_key = k
		gjson.DecodeTo(v, &task)
		break
	}
	task.Deleted = true
	task_b, _ := gjson.New(task).ToJson()
	if t.M("task").Data(task_b).WhereKey(task_key).Update() {
		t.Response.WriteJson(&Response{Success: 0, Msg: "删除成功！"})
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "删除失败！"})
	}
}

//恢复任务到列表
func (t *ControllerTask) RecoverTask() {
	var task_key string
	uuid := GetParams(t).GetString("uuid")
	task := new(Task)
	taskm := t.M("task").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": uuid, "op": "="}}).ToArray()).Select()
	for k, v := range taskm {
		task_key = k
		gjson.DecodeTo(v, &task)
		break
	}
	task.Deleted = false
	task_b, _ := gjson.New(task).ToJson()
	if t.M("task").Data(task_b).WhereKey(task_key).Update() {
		t.Response.WriteJson(&Response{Success: 0, Msg: "恢复成功！"})
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "恢复失败！"})
	}
}

//彻底删除任务
func (t *ControllerTask) RemoveTask() {
	var task_key string
	uuid := GetParams(t).GetString("uuid")
	task := new(Task)
	taskm := t.M("task").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": uuid, "op": "="}}).ToArray()).Select()
	for k, v := range taskm {
		task_key = k
		gjson.DecodeTo(v, &task)
		break
	}

	if !task.Deleted {
		t.Response.WriteJson(&Response{Success: 2, Msg: "移除失败,请先先添加到回收站"})
		return
	}

	if t.M("task").WhereKey(task_key).Delete() {
		t.Response.WriteJson(&Response{Success: 0, Msg: "移除成功!"})
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "移除失败!"})
	}
}

//执行任务
func (t *ControllerTask) Exec() {
	var task_key string
	task_uuid := GetParams(t).GetString("uuid")
	task := new(Task)
	res := t.M("task").Where(gjson.New([]map[string]interface{}{{"name": "uuid", "value": task_uuid, "op": "="}}).ToArray()).Select()
	for k, v := range res {
		task_key = k
		gjson.DecodeTo(v, &task)
		break
	}
	if task.Status == 1 {
		t.Response.WriteJson(&Response{Success: 2, Msg: "任务正在进行中不能重复执行"})
		return
	} else if task.Status == 2 {
		t.Response.WriteJson(&Response{Success: 2, Msg: "任务已执行完成了"})
		return
	}
	task.Status = 1
	task_b, _ := gjson.New(task).ToJson()
	var resp Response
	if t.M("task").Data(task_b).WhereKey(task_key).Update() {
		//解析task任务内容，模拟客户端发送到目标服务器执行任务接口
		r, e := ghttp.NewClient().Post("http://"+task.Content.Serverip+":8888/"+task.Content.Category+"/exec", ghttp.BuildParams(map[string]string{
			"jsoninfo": string(task_b),
		}))
		if e != nil {
			fmt.Println("line229", e.Error())
			t.Response.WriteJson(&Response{Success: 2, Msg: "任务执行出错" + e.Error()})
			return
		}

		gjson.DecodeTo(r.ReadAll(), &resp)
		if resp.Success == 0 {
			t.Response.WriteJson(&Response{Success: 0, Msg: "操作成功,任务正在执行"})
			return
		} else {
			t.Response.WriteJson(&resp)
		}
	} else {
		t.Response.WriteJson(&Response{Success: 2, Msg: "任务执行失败"})
		return
	}

}
