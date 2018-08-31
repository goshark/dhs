package controller

import (
	_ "bytes"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitee.com/johng/gf/g/encoding/gjson"

	"gitee.com/goshark/dhs/model"

	"os/exec"
	"strings"

	"gitee.com/johng/gf/g/frame/gmvc"
	"gitee.com/johng/gf/g/net/ghttp"
	"github.com/Unknwon/goconfig"
	_ "golang.org/x/crypto/ssh"
)

type ControllerSlave struct {
	gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server 从数据库绑定处理
func init() {
	ghttp.GetServer().BindController("/slave", &ControllerSlave{})
}

//重写父类构造函数
//func (slave *ControllerSlave) Init(r *ghttp.Request) {
//	defer slave.Controller.Init(r)

//	//身份验证层
//	//获取action
//	action := GetAction(r, "/slave/")
//fmt.Println("当前请求:slave", action)
//	if action != "login-index" { //登录和注册不存在token验证
//		if r.Session.Get("user") == nil {
//			r.Response.RedirectTo("/user/login-index")
//			r.Exit()
//		}

//	}
//}

//实例化表
func (slave *ControllerSlave) M(tables string) *model.Models {
	return model.NewDB(tables)
}

//修改配置文件
func (slave *ControllerSlave) saveconfig(server_id string) {
	var mysql_config string = "/etc/my.cnf"
	if ok, err := IsExists(mysql_config); !ok {
		fmt.Println("配置文件不存在", err)
	}
	cfg, err := goconfig.LoadConfigFile(mysql_config)
	if err != nil {
		log.Fatal(err)
	}
	cfg.SetValue("mysqld", "log-bin", "mysql-bin")
	cfg.SetValue("mysqld", "server-id", server_id)
	goconfig.SaveConfigFile(cfg, mysql_config)
}

//重启数据库服务
func (slave *ControllerSlave) restart() {
	var c = "service mysqld restart"
	var cmd = exec.Command("sh", "-c", c)
	_, err := cmd.Output()
	//fmt.Println(cmd.ProcessState)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("重启mysql成功")
	}
}

//获取主服务器配置ip
func GetMasterIp() (masterip string) {
	tb_config := model.NewDB("config")
	res := tb_config.Where(gjson.New([]map[string]interface{}{{"name": "type", "value": "master", "op": "="}}).ToArray()).Select()
	if len(res) <= 0 {
		masterip = ""
		return
	}
	for _, v := range res {
		ip, _ := gjson.DecodeToJson(v)
		masterip = ip.GetString("ip")
		break
	}
	return
}
func (slave *ControllerSlave) Exec() {
	fmt.Println("设置slave配置")
	params := GetParams(slave) //从信息+masterip+maser_status
	params_b, _ := params.ToJson()
	conf := new(Config)
	gjson.DecodeTo(params_b, &conf)

	//修改配置文件,server_id为当前服务器ip末位
	lindex := strings.LastIndex(conf.Ip, ".")
	server_ip := string([]rune(conf.Ip)[lindex+1:])
	slave.saveconfig(server_ip)

	//重启从数据库服务
	slave.restart()

	//登陆数据库
	DB, _ := slave.getdb(conf)
	DB.Exec("stop slave")
	masterip := params.GetString("masterip")
	master_status_file := params.GetString("file")
	master_status_pos := params.GetString("position")
	changesql := "change master to master_host='" + masterip + "',master_port=3306,master_user='" + GRANT_USER + "',master_password='" + GRANT_PASS + "',master_log_file='" + master_status_file + "',master_log_pos=" + master_status_pos
	fmt.Println("change master sql", changesql)
	_, errs := DB.Exec(changesql)

	if errs != nil {
		slave.Response.WriteJson(&Response{Success: 2, Msg: errs.Error()})
		return
	}
	_, errs = DB.Exec("start slave")
	if errs != nil {
		slave.Response.WriteJson(&Response{Success: 2, Msg: errs.Error()})
		return
	}
	var res bool
	for i := 0; i <= 4; i++ {
		if slave.checkslavestatus(DB) {
			res = true
			break
		}
		time.Sleep(time.Duration(time.Second * 1))
	}
	if res {
		slave.Response.WriteJson(&Response{Success: 0, Msg: "配置从服务器成功！"})
	} else {
		slave.Response.WriteJson(&Response{Success: 2, Msg: "配置从服务器失败！"})
	}

}

func (slave *ControllerSlave) checkslavestatus(DB *sql.DB) bool {
	//查询从的状态
	var rows, _ = DB.Query("show slave status")
	//字段

	cols, _ := rows.Columns()

	fmt.Println("=================================")
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println("Error")
			return false
		}
		row := make(map[string]string)
		for j, v := range values {
			key := cols[j]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}

	fmt.Println("Slave_IO_Running---->", results[0]["Slave_IO_Running"])
	fmt.Println("Slave_SQL_Running---->", results[0]["Slave_SQL_Running"])
	if results[0]["Slave_IO_Running"] == "Yes" && results[0]["Slave_SQL_Running"] == "Yes" {
		return true
	} else {
		return false
	}
}
func (slave *ControllerSlave) ShowSlaveStatus() {
	conf := new(Config)
	params, _ := GetParams(slave).ToJson()
	gjson.DecodeTo(params, &conf)
	DB, err := slave.getdb(conf)
	if err != nil {
		fmt.Println(err.Error())
		slave.Response.WriteJson(&Response{Success: 2, Msg: err.Error()})
		return
	}

	//查询从的状态
	var rows, _ = DB.Query("show slave status")
	//字段

	cols, _ := rows.Columns()

	fmt.Println("=================================")
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println("Error")
			return
		}
		row := make(map[string]string)
		for j, v := range values {
			key := cols[j]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}

	slave.Response.WriteJson(&Response{Success: 0, Msg: results[0]})

}

func (slave *ControllerSlave) getdb(conf *Config) (DB *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", conf.Dbuser, conf.Dbpass, DBNET, DBSERVER, conf.Dbport, DATABASE)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//登陆数据库
	err = DB.Ping()
	fmt.Println(err, "ping")
	if err != nil {
		return nil, err
	}
	return DB, nil

}
