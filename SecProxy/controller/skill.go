package controller

import "github.com/astaxie/beego"

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	p.Data["json"] = "SecKill"
	p.ServeJSON(true)
}

func (p *SkillController) SecInfo() {
	p.Data["json"] = "SecInfo"
	p.ServeJSON(true)
}
