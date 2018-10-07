package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

func loadBlackList() (err error) {

	secKillConf.ipBlackMap = make(map[string]bool, 10000)
	secKillConf.idBlackMap = make(map[int]bool, 10000)

	err = initBlackRedis()
	if err != nil {
		logs.Error("init black redis failed, err:%v", err)
		return
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "idblacklist")
	idlist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}

	for _, v := range idlist {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id [%v]", id)
			continue
		}
		secKillConf.idBlackMap[id] = true
	}

	reply, err = conn.Do("hgetall", "ipblacklist")
	iplist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}

	for _, v := range iplist {
		secKillConf.ipBlackMap[v] = true
	}

	go SyncIpBlackList()
	go SyncIdBlackList()
	return
}

func initBlackRedis() (err error) {
	secKillConf.blackRedisPool = &redis.Pool{
		MaxIdle:     64,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1")
		},
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		conn := secKillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			secKillConf.RWBlackLock.Lock()
			for _, v := range ipList {
				secKillConf.ipBlackMap[v] = true
			}
			secKillConf.RWBlackLock.Unlock()

			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip[%v]", ipList)
		}
	}
}

func SyncIdBlackList() {
	for {
		conn := secKillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}

		secKillConf.RWBlackLock.Lock()
		secKillConf.idBlackMap[id] = true
		secKillConf.RWBlackLock.Unlock()

		logs.Info("sync id list from redis succ, ip[%v]", id)
	}
}
