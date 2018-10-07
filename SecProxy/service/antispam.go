package service

import (
	"fmt"
	"sync"
)

var secLimitMgr = &SecLimitMgr{
	UserLimitMap: make(map[int]*SecLimit, 10000),
	IpLimitMap:   make(map[string]*SecLimit, 10000),
}

type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	IpLimitMap   map[string]*SecLimit
	lock         sync.RWMutex
}
type SecLimit struct {
	Count   int
	CurTime int64
}

func (p *SecLimit) count(nowTime int64) (curCount int) {
	if p.CurTime != nowTime {
		p.CurTime = nowTime
		p.Count = 1
		return
	}
	p.Count++
	curCount = p.Count
	return
}

func (p *SecLimit) check(nowTime int64) int {
	if p.CurTime != nowTime {
		return 0
	}
	return p.Count
}
func antiSpam(req *SecRequest) (err error) {
	secLimitMgr.lock.Lock()
	defer secLimitMgr.lock.Unlock()

	secLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		secLimit = &SecLimit{}
		secLimitMgr.UserLimitMap[req.UserId] = secLimit
	}
	count := secLimit.count(req.AccessTime.Unix())

	if count > 5 {
		err = fmt.Errorf("请求频繁")
		return
	}

	secIpLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		secIpLimit = &SecLimit{}
		secLimitMgr.IpLimitMap[req.ClientAddr] = secIpLimit
	}
	ipCount := secLimit.count(req.AccessTime.Unix())
	if ipCount > 50 {
		err = fmt.Errorf("同IP请求频繁")
		return
	}
	return
}
