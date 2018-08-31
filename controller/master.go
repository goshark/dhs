package controller

import (
	_ "bytes"
	"database/sql"
	"fmt"

	"gitee.com/goshark/dhs/model"

	"gitee.com/johng/gf/g/encoding/gjson"

	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"

	"os/exec"
	"strconv"
	"strings"

	"github.com/Unknwon/goconfig"

	_ "golang.org/x/crypto/ssh"
)

//数据库配置
const (
	DATABASE   = ""
	DBNET      = "tcp"
	DBSERVER   = "localhost"
	GRANT_USER = "user"
	GRANT_PASS = "goshark!@#$123"
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

/***
statement: 实例化表（gkvdb）
params:tables string
***/

func (master *ControllerMaster) M(tables string) *model.Models {
	return model.NewDB(tables)
}

/***
statement: 执行配置master
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,}}
***/
func (master *ControllerMaster) Exec() {

	//重启数据库
	master.restart()

	//获取masterip
	params := GetParams(master)
	params_b, _ := params.ToJson()
	conf := new(Config)
	gjson.DecodeTo(params_b, &conf)
	//获取数据库句柄
	DB, _ := master.getdb(conf)

	//修改配置文件,server_id为当前服务器ip末位
	lindex := strings.LastIndex(conf.Ip, ".")
	server_ip := string([]rune(conf.Ip)[lindex+1:])
	mos := params.GetString("type")
	if mos == "master" {
		master.saveconfig(map[string]string{
			"log-bin":          "mysql-bin",
			"server-id":        server_ip,
			"binlog-ignore-db": "mysql",
		})
	} else if mos == "master-master" {
		master.saveconfig(map[string]string{
			"log-bin":                  "mysql-bin",
			"binlog_format":            "mixed",
			"relay-log":                "relay-bin",
			"relay-log-index":          "slave-relay-bin.index",
			"auto_increment-increment": "2",
			"auto_increment-offset":    params.GetString("offset"),
			"server-id":                server_ip,
		})
	} else {
		master.Response.WriteJson(&Response{Success: 2, Msg: "未获取到执行对象！"})
		return
	}

	//重启数据库
	master.restart()
	//重新获取句柄
	DB, _ = master.getdb(conf)
	//查询主数据库的状态
	master_status := master.showstatus(DB)
	master.Response.WriteJson(&Response{Success: 0, Msg: master_status})
}

/***
statement: 获取master数据库列表
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,}}
***/

func (master *ControllerMaster) GetDbList() {
	params, _ := GetParams(master).ToJson()
	conf := new(Config)
	gjson.DecodeTo(params, &conf)
	//获取数据库句柄
	DB, _ := master.getdb(conf)
	rows, _ := DB.Query("show databases")
	var databaselist []string
	var database string
	for rows.Next() {
		rows.Scan(&database)
		if database != "mysql" {
			databaselist = append(databaselist, database)
		}
	}
	fmt.Println("该服务器mysql服务下共有" + strconv.Itoa(len(databaselist)) + "个数据库,分别是:" + strings.Join(databaselist, ","))
	if len(databaselist) <= 0 {
		master.Response.WriteJson(&Response{Success: 0, Msg: ""})
	} else {
		master.Response.WriteJson(&Response{Success: 0, Msg: databaselist})
	}

}

//
/***
statement: 忽略记录日志的数据库
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,dblist:[],}}
***/
func (master *ControllerMaster) IgnoreDb() {
	params := GetParams(master)
	dblist := params.GetArray("dblist")
	if len(dblist) <= 0 {
		master.Response.WriteJson(&Response{Success: 2, Msg: "暂未获取到需要忽略的数据库"})
		return
	}
	//转为字符串格式("a,b,c,d")
	var dblistslice []string
	for _, v := range dblist {
		dblistslice = append(dblistslice, v.(string))
	}

	dbliststr := strings.Join(dblistslice, ",")
	//修改数据库配置文件
	var mysql_config string = "/etc/my.cnf"
	if ok, err := IsExists(mysql_config); !ok {
		fmt.Println("配置文件不存在", err)
		master.Response.WriteJson(&Response{Success: 2, Msg: "配置文件不存在" + err.Error()})
		return
	}
	//加载配置文件
	cfg, err := goconfig.LoadConfigFile(mysql_config)
	if err != nil {
		fmt.Println("读取配置文件失败[/etc/my.cnf]")
		master.Response.WriteJson(&Response{Success: 2, Msg: "读取配置文件失败[/etc/my.cnf]" + err.Error()})
		return
	}
	ignore_db, err := cfg.GetValue("mysqld", "binlog-ignore-db")
	if err != nil {
		fmt.Println(err.Error())
		master.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	}
	//追加
	ignore_db = ignore_db + "," + dbliststr
	cfg.SetValue("mysqld", "binlog-ignore-db", ignore_db)
	if err := goconfig.SaveConfigFile(cfg, mysql_config); err != nil {
		fmt.Println("保存时出错！", err.Error())
		master.Response.WriteJson(&Response{Success: 2, Msg: "保存时出错！" + err.Error()})

	} else {
		master.Response.WriteJson(&Response{Success: 0, Msg: "保存成功!"})
	}

}

/***
statement: 获取master状态主要状态信息
params:masterip string
***/
//
func GetMasterStatus(masterip string) map[string]interface{} {
	var resp map[string]interface{} = make(map[string]interface{})
	var resp1 Response
	conf := new(Config)
	tb_master := model.NewDB("config")

	res := tb_master.Where(gjson.New([]map[string]interface{}{{"name": "ip", "value": masterip, "op": "="}}).ToArray()).Select()

	for _, v := range res {
		gjson.DecodeTo(v, &conf)
		break
	}
	//获取数据库句柄
	masterinfo, _ := gjson.New(conf).ToJson()
	fmt.Println("minfomation-------", string(masterinfo))
	r, e := ghttp.NewClient().Post("http://"+conf.Ip+":8888/master/get-master-status", ghttp.BuildParams(map[string]string{
		"jsoninfo": string(masterinfo),
	}))

	if e != nil {
		return resp
	}

	gjson.DecodeTo(r.ReadAll(), &resp1)
	fmt.Println("227tototo", resp1)
	if resp1.Success == 0 {
		if resp1.Msg == nil {
			return nil
		} else {
			return resp1.Msg.(map[string]interface{})
		}

	}

	return resp
}

/***
statement: 获取master主要状态信息api
params:{jsoninfo:{ip:,dbuser:,dbpass:,dbport:,}}
***/
func (master *ControllerMaster) GetMasterStatus() {
	//获取masterip
	params, _ := GetParams(master).ToJson()
	fmt.Println("line-242", string(params))
	conf := new(Config)
	gjson.DecodeTo(params, &conf)
	//获取数据库句柄
	DB, _ := master.getdb(conf)
	fmt.Println(DB)
	//查询主数据库的状态
	fmt.Println("line-248", conf)
	master_status := master.showstatus(DB)
	msg := gjson.New(map[string]string{"file": master_status["File"], "pos": master_status["Position"]})
	fmt.Println("line-251被调用方", conf, msg.ToMap())
	master.Response.WriteJson(&Response{Success: 0, Msg: msg.ToMap()})
}

/***
statement: show master status for calling.
params:DB *sql.DB
***/

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

/***
statement:修改/etc/my.cnf配置文件
params:server_id 服务编号 （当前服务Ip的末位）
***/

func (master *ControllerMaster) saveconfig(params map[string]string) {
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
	for k, v := range params {
		cfg.SetValue("mysqld", k, v)
	}
	goconfig.SaveConfigFile(cfg, mysql_config)
}

/***
statement:重启数据库方法
params:nil
***/

func (master *ControllerMaster) restart() {
	var c = "service mysqld restart"
	cmd := exec.Command("sh", "-c", c)
	_, err := cmd.Output()
	if err != nil {
		master.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	} else {
		fmt.Println("启动mysql成功")
	}
}

/***
statement:给master服务器添加授权用户
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

/***
statement:获取数据库句柄方法
params:conf *Config{conf.Dbuser, conf.Dbpass, DBNET, DBSERVER, conf.Dbport}
***/

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
