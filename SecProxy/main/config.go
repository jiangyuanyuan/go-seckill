package main

import (
	"workspace/go-seckill/SecProxy/service"
)

var secKillConf = &service.SecSkillConf{
	SecProductInfoMap: make(map[int]*service.SecProductInfoConf, 1024),
}
