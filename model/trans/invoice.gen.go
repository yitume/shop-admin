package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

// ReqInvoiceList 你可以把ReqInvoiceList嵌套到需要自行修改的结构体中
type ReqInvoiceList struct {
	ReqPage
	mysql.Invoice
}

// ReqInvoiceCreate 你可以把ReqInvoiceCreate或mysql.Invoice嵌套到需要自行修改的结构体中
type ReqInvoiceCreate = mysql.Invoice

// ReqInvoiceUpdate 你可以把ReqInvoiceUpdate或mysql.Invoice嵌套到需要自行修改的结构体中
type ReqInvoiceUpdate = mysql.Invoice

// ReqInvoiceDelete 你可以把ReqInvoiceDelete嵌套到需要自行修改的结构体中
type ReqInvoiceDelete struct {
	Id int `json:"id"`
}
