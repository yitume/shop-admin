package order

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := struct {
		StateType    string `json:"state_type" form:"state_type"`
		OrderType    int    `json:"order_type" form:"order_type"`
		KeywordsType string `json:"keywords_type" form:"keywords_type"`
		Keywords     string `json:"keywords" form:"keywords"`
		CreateTime   []int  `json:"create_time" form:"create_time[]"`
		trans.ReqOrderList
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	// 关键字查询
	keywordsTypes := map[string]bool{
		"goods_name":     true,
		"order_no":       true,
		"receiver_name":  true,
		"receiver_phone": true,
		"courier_number": true,
	}
	var conds = make(mysql.Conds)
	if _, ok := keywordsTypes[req.KeywordsType]; req.KeywordsType != "" && ok {
		switch req.KeywordsType {
		case "goods_name":
			conds["id"] = mysql.Cond{
				"exp", "in (SELECT GROUP_CONCAT(order_id) FROM fa_order_goods WHERE goods_title LIKE '%" + req.Keywords + "%' GROUP BY order_id)",
			}
		case "order_no":
			conds["sn"] = mysql.Cond{
				"like", req.Keywords,
			}
		case "receiver_name":
			conds["id"] = mysql.Cond{
				"exp", "in (SELECT GROUP_CONCAT(id) FROM fa_order_extend WHERE reciver_name LIKE '%" + req.Keywords + "%' GROUP BY id)",
			}
		case "receiver_phone":
			// "in (SELECT GROUP_CONCAT(id) FROM $table_order_extend WHERE receiver_phone LIKE '%".$this->keywords."%' GROUP BY id)",
			conds["id"] = mysql.Cond{
				"exp", "in (SELECT GROUP_CONCAT(id) FROM fa_order_extend WHERE receiver_phone LIKE '%" + req.Keywords + "%' GROUP BY id)",
			}
		case "courier_number":
			conds["trade_no"] = mysql.Cond{
				"like", req.Keywords,
			}
		}
	}

	// 订单时间查询
	if len(req.CreateTime) == 2 {
		conds["create_time"] = mysql.Cond{
			"between", req.CreateTime,
		}
	}

	// 订单状态查询
	if req.StateType != "" {
		state, ok := mysql.OrderStates[req.StateType]
		if !ok {
			api.JSON(c, api.MsgErr, "request app list params is error")
			return
		}
		conds["state"] = state
	}

	// 订单类型查询
	if req.OrderType != 0 {
		conds["goods_type"] = req.OrderType
	}

	total, orders := service.Order.GetOrderFullsPage(c, req.ReqPage, conds, map[string]bool{
		"order_goods":  true,
		"order_extend": true,
		"user":         true,
	})
	rets := make([]service.OrderFull, 0, len(orders))
	for _, v := range orders {
		rets = append(rets, v)
	}
	api.JSONList(c, rets, total)
}

