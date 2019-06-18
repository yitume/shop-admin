package page

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

type resPage struct {
	mysql.Page
	Body interface{} `json:"body"` // 模板json内容
}

func List(c *gin.Context) {
	req := trans.ReqPageList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	// total, list := service.Page.ListPage(c, req, auth.Default(c).Id)
	total, list := service.Page.ListPage(c, mysql.Conds{}, req.ReqPage)
	res := make([]resPage, 0, len(list))
	for _, v := range list {
		var body interface{}
		if err := json.Unmarshal([]byte(v.Body), &body); err != nil {
			api.JSON(c, api.MsgErr, "info error")
			return
		}
		res = append(res, resPage{
			Page: v,
			Body: v.Body,
		})
	}
	api.JSONList(c, list, total)
}

func Info(c *gin.Context) {
	req := trans.ReqPageList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	value, err := service.Page.Info(c, req.Id)
	if err != nil {
		api.JSON(c, api.MsgErr, "info error")
		return
	}
	var body interface{}
	if err = json.Unmarshal([]byte(value.Body), &body); err != nil {
		api.JSON(c, api.MsgErr, "info error")
		return
	}
	api.JSONOK(c, gin.H{
		"info": resPage{
			Page: value,
			Body: body,
		},
	})
}

func Create(c *gin.Context) {
	create := resPage{}
	if err := c.Bind(&create); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	body, _ := json.Marshal(create.Body)
	tx := model.Db.Begin()
	err := service.Page.Create(c, tx, &mysql.Page{
		Id:              create.Id,
		Name:            create.Name,
		Description:     create.Description,
		Body:            string(body),
		IsPortal:        create.IsPortal,
		IsSystem:        create.IsSystem,
		BackgroundColor: create.BackgroundColor,
		Type:            create.Type,
		CreateTime:      create.CreateTime,
		UpdateTime:      create.UpdateTime,
		Module:          create.Module,
		DeleteTime:      create.DeleteTime,
		CloneFromId:     create.CloneFromId,
	})
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "create error")
		return
	}

	tx.Commit()
	api.JSONOK(c)
}

func Update(c *gin.Context) {
	update := resPage{}
	if err := c.Bind(&update); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	body, _ := json.Marshal(update.Body)

	tx := model.Db.Begin()
	err := service.Page.Update(c, tx, update.Id,
		mysql.Ups{
			"name":             update.Name,
			"description":      update.Description,
			"is_portal":        update.IsPortal,
			"is_system":        update.IsSystem,
			"background_color": update.BackgroundColor,
			"type":             update.Type,
			"update_time":      update.UpdateTime,
			"module":           update.Module,
			"clone_from_id":    update.CloneFromId,
			"body":             body,
		})
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "create error")
		return
	}

	tx.Commit()
	api.JSONOK(c)
}

func SetPortal(c *gin.Context) {
	req := trans.ReqPageUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	var err error
	tx := model.Db.Begin()
	err = service.Page.UpdateX(c, tx, mysql.Conds{}, mysql.Ups{
		"is_portal":   0,
		"update_time": time.Now().Unix(),
	})
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "update error")
		return
	}
	err = service.Page.Update(c, tx, req.Id, mysql.Ups{
		"is_portal":   1,
		"update_time": time.Now().Unix(),
	})
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "update error")
		return
	}

	tx.Commit()
	api.JSONOK(c)
}
