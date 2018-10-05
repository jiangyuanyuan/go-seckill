package service

import (
	"sync"
	"time"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

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
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        string
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	CloseNotify   <-chan bool

	//ResultChan chan *SecResult
}
