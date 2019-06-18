package goods

import (
	"sync"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func SpecList(c *gin.Context) {
	total, list := service.GoodsSpec.ListPage(c, mysql.Conds{}, trans.ReqPage{Sort: "sort desc"})
	resp := make([]trans.RespGoodsSpecList, 0)
	var wg sync.WaitGroup
	wg.Add(len(list))
	for _, value := range list {
		go func(value mysql.GoodsSpec) {
			// TODO auth.Default(c).Id open_id
			// _, list := service.GoodsSpecValue.ListAllBySpecId(value.Id, "sort desc", auth.Default(c).Id)
			_, list := service.GoodsSpecValue.ListPage(c, mysql.Conds{
				"spec_id": value.Id,
			}, trans.ReqPage{Sort: "sort desc"})
			resp = append(resp, trans.RespGoodsSpecList{
				Id:     value.Id,
				Name:   value.Name,
				Values: list,
			})
			wg.Done()
		}(value)
	}
	wg.Wait()
	api.JSONList(c, resp, total)
}

func SpecCreate(c *gin.Context) {
	req := trans.ReqGoodsspecCreate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.GoodsSpec.Create(c, model.Db, &mysql.GoodsSpec{
		Name: req.Name,
	}); err != nil {
		api.JSON(c, api.MsgErr, "create error")
		return
	}
	api.JSONOK(c)
}
