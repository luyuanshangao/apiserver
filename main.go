package main

import (
	"apiserver/config"
	"apiserver/model"
	"apiserver/router"
	"apiserver/router/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"

	v "apiserver/pkg/version"
)
var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)
func main() {
	pflag.Parse()

	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}
	// 初始化配置 cfg 变量值从命令⾏ flag 传⼊  可以传值⽐如 ./apiserver -c config.yaml 也可以为空 如果为空会默认 读取 conf/config.yaml
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	//初始化数据库
	model.DB.Init()
	defer model.DB.Close()

	// 设置gin模式.
	gin.SetMode(viper.GetString("runmode"))
	//创建gin引擎
	g := gin.New()

	// Routes
	router.Load(
		g,
		middleware.Logging(),
		//middleware.RequestId(),
		)
	//Ping服务器以确保路由器正常工作
	go func() {

		if err := pingServer(); err != nil {
			log.Infof("路由器没有响应，或者启动时间太长.", err)
		}
		log.Info("路由器已成功部署.")
	}()

	// Start to listening the incoming requests.
	//cert := viper.GetString("tls.cert")
	//key := viper.GetString("tls.key")
	//if cert != "" && key != "" {
	//	go func() {
	//		log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
	//		log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
	//	}()
	//}


	log.Infof("开始侦听http地址上的传入请求: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer  http服务器以确保路由器正常工作
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		//通过向“/health”发送GET请求Ping服务器
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep一会 然后继续下一个ping
		log.Info("正在等待路由器，请在1秒钟内重试.")
		time.Sleep(time.Second)
	}
	return errors.New("无法连接到路由器")
}