func Info(c *gin.Context) {
	req := trans.ReqOrderList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	if req.Id == 0 {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	var err error

	order, err := service.Order.GetOrderFulls(c, []int{req.Id})
	if len(order) == 0 || err != nil {
		api.JSON(c, api.MsgErr, "OrderLog list query fail")
		return
	}
	logs, err := service.OrderLog.List(c, mysql.Conds{
		"order_id": req.Id,
	})
	if err != nil {
		api.JSON(c, api.MsgErr, "OrderLog list query fail")
		return
	}
	returns, err := service.OrderRefund.List(c, mysql.Conds{
		"order_id":    req.Id,
		"refund_type": 2,
	})
	if err != nil {
		api.JSON(c, api.MsgErr, "OrderRefund list query fail")
		return
	}
	refunds, err := service.OrderRefund.List(c, mysql.Conds{
		"order_id":    req.Id,
		"refund_type": 1,
	})
	if err != nil {
		api.JSON(c, api.MsgErr, "OrderRefund list query fail")
		return
	}

	api.JSONOK(c, gin.H{
		"info":        order[0],
		"order_log":   logs,
		"return_list": returns,
		"refund_list": refunds,
	})
}

func ChangePrice(c *gin.Context) {
	type Req struct {
		trans.ReqOrderUpdate
		ReviseGoods []struct {
			Id              int     `json:"id"`
			DifferencePrice float64 `json:"difference_price"`
			GoodsPayPrice   float64 `json:"goods_pay_price"`
		} `json:"revise_goods"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "参数错误")
		return
	}

	gids := make([]int, 0, len(req.ReviseGoods))
	for _, v := range req.ReviseGoods {
		if funk.Contains(gids, v.Id) {
			api.JSON(c, api.MsgErr, "订单参数错误")
			return
		}
		gids = append(gids, v.Id)
	}
	orderGoods, _ := service.OrderGoods.List(c, mysql.Conds{
		"id": gids,
	})
	if len(orderGoods) != len(gids) {
		api.JSON(c, api.MsgErr, "没有该商品")
		return
	}

	oid := orderGoods[0].OrderId
	cond := mysql.Conds{
		"id":           oid,
		"state":        mysql.OrderStateNew,
		"payment_time": 0,
		"goods_type":   1, // 目前只支持普通订单
	}
	order, _ := service.Order.InfoX(c, cond)
	if order.Id == 0 {
		api.JSON(c, api.MsgErr, "没有可修改的订单")
		return
	}

	for key, val := range req.ReviseGoods {
		for _, v := range orderGoods {
			if val.Id == v.Id {
				req.ReviseGoods[key].GoodsPayPrice = v.GoodsPayPrice
			}
		}
	}

	tx := model.Db.Begin()
	var sum float64
	for _, val := range req.ReviseGoods {
		price := val.GoodsPayPrice + val.DifferencePrice
		sum += price
		if price < 0 {
			api.JSON(c, api.MsgErr, "商品实际支付金额不可以小于0")
			return
		}
		if err := service.OrderGoods.Update(c, tx, val.Id, mysql.Ups{
			"goods_revise_price": gorm.Expr("goods_revise_price+100"),
		}); err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "修改失败")
			return
		}
	}

	if err := service.Order.Update(c, tx, oid, mysql.Ups{
		"revise_amount":      sum + req.ReviseFreightFee,
		"revise_freight_fee": req.ReviseFreightFee,
	}); err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "修改失败")
		return
	}

	tx.Commit()
	api.JSONOK(c)
	return
}

func SetSend(c *gin.Context) {
	type Req struct {
		mysql.OrderExtend
		IsCommonlyUse bool `json:"is_commonly_use" form:"is_commonly_use"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "参数错误")
		return
	}

	orders, err := service.Order.GetOrderFulls(c, []int{req.Id})
	if len(orders) == 0 || err != nil {
		api.JSON(c, api.MsgErr, "未找到该订单")
		return
	}
	order := orders[0]
	if order.Order.State != mysql.OrderStatePay && order.Order.State != mysql.OrderStateSend {
		api.JSON(c, api.MsgErr, "非未发货状态")
		return
	}

	if order.Order.RefundState != 0 {
		api.JSON(c, api.MsgErr, "退款状态中不可设置发货")
		return
	}

	tx := model.Db.Begin()
	ups := mysql.Ups{
		"deliver_name":    req.DeliverName,
		"deliver_phone":   req.DeliverPhone,
		"deliver_address": req.DeliverAddress,
		"need_express":    req.NeedExpress,
	}
	if req.Remark != "" {
		ups["remart"] = req.Remark
	}
	if req.ExpressId != 0 {
		ups["express_id"] = req.ExpressId
	}
	if req.TrackingNo != "" {
		ups["tracking_no"] = req.TrackingNo
	}
	now := time.Now().Unix()
	ups["tracking_time"] = now
	ups["id"] = req.Id
	if err := service.OrderExtend.Update(c, tx, req.Id, ups); err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "修改失败")
		return
	}

	ups["state"] = mysql.OrderStateSend
	ups["delay_time"] = now
	if err := service.Order.Update(c, tx, req.Id, ups); err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "订单状态修改失败")
		return
	}
	tx.Commit()

	user := auth.Default(c)
	service.OrderLog.Create(c, model.Db, &mysql.OrderLog{
		User:    user.Nickname,
		OrderId: req.Id,
		Role:    "seller",
	})
	api.JSONOK(c)
}
