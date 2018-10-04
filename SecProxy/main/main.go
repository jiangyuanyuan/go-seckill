package main

import (
	"github.com/astaxie/beego"
	_ "workspace/go-seckill/SecProxy/router"
)

func main() {
	//err := intConfig()
	//if err != nil {
	//	panic(err)
	//	return
	//}
	err := initSec()
	if err != nil {
		panic(err)
		return
	}
	beego.Run()
}
