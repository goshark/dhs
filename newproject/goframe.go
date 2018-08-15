package main

import (
	_ "newproject/controller"

	"gitee.com/johng/gf/g"
	"gitee.com/johng/gf/g/net/ghttp"
	"gitee.com/johng/gf/g/os/glog"
)

const SERVER_PORT = 8888

func main() {
	g.View().SetPath("view")
	g.View().SetDelimiters("${", "}")
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {

		if r.Session.Get("user") == nil {
			r.Response.RedirectTo("/user/login-index")
		} else {
			content, _ := g.View().Parse("index.html", map[string]interface{}{
				"user": r.Session.Get("user"),
			})
			r.Response.Write(content)
		}

	})
	s.SetPort(SERVER_PORT)
	s.Run()
	glog.SetPath("log")
	glog.SetDebug(true)
}
