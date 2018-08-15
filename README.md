# dhs
go语言结合[gitee.com/johng/gf](http://gitee.com/johng/gf)的web框架+vue前端实现mysql双机热备管理平台，支持一主多从配置。前期功能简单，持续更新。

# 架构设计

本项目(newproject)是由B/C/S架构设计所实现。轻巧便捷，将项目放到目标主机(需要进行主从配置的机器中，运行可执行文件即可开启服务)。注：目前主从配置只支持linux，windows运行可执行文件后也可直接客户端进行主从配置管理)。

# 使用教程

下载本项目(newproject)到控制端和服务端（服务端也可做控制端）,根据对应操作系统，选择对应可执行文件。链接地址https://github.com/yanyuxuanz/dhs/newproject

## 下载 ##
 ```
 go get github.com/yanyuxuanz/dhs/newproject

 ```
## 部署到服务器 ##
1.将下载后在对应工作区目录(src/github.com)找到yanyuxuanz/dhs/目录；


2.对该目录中 newproject 项目进行压缩打包；

3.上传至分别上传至主、从、和客户端（因为B/C/S架构，所以可直接用服务端来提供无需额外提供）,接下来对项目解压、设置项目权限、开启对应服务端口即可。

4.前期功能较少，源码开放欢迎大家直接fork和更新。




