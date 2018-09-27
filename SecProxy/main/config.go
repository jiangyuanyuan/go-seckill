package main

var secKillConf = &SeckillConf{}

type SeckillConf struct {
	redisAddr string
	etcdAddr  string
}
