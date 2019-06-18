package mysql

var OrderStates = map[string]int{
	"state_close":      1, // TODO
	"state_new":        OrderStateNew,
	"state_pay":        OrderStatePay,
	"state_send":       OrderStateSend,
	"state_success":    OrderStateSuccess,
	"state_cancel":     OrderStateCancel,
	"state_refund":     1,                    // TODO user表handle_state
	"state_unevaluate": OrderStateUnevaluate, // TODO user表evaluate_state
}

var OrderComments = map[int]string{
	OrderStateNew:     "待付款",
	OrderStatePay:     "待发货",
	OrderStateSend:    "待收货",
	OrderStateSuccess: "交易完成",
	OrderStateCancel:  "已取消",
}

const (
	OrderStateCancel     = 0  // 取消订单
	OrderStateNew        = 10 // 未支付订单
	OrderStatePay        = 20 // 已支付
	OrderStateSend       = 30 // 已发货
	OrderStateSuccess    = 40 // 已收货，交易成功
	OrderStateUnevaluate = 40 // 未评价

)

var OrderGroupComments = map[int]string{
	OrderGroupStateNew:     "待付款",
	OrderGroupStatePay:     "待开团",
	OrderGroupStateSuccess: "拼团成功",
	OrderGroupStateFail:    "拼团失败",
}

const (
	OrderGroupStateNew     = 0 // 正在进行中(待开团)
	OrderGroupStatePay     = 1 // 待付款
	OrderGroupStateSuccess = 2 // 拼团成功
	OrderGroupStateFail    = 3 // 拼团失败
)
