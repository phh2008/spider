package main

import (
	_ "com.phh/spider/rules"
	"github.com/henrylee2cn/pholcus/config"
	"github.com/henrylee2cn/pholcus/exec"
	"github.com/henrylee2cn/pholcus/runtime/cache"
	//_ "github.com/pholcus/spider_lib" // 此为公开维护的spider规则库
)

func init() {
	// 标记当前init()已执行完毕
	defer cache.ExecInit(0)
	//数据库
	config.DB_NAME = "pholcus"
	//mongodb链接字符串
	config.MGO_CONN_STR = "127.0.0.1:27017"
	//mongodb连接池容量
	config.MGO_CONN_CAP = 1024
	//mysql服务器地址
	config.MYSQL_CONN_STR = "root:root@tcp(192.168.1.216:3306)"
	//mysql连接池容量
	config.MYSQL_CONN_CAP = 1024

}

func main() {
	// 设置运行时默认操作界面，并开始运行
	// 运行软件前，可设置 -a_ui 参数为"web"、"gui"或"cmd"，指定本次运行的操作界面
	// 其中"gui"仅支持Windows系统
	exec.DefaultRun("web")
}
