package model

import (
	"fmt"

	"gitee.com/johng/gf/g/encoding/gjson"
	"gitee.com/johng/gf/g/os/glog"
	"gitee.com/johng/gkvdb/gkvdb"
	//"errors"
	"strings"
)

const (
	DBFILE = "MyDB"
)

type Models struct {
	db         *gkvdb.DB   // 数据库操作对象
	tables     string      // 数据库操作表
	key        string      //db.key
	keys       []string    // db.keys
	fields     []string    //保留字段
	data       []byte      //数据
	conditions []Condition //查询条件
}

type Condition struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Op    string      `json:"op,omitempty"` //"=","!=","~="
}

func NewDB(tables string) *Models {
	db, err := gkvdb.New(DBFILE)
	if err != nil {
		glog.Error(err)
	}
	return &Models{
		db:     db,
		tables: tables,
	}
}

func (md *Models) WhereKeys(args ...string) *Models {
	md.keys = append(md.keys, args...)
	return md
}

func (md *Models) WhereKey(arg string) *Models {
	md.key = arg
	return md
}

func (md *Models) Fields(args ...string) *Models {
	md.fields = append(md.fields, args...)
	return md
}

func (md *Models) Delete() bool {
	if len(md.keys) > 0 {
		for k, v := range md.keys {
			if err := md.db.RemoveFrom([]byte(v), md.tables); err != nil {
				glog.Errorfln("The operations of delete has an error when running.Error : The ", k+1, ".Th operate is failed, details:"+err.Error())
			}
		}
	}

	if len(md.key) > 0 {
		if err := md.db.RemoveFrom([]byte(md.key), md.tables); err != nil {
			glog.Errorfln("The operations of delete has an error when running.Error : operate is failed, details:" + err.Error())
			return false
		}

	}
	return true
}

func (md *Models) GetKey() []string {
	if md.tables == "" {
		glog.Printfln("method: getkey()", "Table can not be None.")
		return []string{`{Success:0,Msg:"Field Table can not be None."}`}

	}
	table, _ := md.db.Table(md.tables)
	return table.Keys(-1)
}

func (md *Models) GetVal() [][]byte {
	if md.tables == "" {
		glog.Printfln("method: getval()", "Table can not be None.")
		return [][]byte{[]byte(`{Success:0,Msg:"Field Table can not be None."}`)}

	}
	table, _ := md.db.Table(md.tables)
	return table.Values(-1)
}

func (md *Models) Update() bool { //key如果相同则覆盖更新，如果不同则插入
	if md.tables == "" {
		glog.Printfln("method: update()", "Table can not be None.")
		return false

	}
	table, _ := md.db.Table(md.tables)
	if err := table.Set([]byte(md.key), md.data); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func (md *Models) Insert() bool { //如果不同则插入
	if md.tables == "" {
		glog.Printfln("method: insert()", "Table can not be None.")
		return false

	}
	table, _ := md.db.Table(md.tables)
	if err := table.Set([]byte(md.key), md.data); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//返回一条
func (md *Models) Find() []byte {
	if md.tables == "" {
		glog.Printfln("method: find()", "Table can not be None.")
		return []byte(`{Success:0,Msg:"Field Table can not be None."}`)
	}
	fmt.Println(string(md.db.GetFrom([]byte(md.key), md.tables)))
	return md.db.GetFrom([]byte(md.key), md.tables)

}

//kvdb数据
func (md *Models) Data(data []byte) *Models {
	md.data = data
	return md
}

func (md *Models) SelectAll() [][]byte {

	table, _ := md.db.Table(md.tables)
	return table.Values(-1)
}

func (md *Models) Where(data []interface{}) *Models {
	var condition Condition //name,value,op
	var conditions []Condition
	for _, v := range data {
		condition_b, _ := gjson.Encode(v)
		gjson.DecodeTo(condition_b, &condition)
		conditions = append(conditions, condition)
	}
	md.conditions = conditions
	return md
}

//查询
func (md *Models) Select() map[string][]byte {

	var res map[string][]byte = make(map[string][]byte)
	if md.tables == "" {
		glog.Printfln("method: select()", "Table can not be None.")
		return map[string][]byte{
			"Success": []byte("2"),
			"Msg":     []byte("Field Table can not be None."),
		}

	}
	table, _ := md.db.Table(md.tables)
	result := table.Items(-1)
	if len(md.conditions) > 0 {
		for _, v := range md.conditions {
			for k1, v1 := range result {

				j, err := gjson.DecodeToJson(v1)

				if err != nil {
					fmt.Println(err)
				}

				if v.Op == "=" {

					if j.Get(v.Name) == v.Value {

						res[k1] = v1

					}
				}
				if v.Op == "!=" {
					if j.Get(v.Name) != v.Value {

						res[k1] = v1

					}
				}
				if v.Op == "~=" {
					if strings.Contains(j.Get(v.Name).(string), v.Value.(string)) {

						res[k1] = v1

					}
				}
			}
		}
	} else {
		res = result
	}

	if len(md.fields) > 0 {
		for k, v := range res {
			res[k] = fields(md.fields, v)
		}
	}
	return res

}

func fields(args []string, data []byte) []byte {
	b, _ := gjson.DecodeToJson(data)
	var res map[string]interface{} = make(map[string]interface{})
	for _, v := range args {
		for k1, v1 := range b.ToMap() {
			if v == k1 {
				res[v] = v1
			}
		}
	}
	bres, _ := gjson.Encode(res)
	return bres
}
