package express

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := struct {
		KeywordsType string `json:"keywords_type" form:"keywords_type"`
		Keywords     string `json:"keywords" form:"keywords"`
		trans.ReqExpressList
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	// 关键字查询
	keywordsTypes := map[string]bool{
		"company_name":   true,
		"kuaidi100_code": true,
		"taobao_code":    true,
	}
	var conds = make(mysql.Conds)
	if _, ok := keywordsTypes[req.KeywordsType]; req.KeywordsType != "" && ok {
		switch req.KeywordsType {
		case "company_name", "kuaidi100_code", "taobao_code":
			conds[req.KeywordsType] = mysql.Cond{
				"like", req.Keywords,
			}
		}
	}

	req.Sort = "is_commonly_use desc,id desc"
	total, rets := service.Express.ListPage(c, conds, req.ReqPage)
	api.JSONList(c, rets, total)
}
