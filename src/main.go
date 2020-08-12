package main

import (
	"net/http"
	"runtime"

	"controllers"
	"libs/utils"
	"models/mysql"
	"models/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// sys
	runtime.GOMAXPROCS(runtime.NumCPU())
	time.LoadLocation("Local")
	// config
	utils.InitConfig("conf", "app.conf")
	utils.InitConfig("conf", "model.conf")
	utils.InitConfig("conf", "limit.conf")
	utils.InitConfig("conf", "third.conf")
	// language
	utils.InitLanguage("static/language", "zh-cn")
	//utils.InitLanguage("static/language", "zh-hk")
	//utils.InitLanguage("static/language", "zh-tw")
	//utils.InitLanguage("static/language", "en")
	// log
	logLevel := utils.GetConfigInt("conf", "app.conf", "server", "log_level")
	utils.InitLog(logLevel, "src/log")
	// outers
	controllers.InitRoutes()
	// mysql
	mysql.InitPool()
	// redis
	redis.InitPool()
	// 服务进程状态监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		mysql.ClosePool() // 关闭mysql
		redis.ClosePool() // 关闭redis
		os.Exit(0)
	}()
	// 开启server
	host := utils.GetConfigStr("conf", "app.conf", "server", "http_domain")
	port := utils.GetConfigStr("conf", "app.conf", "server", "http_port")
	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: nil,
	}
	err := server.ListenAndServe()
	//	err := server.ListenAndServeTLS("conf/214488833170078.pem", "conf/214488833170078.key")
	utils.LogFatal("服务器启动出错 ======> ", err)
}
