package service

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

var (
	secKillConf *SecSkillConf
)

func InitService(serviceConf *SecSkillConf) (err error) {
	secKillConf = serviceConf
	err = loadBlackList()
	if err != nil {
		logs.Error("load black list err:%v", err)
		return
	}
	logs.Debug("init service succ, config:%v", secKillConf)

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("load proxy2layer redis pool failed, err:%v", err)
		return
	}

	secKillConf.secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	secKillConf.SecReqChan = make(chan *SecRequest, secKillConf.SecReqChanSize)
	secKillConf.UserConnMap = make(map[string]chan *SecResult, 10000)

}
func SecInfoList() (data []map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.Unlock()

	for _, v := range secKillConf.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			continue
		}
		data = append(data, item)
	}

	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.Unlock()
	item, code, err := SecInfoById(productId)
	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.Unlock()

	v, ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("没发现商品%d", productId)
	}

	start := false
	end := false
	status := "success"
	now := time.Now().Unix()
	if (now - v.StartTime) < 0 {
		start = false
		end = false
		status = "还没开始"
	}

	if (now - v.StartTime) > 0 {
		start = true
	}

	if now-v.EndTime > 0 {
		start = false
		end = true
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "卖完了"
	}

	data = make(map[string]interface{})
	data["product_id"] = v.ProductId
	data["start"] = start
	data["end"] = end
	data["status"] = status
	return
}

func userCherk(req *SecRequest) (err error) {
	found := false
	var white = [10]string{"baidu.com", "mm900.cn"}
	for _, v := range white {
		if v == req.ClientRefence {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("不是白名单")
		return
	}

	authData := fmt.Sprintf("%d%s", req.UserId, "key")
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
	if authSign != req.UserAuthSign {
		err = fmt.Errorf("非法的cookie")
		return

	}
	return
}

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.Unlock()
	err = userCherk(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		return
	}
	err = antiSpam(req)

	return
}
