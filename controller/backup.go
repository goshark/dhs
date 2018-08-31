package controller

import (
	_ "database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/goshark/dhs/model"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
)

type ControllerBackup struct {
	gmvc.Controller
}

type TaskBackupContent struct {
	Serverip     string
	Dbuser       string
	Dbpass       string
	Dbport       int
	BackFileName string
	Dababases    []string
}

type ErrorRepoting struct {
	Error_info string
	Task_info  *Task
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/backup", new(ControllerBackup))
}

//实例化表
func (bak *ControllerBackup) M(tables string) *model.Models {
	return model.NewDB(tables)
}

func (bak *ControllerBackup) Exec() {
	bak.Response.WriteJson(&Response{Success: 0, Msg: "操作成功,任务已启动"})
	remoteaddr := strings.Split(bak.Request.RemoteAddr, ":")[0] //获取调用方ip
	task := new(Task)
	task_backup := new(TaskBackupContent)
	nowtime := time.Now().Unix()
	params_b, _ := GetParams(bak).ToJson()
	gjson.DecodeTo(params_b, &task)                            //解析请求数据到task结构
	task_backup_b, _ := gjson.New(task.Content.Value).ToJson() //将任务内容的value字段编码
	gjson.DecodeTo(task_backup_b, &task_backup)                //解析到备份结构

	if task.Starttime > 0 && task.Starttime > nowtime {
		time.Sleep(time.Duration((task.Starttime - nowtime)) * time.Second) //开始时间-当前时间 = 段时间,阻塞一段时间
	}
	//根据task_backup.BackFileName结合时间动态生成备份文件名称
	task_backup.BackFileName = time.Now().Format("2006_01_02_") + task_backup.BackFileName
	time_now := time.Now()
	td, err := time.ParseDuration(strconv.Itoa(task.Interval.Val) + task.Interval.Unit)
	if err != nil {
		fmt.Println("转换错误", err.Error())
	}
	var wg sync.WaitGroup
	for i := 0; i < task.Repeatcount; i++ {
		//程序体...
		wg.Add(1)
		go func() {
		INTERVALLabel:

			if time.Now().After(time_now.Add(td)) || i == 0 {
				cmd0 := exec.Command("mysqldump", "-h"+task_backup.Serverip, "-u"+task_backup.Dbuser, "-p"+task_backup.Dbpass, "-P"+strconv.Itoa(task_backup.Dbport), "-d "+strings.Join(task_backup.Dababases, " "), ">", task_backup.BackFileName)
				err := cmd0.Start()

				if err != nil {
					errrepoting_b, _ := gjson.New(&ErrorRepoting{Error_info: err.Error(), Task_info: task}).ToJson()
					ghttp.NewClient().Post("http://"+remoteaddr+":8888/task/errorepoting", ghttp.BuildParams(map[string]string{
						"jsoninfo": string(errrepoting_b),
					}))
					return
				}
				err = cmd0.Wait()
				if err != nil {
					errrepoting_b, _ := gjson.New(&ErrorRepoting{Error_info: err.Error(), Task_info: task}).ToJson()
					ghttp.NewClient().Post("http://"+remoteaddr+":8888/task/errorepoting", ghttp.BuildParams(map[string]string{
						"jsoninfo": string(errrepoting_b),
					}))
					return
				}

			}

			time_now = time.Now() //刷新现在的时间

			if i < task.Repeatcount {
				wg.Done()
				goto INTERVALLabel
			}

		}()
		wg.Wait()
	}

}

//func (bak *ControllerBackup) AppendTask() {

//}

//func (bak *ControllerBackup) DeleteTask() {

//}
