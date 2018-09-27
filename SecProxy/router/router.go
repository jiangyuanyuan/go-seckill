package router

import (
	"github.com/astaxie/beego"
	"workspace/go-seckill/SecProxy/controller"
)

func init() {
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}
