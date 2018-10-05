package service

import "sync"

type SecSkillConf struct {
	RedisConf         RedisConf
	EtcdConf          EtcdConf
	LogPath           string
	Loglevel          string
	SecProductInfoMap map[int]*SecProductInfoConf
	RWSecProductLock  sync.RWMutex
}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr          string
	TimeOut           int
	EtcdSecKey        string
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}
