package goodsevaluate

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := struct {
		Type         string `json:"type" form:"type"`
		KeywordsType string `json:"keywords_type" form:"keywords_type"`
		Keywords     string `json:"keywords" form:"keywords"`
		CreateTime   []int  `json:"create_time" form:"create_time[]"`
		trans.ReqGoodsEvaluateList
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	// 关键字查询
	keywordsTypes := map[string]string{
		"goods_name":     "goods_title",
		"user_nicknname": "goods_title",
		"user_phone":     "receiver_phone",
	}
	conds := make(mysql.Conds)
	if alias, ok := keywordsTypes[req.KeywordsType]; ok {
		switch req.KeywordsType {
		case "goods_name":
			conds[alias] = mysql.Cond{"like", req.Keywords}
		case "user_nicknname":
			conds["user_id"] = mysql.Cond{"exp", "in (SELECT uid FROM fa_user_profile WHERE nickname like '%" + req.Keywords + "%'  GROUP BY id)"}
		case "user_phone":
			conds["user_id"] = mysql.Cond{"exp", "in (SELECT id FROM fa_user WHERE phone like '%" + req.Keywords + "%' GROUP BY id)"}
		}
	}

	// 订单时间查询
	if len(req.CreateTime) == 2 {
		conds["create_time"] = mysql.Cond{
			"between", req.CreateTime,
		}
	}
	// 评价查询
	if req.Type != "" {
		switch req.Type {
		case "positive":
			conds["score"] = 5
		case "moderate":
			conds["score"] = []int{3, 4}
		case "negative":
			conds["score"] = []int{1, 2}
		}
	}

	// total, list := service.GoodsEvaluate.ListPage(c, req, auth.Default(c).Id)
	total, list := service.GoodsEvaluate.ListPage(c, conds, req.ReqPage)
	api.JSONList(c, list, total)
}

func Display(c *gin.Context) {
	req := struct {
		trans.ReqGoodsEvaluateUpdate
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error,"+err.Error())
		return
	}
	if err := service.GoodsEvaluate.Update(c, model.Db, req.Id, mysql.Ups{"display": req.Display}); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error,"+err.Error())
		return
	}

	api.JSONOK(c)
	return
}
