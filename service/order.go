package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
)

type OrderFull struct {
	mysql.Order
	StateDesc         string `json:"state_desc"`
	GroupStateDesc    string `json:"group_state_desc"`
	PaymentName       string `json:"payment_name"`
	mysql.UserOpen    `json:"extend_user"`
	mysql.OrderExtend `json:"extend_order_extend"`
	ExtendOrderGoods  []ExtendOrderGoods `json:"extend_order_goods"`
}

type ExtendOrderGoods struct {
	mysql.OrderGoods
	IfRefund bool `json:"if_refund"`
}

func (order) PaymentName(in string) string {
	switch in {
	case "offline":
		return "货到付款"
	case "online":
		return "在线付款"
	case "alipay":
		return "支付宝"
	case "tenpay":
		return "网银在线"
	case "chinabank":
		return "网银在线"
	case "predeposit":
		return "预存款"
	}
	return ""
}

const (
	// OrderGoodsTypeNormal 1普通订单
	OrderGoodsTypeNormal = 1
	// OrderGoodsTypeGroupon 2团购订单
	OrderGoodsTypeGroupon = 2
	// OrderGoodsTypeDiscount 3限时折扣商品
	OrderGoodsTypeDiscount = 3
	// OrderGoodsTypeComb 4组合套装
	OrderGoodsTypeComb = 4
	// OrderGoodsTypeGift 5赠品
	OrderGoodsTypeGift = 5
)

// GetOrderFulls 改写getOrderInfo
func (order) GetOrderFulls(c *gin.Context, orderIds []int) ([]OrderFull, error) {
	rets := make([]OrderFull, 0, len(orderIds))
	orderGoods, _ := OrderGoods.ListMap(c, mysql.Conds{
		"order_id": orderIds,
	})
	extendOrderGoods, _ := OrderExtend.ListMap(c, mysql.Conds{
		"id": orderIds,
	})

	orders, _ := Order.ListMap(c, mysql.Conds{
		"id": orderIds,
	})

	for _, v := range orders {
		for _, val := range orderGoods {
			if v.Id == val.OrderId {
				ret := OrderFull{
					Order:            v,
					OrderExtend:      extendOrderGoods[v.Id],
					ExtendOrderGoods: make([]ExtendOrderGoods, 0),
				}
				tmp := ExtendOrderGoods{OrderGoods: val}
				// 退款平台处理状态 默认0处理中(未处理) 10拒绝(驳回) 20同意 30成功(已完成) 50取消(用户主动撤销) 51取消(用户主动收货)
				// 不可退款
				if val.RefundId > 0 && funk.Contains([]int{20, 30, 51}, val.RefundHandleState) {
					tmp.IfRefund = false
				} else {
					tmp.IfRefund = true
				}

				// refund_state 0不显示申请退款按钮 1显示申请退款按钮 2显示退款中按钮 3显示退款完成
				refundState := 0
				if v.State <= 10 {
					refundState = 0
				} else {
					if val.LockState == 0 && val.RefundId == 0 {
						refundState = 1
					} else {
						if val.RefundHandleState == 30 {
							refundState = 3
						} else {
							refundState = 2
						}
					}
				}
				ret.ExtendOrderGoods = append(ret.ExtendOrderGoods, tmp)
				ret.RefundState = refundState
				rets = append(rets, ret)
			}
		}
	}

	return rets, nil
}

// GetOrderFullsPage 改写getOrderList
func (order) GetOrderFullsPage(c *gin.Context, reqList trans.ReqPage, conds mysql.Conds, extend map[string]bool) (int, map[int]OrderFull) {
	total, orders := Order.ListPage(c, conds, reqList)
	tmpMap := make(map[int]OrderFull, 0)
	orderIds := make([]int, 0, total)
	uids := make([]int, 0, total)
	uidstmp := make(map[int]bool, 0)
	for _, v := range orders {
		ret := OrderFull{Order: v}
		// 1默认2拼团商品3限时折扣商品4组合套装5赠品
		ret.StateDesc = mysql.OrderComments[v.State]
		if ret.Order.GoodsType == 2 {
			ret.StateDesc = mysql.OrderGroupComments[v.State]
		}
		ret.PaymentName = Order.PaymentName(v.PaymentCode)
		tmpMap[v.Id] = ret

		orderIds = append(orderIds, v.Id)
		if _, ok := uidstmp[v.Uid]; !ok {
			uids = append(uids, v.Uid)
		}
	}

	if _, ok := extend["user"]; ok {
		users, _ := UserOpen.ListMap(c, mysql.Conds{
			"uid": uids,
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			ret.UserOpen = users[ret.UserOpen.Uid]
			tmpMap[oid] = ret
		}
	}

	if _, ok := extend["order_extend"]; ok {
		extendOrderGoods, _ := OrderExtend.ListMap(c, mysql.Conds{
			"id": orderIds,
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			ret.OrderExtend = extendOrderGoods[oid]
			tmpMap[oid] = ret
		}
	}

	if _, ok := extend["order_goods"]; ok {
		orderGoods, _ := OrderGoods.ListMap(c, mysql.Conds{
			"order_id": orderIds,
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			if ret.ExtendOrderGoods == nil {
				ret.ExtendOrderGoods = make([]ExtendOrderGoods, 0)
			}
			ret.ExtendOrderGoods = append(ret.ExtendOrderGoods, ExtendOrderGoods{OrderGoods: orderGoods[oid]})
			tmpMap[oid] = ret
		}
	}

	return total, tmpMap
}
