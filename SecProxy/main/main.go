package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "workspace/go-seckill/SecProxy/router"
)

func intConfig() (err error) {
	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	logs.Debug("加载配置成功 redis:%v", redisAddr)
	logs.Debug("加载配置成功 etcd:%v", etcdAddr)
	secKillConf.redisAddr = redisAddr
	secKillConf.etcdAddr = etcdAddr

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("初始化失败，redis [%s]|etcd [%s]加载失败", redisAddr, etcdAddr)
		return
	}
	return
}

func main() {
	err := intConfig()
	if err != nil {
		panic(err)
		return
	}
	err = initSec()
	if err != nil {
		panic(err)
		return
	}
	beego.Run()
}
