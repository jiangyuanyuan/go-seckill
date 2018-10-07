package service

import (
	"github.com/gomodule/redigo/redis"
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
	CookieSecretKey   string

	ReferWhiteList []string

	ipBlackMap map[string]bool
	idBlackMap map[int]bool

	AccessLimitConf      AccessLimitConf
	blackRedisPool       *redis.Pool
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool

	secLimitMgr *SecLimitMgr

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	SecReqChan     chan *SecRequest
	SecReqChanSize int

	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
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
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	CloseNotify   <-chan bool

	ResultChan chan *SecResult
}
type AccessLimitConf struct {
	IPSecAccessLimit   int
	UserSecAccessLimit int
	IPMinAccessLimit   int
	UserMinAccessLimit int
}
type SecResult struct {
	ProductId int
	UserId    int
	Code      int
	Token     string
}
