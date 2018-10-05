package service

import (
	"fmt"
	"time"
)

var (
	secKillConf *SecSkillConf
)

func InitService(serviceConf *SecSkillConf) {
	secKillConf = serviceConf
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
