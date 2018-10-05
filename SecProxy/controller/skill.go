package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"workspace/go-seckill/SecProxy/service"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	p.Data["json"] = "SecKill"
	p.ServeJSON(true)
}

func (p *SkillController) SecInfo() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 1000
	result["message"] = "成功"
	defer func() {
		p.Data["json"] = result
		p.ServeJSON(true)
	}()
	if err != nil {
		//result["code"] = service.ErrInvalidRequest
		//result["message"] = "非法ID"
		//p.Data["json"] = result
		//p.ServeJSON(true)
		//return
		data, code, err := service.SecInfoList()
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()

			logs.Error("invalid request, get product_id failed, err:%v", err)
			return
		}

		result["code"] = code
		result["data"] = data
	} else {
		data, code, err := service.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			p.Data["json"] = result
			p.ServeJSON(true)
			return
		}
		result["data"] = data
	}

}
