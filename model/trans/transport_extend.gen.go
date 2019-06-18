package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

// ReqTransportExtendList 你可以把ReqTransportExtendList嵌套到需要自行修改的结构体中
type ReqTransportExtendList struct {
	ReqPage
	mysql.TransportExtend
}

// ReqTransportExtendCreate 你可以把ReqTransportExtendCreate或mysql.TransportExtend嵌套到需要自行修改的结构体中
type ReqTransportExtendCreate = mysql.TransportExtend

// ReqTransportExtendUpdate 你可以把ReqTransportExtendUpdate或mysql.TransportExtend嵌套到需要自行修改的结构体中
type ReqTransportExtendUpdate = mysql.TransportExtend

// ReqTransportExtendDelete 你可以把ReqTransportExtendDelete嵌套到需要自行修改的结构体中
type ReqTransportExtendDelete struct {
	Id int `json:"id"`
}
