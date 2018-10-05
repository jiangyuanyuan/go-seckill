package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"github.com/gomodule/redigo/redis"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"workspace/go-seckill/SecProxy/service"

	"sync"
	"time"
)

var (
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
	rwLock     sync.RWMutex
)

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = "./log.log"
	config["level"] = convertLogLevel("debug")

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
func initRedis() (err error) {
	//redisPool = &redis.Pool{
	//	MaxIdle:     secKillConf.RedisConf.RedisMaxIdle,
	//	MaxActive:   secKillConf.RedisConf.RedisMaxActive,
	//	IdleTimeout: time.Duration(secKillConf.RedisConf.RedisIdleTimeout) * time.Second,
	//	Dial: func() (redis.Conn, error) {
	//		return redis.Dial("tcp", secKillConf.RedisConf.RedisAddr)
	//	},
	//}

	redisPool = &redis.Pool{
		MaxIdle:     64,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1")
		},
	}
	//conn := redisPool.Get()
	//defer conn.Close()
	//_, err = conn.Do("ping")
	//if err != nil {
	//	logs.Error("redis ping连接异常")
	//	return
	//}
	return
}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	return
}

func loadSecConf() (err error) {
	resp, err := etcdClient.Get(context.Background(), "productKey")
	if err != nil {
		logs.Error("connect loadSecConf failed, err:", err)
		return
	}

	var secProductInfo []service.SecProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}

		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	updateSecProductInfo(secProductInfo)
	return
}
func initSec() (err error) {
	err = initLogger()
	if err != nil {
		logs.Error("初始化logs失败[%s]", err)
		return
	}

	err = initRedis()
	if err != nil {
		logs.Error("初始化redis失败[%s]", err)
		return
	}
	err = initEtcd()
	if err != nil {
		logs.Error("初始化etcd失败[%s]", err)
		return
	}
	logs.Debug("初始化成功")

	err = loadSecConf()
	if err != nil {
		logs.Error("加载秒杀配置失败[%s]", err)
		return
	}
	return
}

func intConfig() (err error) {
	RedisAddr := beego.AppConfig.String("redisAddr")
	EtcdAddr := beego.AppConfig.String("etcdAddr")

	logs.Debug("加载配置成功 redis:%v", RedisAddr)
	logs.Debug("加载配置成功 etcd:%v", EtcdAddr)
	secKillConf.RedisConf.RedisAddr = RedisAddr
	secKillConf.EtcdConf.EtcdAddr = EtcdAddr
	if len(RedisAddr) == 0 || len(EtcdAddr) == 0 {
		err = fmt.Errorf("初始化失败，redis [%s]|etcd [%s]加载失败", RedisAddr, EtcdAddr)
		return
	}

	RedisMaxIdle, err := beego.AppConfig.Int("redisMaxIdle")
	if err != nil {
		err = fmt.Errorf("初始化失败，redis RedisMaxIdle 加载失败")
		return
	}
	secKillConf.RedisConf.RedisMaxIdle = RedisMaxIdle

	RedisMaxActive, err := beego.AppConfig.Int("redisMaxActive")
	if err != nil {
		err = fmt.Errorf("初始化失败，redis RedisMaxActive")
		return
	}
	secKillConf.RedisConf.RedisMaxActive = RedisMaxActive

	RedisIdleTimeout, err := beego.AppConfig.Int("redisIdleTimeout")
	if err != nil {
		err = fmt.Errorf("初始化失败，redis RedisIdleTimeout加载失败")
		return
	}
	secKillConf.RedisConf.RedisIdleTimeout = RedisIdleTimeout

	service.InitService(secKillConf)
	initSecProductWatcher()
	return
}

func initSecProductWatcher() {
	go watchSecProductKey(secKillConf.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string) {

	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	logs.Debug("begin watch key:%s", key)
	for {
		rch := cli.Watch(context.Background(), key)
		var secProductInfo []service.SecProductInfoConf
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}

	}
}

func updateSecProductInfo(secProductInfo []service.SecProductInfoConf) {

	var tmp = make(map[int]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		productInfo := v
		tmp[v.ProductId] = &productInfo
	}

	secKillConf.RWSecProductLock.Lock()
	secKillConf.SecProductInfoMap = tmp
	secKillConf.RWSecProductLock.Unlock()

}
