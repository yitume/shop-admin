package orderrefund

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := struct {
		RefundState  int    `json:"refund_state" form:"refund_state"`
		SortType     int    `json:"sort_type" form:"sort_type"`
		KeywordsType string `json:"keywords_type" form:"keywords_type"`
		Keywords     string `json:"keywords" form:"keywords"`
		CreateTime   []int  `json:"create_time" form:"create_time[]"`
		trans.ReqOrderRefundList
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	// 关键字查询
	keywordsTypes := map[string]string{
		"goods_name":     "goods_title",
		"order_no":       "order_sn",
		"receiver_name":  "receiver_name",
		"receiver_phone": "receiver_phone",
		"refund_sn":      "refund_sn",
	}
	conds := make(mysql.Conds)
	if alias, ok := keywordsTypes[req.KeywordsType]; ok {
		switch req.KeywordsType {
		case "goods_name", "order_no", "refund_sn":
			conds[alias] = mysql.Cond{"like", req.Keywords}
		case "receiver_name":
			conds["order_id"] = mysql.Cond{"exp", "in (SELECT GROUP_CONCAT(id) FROM fa_order_extend WHERE reciver_name LIKE '%" + req.Keywords + "%' GROUP BY id)"}
		case "receiver_phone":
			conds["order_id"] = mysql.Cond{"exp", "in (SELECT GROUP_CONCAT(id) FROM fa_order_extend WHERE receiver_phone LIKE '%" + req.Keywords + "%' GROUP BY id)"}
		}
	}

	// 退款类型查询
	if req.RefundType != 0 {
		conds["refund_type"] = req.RefundType
	}

	// 退款状态查询
	if req.RefundState != 0 {
		conds["handle_state"] = req.RefundState
	}

	// 订单时间查询
	if len(req.CreateTime) == 2 {
		conds["create_time"] = mysql.Cond{
			"between", req.CreateTime,
		}
	}

	// 设置排序参数
	sortTypes := map[int]string{
		1: "id asc",
		2: "id desc",
	}
	if req.SortType != 0 {
		req.Sort = sortTypes[req.SortType]
	}

	total, list := service.OrderRefund.ListPage(c, conds, req.ReqPage)
	api.JSONList(c, list, total)
}

func Info(c *gin.Context) {
	req := trans.ReqOrderRefundList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	if req.Id == 0 {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	res, err := service.OrderRefund.Info(c, req.Id)
	if err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c, gin.H{
		"info": res,
	})
}
