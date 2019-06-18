package user

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := trans.ReqUserList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	total, list := service.UserOpen.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)
}

func Info(c *gin.Context) {
	req := trans.ReqUserList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	user, _ := service.UserOpen.Info(c, req.Uid)
	p, _ := service.UserProfile.InfoX(c, mysql.Conds{
		"uid": user.Uid,
	})
	type res struct {
		mysql.UserOpen
		Name       string `json:"name" form:"name" `               // 客户姓名
		Nickname   string `json:"nickname" form:"nickname" `       // 昵称
		Avatar     string `json:"avatar" form:"avatar" `           // 头像
		Sex        int    `json:"sex" form:"sex" `                 // 1男0女
		Birthday   int    `json:"birthday" form:"birthday" `       // 生日
		Qq         string `json:"qq" form:"qq" `                   // QQ
		DeleteTime int64  `json:"delete_time" form:"delete_time" ` // 软删除时间
	}
	api.JSONOK(c, gin.H{
		"info": res{user, p.Name, p.Nickname, p.Avatar, p.Sex, p.Birthday, p.Qq, p.DeleteTime},
	})
}

func Statistic(c *gin.Context) {
	commonStr := "SELECT GROUP_CONCAT(distinct user_id SEPARATOR '_') FROM fa_user_open WHERE user_id=fa_user.id"
	// 退款次数
	refundTimesStr := "(SELECT COUNT(*) FROM fa_order_goods WHERE lock_state=1 AND user_id IN (" + commonStr + "))"
	// 退款金额
	refundTotalStr := "(SELECT SUM(goods_pay_price) FROM fa_order_goods WHERE lock_state=1 AND user_id IN (" + commonStr + "))"
	// 购买次数，计算总订单的所有的已付款的购买次数
	buyTimesStr := "(SELECT COUNT(*) FROM fa_order WHERE state>=20 AND uid IN (" + commonStr + "))"
	// 客单价(平均消费)，计算子订单的未退款的平均消费
	costAverageStr := "(SELECT TRUNCATE(IFNULL(AVG(goods_pay_price),0),2) FROM fa_order_goods WHERE lock_state=0 AND user_id=fa_user.id AND order_id IN (SELECT id FROM fa_order WHERE user_id IN (" + commonStr + ") AND state>=20))"
	// 累计消费，计算子订单的未退款的累计订单金额
	costTotalStr := "(SELECT SUM(goods_pay_price) FROM fa_order_goods WHERE lock_state=0 AND user_id=fa_user.id AND order_id IN (SELECT id FROM fa_order WHERE user_id IN (" + commonStr + ") AND state>=20))"
	field := fmt.Sprintf("id,%s AS refund_times,%s AS refund_total,%s AS buy_times,%s AS cost_average,%s AS cost_total", refundTimesStr, refundTotalStr, buyTimesStr, costAverageStr, costTotalStr)

	type result struct {
		ID          string  `json:"id"`
		RefundTimes int     `json:"refund_times"`
		RefundTotal float32 `json:"refund_total"`
		BuyTimes    int     `json:"buy_times"`
		CostAverage float32 `json:"cost_average"`
		CostTotal   float32 `json:"cost_total"`
	}
	var r result
	if e := model.Db.Table("fa_user").Select(field).Where("id = ?", c.Query("id")).Scan(&r).Error; e != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c, gin.H{
		"info": r,
	})
}

func Address(c *gin.Context) {
	var conds = make(mysql.Conds)
	if c.Query("id") == "" {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	conds["uid"] = mysql.Cond{
		"exp",
		"in (SELECT origin_user_id from fa_user_open where user_id = " + c.Query("id") + " group by origin_user_id)"}
	total, res := service.Address.ListPage(c, conds, trans.ReqPage{Sort: "id desc"})
	api.JSONOK(c, gin.H{
		"total_number": total,
		"list":         res,
	})
}
