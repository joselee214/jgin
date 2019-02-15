package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"jose/josegin/jgin"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	sessions "github.com/tommy351/gin-sessions"

	//"fmt"
	//"time"
	//"net/http"

	// etcd v3 registry
	_ "github.com/micro/go-plugins/registry/etcdv3"
	//"github.com/micro/go-web"
	//"github.com/micro/go-micro"
	//"log"

	"jose/josegin/controller"
)

func registerRouter(ginServer *gin.Engine) {
	new(controller.ApisController).Router(ginServer)
	//new(controller.PageController).Router(ginServer)
	new(controller.DocsController).Router(ginServer)
	//new(controller.UserController).Router(ginServer)
	//new(controller.TestController).Router(router)
	//new(controller.UserController).Router(router)
}

func main() {

	//配置读取
	cfg := new(jgin.Config)
	cfg.Parse("config/app.properties")
	jgin.SetCfg(cfg)

	//设定日志路径//
	jgin.Configuration(cfg.Logger["filepath"])

	//设置运行模式
	gin.SetMode(cfg.App["mode"])

	//连接 Redis
	for k, ds := range cfg.Redis {
		db, _ := strconv.Atoi(ds["db"])
		maxIdle, _ := strconv.Atoi(ds["maxIdle"])
		maxActive, _ := strconv.Atoi(ds["maxActive"])
		idleTimeoutSecond, _ := strconv.Atoi(ds["idleTimeoutSecond"])
		jgin.ConnPools.AddRedis(k, ds["host"], ds["passwd"], db, maxIdle, maxActive, idleTimeoutSecond)
	}

	//获取Redis使用
	//conn := jgin.ConnPools.GetRedisA("default")
	//defer conn.Close()
	//
	//conn.Do("SELECT", 0)
	//conn.Do("FLUSHALL")

	////	conn.Do("SET","jose","22222")
	//	conn.Do("DEL","jose")
	//	z,_ := redis.String(conn.Do("GET","jose"))
	//	fmt.Println("==================",z)

	//设置 Xrom 缓存
	// cacher := jgin.XromNetCacher(jgin.XromRedisCacheStore(conn,-1))
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 10000)

	for k, ds := range cfg.Datasource {
		e, _ := xorm.NewEngine(ds["driveName"], ds["dataSourceName"])
		// e.SetDefaultCacher(cacher) //设置缓存
		// e.SetLogLevel(9) //Xrom信息debug
		e.ShowSQL(ds["showSql"] == "true") //显示sql查询

		n, _ := strconv.Atoi(ds["maxIdle"])
		e.SetMaxIdleConns(n)
		n, _ = strconv.Atoi(ds["maxLifetime"])
		var duration_Seconds time.Duration = time.Second * time.Duration(n)
		e.SetConnMaxLifetime(duration_Seconds)
		n, _ = strconv.Atoi(ds["maxOpen"])
		e.SetMaxOpenConns(n)
		jgin.SetEngin(k, e)

		//e.Sync2(new (entity.User))  //初始化数据？
	}

	ginServer := gin.Default()

	//session 相关支持
	//ginServer.Use(libs.Auth())
	sessionmaxLifetime, _ := strconv.Atoi(cfg.Session["timelive"])
	store := jgin.NewCookieStore(sessionmaxLifetime, []byte("1"))
	ginServer.Use(sessions.Middleware(cfg.Session["name"], store))

	//静态路由映射
	//for k,v :=range cfg.Static{
	//	ginServer.Static(k, v)
	//}∫
	for k, v := range cfg.StaticFile {
		ginServer.StaticFile(k, v)
	}

	//路由相关//静态文件
	ginServer.SetFuncMap(jgin.GetFuncMap())
	ginServer.NoRoute(jgin.NoRoute)
	ginServer.NoMethod(jgin.NoMethod)
	registerRouter(ginServer)

	//view相关
	//ginServer.LoadHTMLGlob(cfg.View["path"]+"/**/*")
	//ginServer.Delims(cfg.View["deliml"],cfg.View["delimr"])

	//fmt.Println(ginServer)

	//fmt.Println(cfg.App)

	//测试
	ginServer.GET("/jgin/ping", func(c *gin.Context) {
		//c.HTML(http.StatusOK,"user/login.html","")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//端口监听
	s := &http.Server{
		Addr:           cfg.App["addr"] + ":" + cfg.App["port"],
		Handler:        ginServer,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[ok] app run", cfg.App["addr"]+":"+cfg.App["port"])
	}
	ginServer.Run(cfg.App["addr"] + ":" + cfg.App["port"])

	//整合进 integrate micro web/api
	//service := web.NewService(
	//	//web.Name("com.7yes.web.jgin"),  // http://localhost:8082/jgin/user/findOne?userId=1 //micro会截断amespace=jgin
	//	web.Name("com.7yes.api.jgin"), // http://localhost:8080/jgin/user/findOne?userId=1 //需要完整处理路径
	//)
	//service.Init()
	//service.Handle("/",ginServer)
	//if err := service.Run(); err != nil {
	//	log.Fatal(err)
	//}

}
