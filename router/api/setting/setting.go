package setting

import (
	"time"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func Info(c *gin.Context) {
	var req trans.ReqSettingList
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	value, err := service.Setting.Info(c, req.Key)
	if err != nil {
		api.JSON(c, api.MsgErr, "info error"+err.Error())
		return
	}
	resp := gin.H{
		"key":    value.Key,
		"name":   value.Name,
		"config": value.Config.Data,
		"status": value.Status,
		"remark": value.Remark,
	}
	api.JSONOK(c, gin.H{"info": resp})
}

func Update(c *gin.Context) {
	var req struct {
		trans.ReqSettingUpdate
		ConfigData interface{} `json:"config" form:"config"`
	}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	value, err := service.Setting.Info(c, req.Key)
	if err != nil {
		api.JSON(c, api.MsgErr, "info error"+err.Error())
		return
	}

	// not exist
	if value.CreateTime == 0 {
		if err := service.Setting.Create(c, model.Db, &mysql.Setting{
			Key:    req.Key,
			Name:   req.Name,
			Config: mysql.SettingConfigJson{req.Key, req.ConfigData},
			Status: req.Status,
		}); err != nil {
			api.JSON(c, api.MsgErr, "创建失败")
			return
		}
	} else {
		if err := service.Setting.Update(c, model.Db, req.Key, mysql.Ups{
			"config":      mysql.SettingConfigJson{req.Key, req.ConfigData},
			"status":      req.Status,
			"update_time": time.Now().Unix(),
		}); err != nil {
			api.JSON(c, api.MsgErr, "更新失败")
			return
		}
	}

	api.JSONOK(c, gin.H{"info": value})
}
