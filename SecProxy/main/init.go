package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPool *redis.Pool
)

func initRedis() (err error) {
	redisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisConf.RedisAddr)
		},
	}
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("redis ping连接异常")
		return
	}
	return
}

func initEtcd() (err error) {

	return
}

func initSec() (err error) {
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
	return
}
