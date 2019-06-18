package group

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := trans.ReqGroupList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	total, list := service.Group.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)
}

func PageGoodsList(c *gin.Context) {
	req := trans.ReqGroupGoodsList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	// total, groupList := service.GroupGoods.ListPage(c, req, auth.Default(c).Id)
	total, groupList := service.GroupGoods.ListPage(c, mysql.Conds{}, req.ReqPage)
	if len(groupList) == 0 {
		api.JSONList(c, groupList, total)
		return
	}

	// todo
	groupIds := make([]int, 0)
	groupGoodsIds := make([]int, 0)
	for _, group := range groupList {
		groupIds = append(groupIds, group.Id)
		groupGoodsIds = append(groupGoodsIds, group.GoodsId)
	}
	api.JSONList(c, groupList, total)
}

func SelectableGoods(c *gin.Context) {
	req := struct {
		Title       string `form:"title"`
		CategoryIds []int  `form:"category_ids[]"`
		trans.ReqPage
	}{}
	ts := time.Now().Unix()
	conds := mysql.Conds{
		"is_show": 1,
		"start_time": mysql.Cond{
			"exp",
			fmt.Sprintf("=fa_group.start_time AND (fa_group.start_time>%d) OR (fa_group.start_time<=%d AND fa_group.end_time>=%d)", ts, ts, ts),
		},
	}
	groups, err := service.Group.List(c, conds)
	if err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	var gids []int
	for _, g := range groups {
		gids = append(gids, g.Id)
	}

	goodsSpecCond := mysql.Conds{}
	// 查询活动商品ids
	if len(gids) > 0 {
		goods, err := service.Group.List(c, mysql.Conds{
			"group_id": gids,
		})
		if err != nil {
			api.JSON(c, api.MsgErr, "request app list params is error")
			return
		}
		if len(goods) > 0 {
			var goodIds []int
			for _, g := range goods {
				goodIds = append(goodIds, g.Id)
			}
			goodsSpecCond["ids"] = mysql.Cond{
				"not in", goodIds,
			}
		}
	}

	goodsSpecCond["is_on_sale"] = 1
	if req.Title != "" {
		goodsSpecCond["title"] = req.Title
	}
	if len(req.CategoryIds) != 0 {
		goodsSpecCond["category_ids"] = req.CategoryIds
	}
	total, goodsSpec := service.GoodsSpec.ListPage(c, goodsSpecCond, req.ReqPage)
	api.JSONList(c, goodsSpec, total)
}
