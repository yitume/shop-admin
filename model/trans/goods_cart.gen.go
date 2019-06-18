package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

// ReqGoodsCartList 你可以把ReqGoodsCartList嵌套到需要自行修改的结构体中
type ReqGoodsCartList struct {
	ReqPage
	mysql.GoodsCart
}

// ReqGoodsCartCreate 你可以把ReqGoodsCartCreate或mysql.GoodsCart嵌套到需要自行修改的结构体中
type ReqGoodsCartCreate = mysql.GoodsCart

// ReqGoodsCartUpdate 你可以把ReqGoodsCartUpdate或mysql.GoodsCart嵌套到需要自行修改的结构体中
type ReqGoodsCartUpdate = mysql.GoodsCart

// ReqGoodsCartDelete 你可以把ReqGoodsCartDelete嵌套到需要自行修改的结构体中
type ReqGoodsCartDelete struct {
	Id int `json:"id"`
}
