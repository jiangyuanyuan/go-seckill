package main

var secKillConf = &SeckillConf{}

type SeckillConf struct {
	RedisConf RedisConf
	EtcdConf  EtcdConf
}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr   string
	TimeOut    int
	EtcdSecKey string
}

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Count     int
	Left      int
}
