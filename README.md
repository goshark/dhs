# dhs
go语言结合[gitee.com/johng/gf](http://gitee.com/johng/gf)的web框架+vue前端实现mysql双机热备管理平台，支持一主多从配置。前期功能简单，持续更新。

# 架构设计

本项目(newproject)是由B/C/S架构设计所实现。轻巧便捷，将项目放到目标主机(需要进行主从配置的机器中，运行可执行文件即可开启服务)。注：目前主从配置只支持linux，windows运行可执行文件后可做客户端进行主从配置管理。

# 使用教程

## 下载 ##


 ```
go get -u gitee.com/goshark/dhs

 ```
## 部署到服务器 ##
1.将下载后在对应工作区目录找到goshark/dhs/目录；


2.编译运行该项目(window/linux注意env)；

```
go build -x -v

```

3.将该项目打包（由于项目包含静态文件，建议项目整体打包）分别上传至主、从、和客户端（因为B/C/S架构，所以可直接用服务端来提供无需额外提供）；

4.接下来对项目解压、设置项目执行权限、检查服务端网络连通性,运行可执行文件即可。


5.通过任意端（部署完成后的机器）访问web（http://address:8888）即可跳转注册，登录，完成后进入管理平台首页。

6.添加服务器至列表。

7.选择对应服务器，进行主从设置。

8.完成后最后一步进入各服务端进行验证。(主服务器登录mysqld创建测试数据库和数据，看从服务器是否更新)

## 引用第三方库 ##
[gitee.com/johng/gf](http://gitee.com/johng/gf)

[github.com/Unknwon/goconfig](http://"github.com/Unknwon/goconfig")




