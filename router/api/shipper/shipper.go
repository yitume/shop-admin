package shipper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := trans.ReqShipperList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	// TODO 从内存缓存取
	req.ReqPage.PageSize = 10000
	total, list := service.Shipper.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)
}

func Add(c *gin.Context) {
	req := trans.ReqShipperCreate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	areas, err := service.Area.GetAddr(c, req.AreaId)
	if err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.Shipper.Create(c, model.Db, &mysql.Shipper{
		Name:          req.Name,
		ProvinceId:    areas[2].Id,
		CityId:        areas[1].Id,
		AreaId:        areas[0].Id,
		CombineDetail: fmt.Sprintf("{%s} {%s} {%s}", areas[2].Name, areas[1].Name, areas[0].Name),
		Address:       req.Address,
		ContactNumber: req.ContactNumber,
	}); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c)
}

func Info(c *gin.Context) {
	res, err := service.Shipper.Info(c, cast.ToInt(c.Query("id")))
	if err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c, gin.H{
		"info": res,
	})
}

func Edit(c *gin.Context) {
	req := trans.ReqShipperCreate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	areas, err := service.Area.GetAddr(c, req.AreaId)
	if err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.Shipper.Update(c, model.Db, req.Id, mysql.Ups{
		"name":           req.Name,
		"province_id":    areas[2].Id,
		"city_id":        areas[1].Id,
		"area_id":        areas[0].Id,
		"combine_detail": fmt.Sprintf("{%s} {%s} {%s}", areas[2].Name, areas[1].Name, areas[0].Name),
		"address":        req.Address,
		"contact_number": req.ContactNumber,
	}); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c)
}

func Del(c *gin.Context) {
	req := trans.ReqShipperDelete{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	if err := service.Shipper.Delete(c, model.Db, req.Id); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	api.JSONOK(c)
}
