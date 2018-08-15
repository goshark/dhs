package controller

import (
	_ "bytes"
	"database/sql"
	"fmt"
	//"newproject/gtoken"
	"newproject/model"

	"gitee.com/johng/gf/g/encoding/gjson"

	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
	//"os"
	"os/exec"
	"strconv"
	_ "strings"

	"github.com/Unknwon/goconfig"

	_ "golang.org/x/crypto/ssh"
)

//数据库配置
const (
	DATABASE   = ""
	DBNET      = "tcp"
	DBSERVER   = "localhost"
	GRANT_USER = "user"
	GRANT_PASS = "wanjianning!@#$123"
)

type ControllerMaster struct {
	gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	ghttp.GetServer().BindController("/master", &ControllerMaster{})
}

//重写父类构造函数
func (master *ControllerMaster) Init(r *ghttp.Request) {
	defer master.Controller.Init(r)

	//身份验证层
	//获取action
	action := GetAction(r, "/master/")
	fmt.Println("当前请求:master", action)
	//	if r.Session.Get("user") == nil {
	//		r.Response.RedirectTo("/user/login-index")
	//		r.Exit()
	//	}

}

//实例化表
func (master *ControllerMaster) M(tables string) *model.Models {
	return model.NewDB(tables)
}

//配置master
/***
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,}}
***/

func (master *ControllerMaster) Exec() {

	//重启数据库
	master.restart()

	//获取masterip
	params, _ := GetParams(master).ToJson()
	conf := new(Config)
	gjson.DecodeTo(params, &conf)
	//获取数据库句柄
	DB, _ := master.getdb(conf)
	fmt.Println(conf)
	//修改为master入库
	//修改配置文件
	master.saveconfig("1")

	master.restart()
	fmt.Println("restart ok")
	DB, _ = master.getdb(conf)

	rows, _ := DB.Query("show databases")

	var databases = make([]string, 10, 10)
	var i int = 1
	for rows.Next() {
		var Database string

		rows.Scan(&Database)
		databases[i] = Database

		i++
		fmt.Println(Database)
	}
	fmt.Println("总共有" + strconv.Itoa(i) + "个数据库")

	//查询主数据库的状态
	master_status := master.showstatus(DB)
	master.Response.WriteJson(&Response{Success: 0, Msg: master_status})
}

//获取master状态主要信息
func GetMasterStatus(masterip string) map[string]interface{} {
	var resp map[string]interface{} = make(map[string]interface{})
	var resp1 Response
	conf := new(Config)
	tb_master := model.NewDB("config")
	condition, _ := gjson.Encode([]map[string]interface{}{{"name": "ip", "value": masterip, "op": "="}})
	cond, _ := gjson.DecodeToJson(condition)
	res := tb_master.Where(cond.ToArray()).Select()
	for _, v := range res {
		gjson.DecodeTo(v, &conf)
	}
	//获取数据库句柄
	masterinfo, _ := gjson.Encode(conf)
	params_b := ghttp.BuildParams(map[string]string{
		"jsoninfo": string(masterinfo),
	})
	r, e := ghttp.NewClient().Post("http://"+masterip+":8888/master/get-master-status", params_b)

	if e != nil {
		return resp
	}

	gjson.DecodeTo(r.ReadAll(), &resp1)
	if resp1.Success == 0 {
		return resp1.Msg.(map[string]interface{})
	}

	return resp
}

/***
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,}}
***/

func (master *ControllerMaster) GetMasterStatus() {
	//获取masterip
	params, _ := GetParams(master).ToJson()
	conf := new(Config)
	gjson.DecodeTo(params, &conf)
	//获取数据库句柄
	DB, _ := master.getdb(conf)
	//查询主数据库的状态

	master_status := master.showstatus(DB)
	msg := gjson.New(map[string]string{"file": master_status["File"], "pos": master_status["Position"]})
	master.Response.WriteJson(&Response{Success: 0, Msg: msg.ToMap()})
}

func (master *ControllerMaster) showstatus(DB *sql.DB) map[string]string {
	rows, err := DB.Query("show master status")
	fmt.Println(rows, err)
	//master status
	var master_status = make(map[string]string)
	for rows.Next() {
		var File string
		var Position string
		var Binlog_Do_DB string
		var Binlog_Ignore_DB string
		var Executed_Gtid_Set string
		err = rows.Scan(&File, &Position, &Binlog_Do_DB, &Binlog_Ignore_DB, &Executed_Gtid_Set)
		master_status["File"] = File
		master_status["Position"] = Position
		master_status["Binlog_Do_DB"] = Binlog_Do_DB
		master_status["Binlog_Ignore_DB"] = Binlog_Ignore_DB
		master_status["Executed_Gtid_Set"] = Executed_Gtid_Set

	}
	return master_status
}

//修改配置文件

func (master *ControllerMaster) saveconfig(server_id string) {
	var mysql_config string = "/etc/my.cnf"
	if ok, err := IsExists(mysql_config); !ok {
		fmt.Println("配置文件不存在", err)
	}
	//加载配置文件
	cfg, err := goconfig.LoadConfigFile(mysql_config)
	if err != nil {
		fmt.Println("读取配置文件失败[/etc/my.cnf]")
		return
	}
	cfg.SetValue("mysqld", "log-bin", "mysql-bin")
	cfg.SetValue("mysqld", "server-id", server_id)
	cfg.SetValue("mysqld", "binlog-ignore-db", "sys ")
	cfg.SetValue("mysqld", "binlog-ignore-db", "mysql ")
	goconfig.SaveConfigFile(cfg, mysql_config)
}

//重启数据库
func (master *ControllerMaster) restart() {
	var c = "service mysqld restart"
	cmd := exec.Command("sh", "-c", c)
	_, err := cmd.Output()
	//fmt.Println(cmd.ProcessState)
	if err != nil {
		master.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	} else {
		fmt.Println("启动mysql成功")
	}
}

//给master服务器添加授权用户
/***
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,slaveip:,}}

***/
func (master *ControllerMaster) GrantUser() {
	conf := new(Config)
	params := GetParams(master)
	params_b, _ := params.ToJson()
	gjson.DecodeTo(params_b, &conf)
	DB, err := master.getdb(conf)
	if err != nil {
		master.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	}
	//授权用户命令
	_, errs := DB.Exec("grant replication slave  on *.* to '" + GRANT_USER + "'@'" + params.GetString("slaveip") + "' identified by '" + GRANT_PASS + "' with grant option")
	if errs != nil {
		fmt.Println(errs)
		master.Response.WriteJson(&Response{Success: 2, Msg: "授权失败！," + errs.Error()})
		return
	}
	//刷新授权命令
	if _, err := DB.Exec("flush privileges"); err != nil {
		master.Response.WriteJson(&Response{Success: 2, Msg: "登录数据库失败," + err.Error()})
		return
	}
	master.Response.WriteJson(&Response{Success: 0, Msg: "授权该服务器成功！"})
}

func (master *ControllerMaster) getdb(conf *Config) (DB *sql.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", conf.Dbuser, conf.Dbpass, DBNET, DBSERVER, conf.Dbport, DATABASE)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//登陆数据库
	if err := DB.Ping(); err != nil {
		return nil, err
	} else {
		fmt.Println("数据库登录成功！")
		return DB, nil
	}
}
