package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

// ReqUserList 你可以把ReqUserList嵌套到需要自行修改的结构体中
type ReqUserList struct {
	ReqPage
	mysql.UserOpen
}

// ReqUserCreate 你可以把ReqUserCreate或mysql.User嵌套到需要自行修改的结构体中
type ReqUserCreate = mysql.UserOpen

// ReqUserUpdate 你可以把ReqUserUpdate或mysql.User嵌套到需要自行修改的结构体中
type ReqUserUpdate = mysql.UserOpen

// ReqUserDelete 你可以把ReqUserDelete嵌套到需要自行修改的结构体中
type ReqUserDelete struct {
	Uid int `json:"uid"`
}
