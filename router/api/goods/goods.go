package goods

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/pkg/util"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"
	"git.yitum.com/saas/shop-admin/service"
)

type createMysql struct {
	Id                int
	Title             string
	ImageSpecImages   mysql.GoodsImageSpecImagesJson
	ImageSpecId       int
	Body              mysql.GoodsBodyJson
	Stock             int
	Freight           float64
	SaleTime          int64
	SpecList          mysql.GoodsSpecListJson
	SpecMap           map[int]mysql.GoodsSpecList
	BaseSaleNum       int
	Price             float64
	Img               string
	Skus              mysql.GoodsSkuListJson
	CategoryIds       mysql.GoodsCategoryIdsJson
	Images            mysql.GoodsImagesJson
	FreightTemplateId int
	IsOnSale          int
	PayType           int
}

func Info(c *gin.Context) {
	req := struct {
		trans.ReqGoodsList
		CategoryIds mysql.GoodsCategoryIdsJson `json:"category_ids" form:"category_ids[]" ` // 商品分类
	}{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	value, err := service.Goods.InfoX(c,
		mysql.Conds{
			"id":          req.Id,
			"open_id":     auth.Default(c).Id,
			"delete_time": 0,
		})
	if err != nil {
		api.JSON(c, api.MsgErr, "info error")
		return
	}
	resp := trans.Goods{
		Id:                 value.Id,
		Title:              value.Title,
		Images:             value.Images,
		CategoryIds:        value.CategoryIds,
		BaseSaleNum:        value.BaseSaleNum,
		Body:               value.Body,
		IsOnSale:           value.IsOnSale,
		ImageSpecId:        value.ImageSpecId,
		ImageSpecImages:    value.ImageSpecImages,
		SkuList:            value.SkuList,
		CreateTime:         value.CreateTime,
		Price:              value.Price,
		UpdateTime:         value.UpdateTime,
		EvaluationCount:    value.EvaluationCount,
		EvaluationGoodStar: value.EvaluationGoodStar,
		Stock:              value.Stock,
		SaleNum:            value.SaleNum,
		GroupSaleNum:       value.GroupSaleNum,
		SaleTime:           value.SaleTime,
		SpecList:           value.SpecList,
		Img:                value.Img,
		PayType:            value.PayType,
		FreightFee:         value.FreightFee,
		FreightId:          value.FreightId,
	}
	api.JSONOK(c, gin.H{
		"info": resp,
	})

}

func List(c *gin.Context) {
	req := trans.ReqGoodsList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	gids := make([]int, 0)
	gidsM := make(map[int]bool, 0)
	conds := mysql.Conds{}

	if len(req.CategoryIds) != 0 {
		conds["category_id"] = []int(req.CategoryIds)
	}
	goodsCategoryIds, _ := service.GoodsCategoryIds.List(c, conds)
	for _, v := range goodsCategoryIds {
		if gidsM[v.GoodsId] {
			continue
		}
		gidsM[v.GoodsId] = true
		gids = append(gids, v.GoodsId)
	}

	total, list := service.Goods.ListPage(c, mysql.Conds{
		"id": gids,
	}, req.ReqPage)
	resp := make([]trans.Goods, 0)
	for _, value := range list {
		resp = append(resp, trans.Goods{
			Id:                 value.Id,
			Title:              value.Title,
			Images:             value.Images,
			CategoryIds:        value.CategoryIds,
			BaseSaleNum:        value.BaseSaleNum,
			Body:               value.Body,
			IsOnSale:           value.IsOnSale,
			ImageSpecId:        value.ImageSpecId,
			ImageSpecImages:    value.ImageSpecImages,
			SkuList:            value.SkuList,
			CreateTime:         value.CreateTime,
			Price:              value.Price,
			UpdateTime:         value.UpdateTime,
			EvaluationCount:    value.EvaluationCount,
			EvaluationGoodStar: value.EvaluationGoodStar,
			Stock:              value.Stock,
			SaleNum:            value.SaleNum,
			GroupSaleNum:       value.GroupSaleNum,
			SaleTime:           value.SaleTime,
			SpecList:           value.SpecList,
			Img:                value.Img,
			PayType:            value.PayType,
			FreightFee:         value.FreightFee,
			FreightId:          value.FreightId,
		})
	}
	api.JSONList(c, resp, total)
}

func Create(c *gin.Context) {
	req := trans.ReqGoodsCreateOrUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	create, err := filterCreateParam(req)
	if err != nil {
		api.JSON(c, api.MsgErr, "filterCreateParam marshal error")
		return
	}

	tx := model.Db.Begin()
	ret := mysql.Goods{create.Id, create.Title,
		create.Images, create.CategoryIds,
		create.BaseSaleNum, create.Body,
		create.IsOnSale, create.ImageSpecId,
		create.ImageSpecImages, create.Skus,
		time.Now().Unix(), create.Price,
		time.Now().Unix(), 0, 0, create.Stock, 0, 0,
		create.SaleTime, 0, create.SpecList, create.Img, create.PayType, create.Freight, create.FreightTemplateId, auth.Default(c).Id}
	err = service.Goods.Create(c, tx, &ret)
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "create error")
		return
	}
	id := ret.Id

	for i := range create.Skus {
		_, specSign, specValueSign, title := skuSign(create.Skus[i].Spec)
		create.Skus[i].GoodsId = id
		create.Skus[i].Title = title
		create.Skus[i].OpenId = auth.Default(c).Id
		create.Skus[i].SpecSign = specSign
		create.Skus[i].SpecValueSign = specValueSign
		err = service.GoodsSku.Create(c, tx, &create.Skus[i])
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}

	for _, img := range create.Images {
		err = service.GoodsImage.Create(c, tx, &mysql.GoodsImage{
			Img:     img,
			GoodsId: id,
			OpenId:  auth.Default(c).Id,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}

	for _, cid := range create.CategoryIds {
		err = service.GoodsCategoryIds.Create(c, tx, &mysql.GoodsCategoryIds{
			GoodsId:    id,
			CategoryId: cid,
			OpenId:     auth.Default(c).Id,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}

	tx.Commit()
	api.JSONOK(c)

}

func Update(c *gin.Context) {
	req := trans.ReqGoodsCreateOrUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	update, err := filterCreateParam(req)
	if err != nil {
		api.JSON(c, api.MsgErr, "filterCreateParam marshal error")
		return
	}

	tx := model.Db.Begin()
	err = service.Goods.Update(c, tx, update.Id, mysql.Ups{
		"title":             update.Title,
		"images":            update.Images,
		"category_ids":      update.CategoryIds,
		"base_sale_num":     update.BaseSaleNum,
		"body":              update.Body,
		"is_on_sale":        update.IsOnSale,
		"image_spec_id":     update.ImageSpecId,
		"image_spec_images": update.ImageSpecImages,
		"sku_list":          update.Skus,
		"price":             update.Price,
		"update_time":       time.Now().Unix(),
		"stock":             update.Stock,
		"sale_time":         update.SaleTime,
		"spec_list":         update.SpecList,
		"img":               update.Img,
		"pay_type":          update.PayType,
		"freight_fee":       update.Freight,
		"freight_id":        update.FreightTemplateId,
	})
	if err != nil {
		tx.Rollback()
		return
	}

	for i := range update.Skus {
		var specSign []int
		var specValueSign []int
		for _, v := range update.Skus[i].Spec {
			specSign = append(specSign, v.ID)
			specValueSign = append(specValueSign, v.ValueID)
		}
		specSignJ, _ := json.Marshal(specSign)
		specValueSignJ, _ := json.Marshal(specValueSign)
		// update.Skus[i].GoodsId = id
		update.Skus[i].Title = update.Title
		update.Skus[i].OpenId = auth.Default(c).Id
		update.Skus[i].SpecSign = string(specSignJ)
		update.Skus[i].SpecValueSign = string(specValueSignJ)
		err = service.GoodsSku.Create(c, tx, &update.Skus[i])
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}
	for _, sku := range update.Skus {
		err = service.GoodsSku.Update(c, tx, sku.Id, mysql.Ups{
			"title": update.Title,
		})
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = service.GoodsImage.Delete(c, tx, update.Id)
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "create error")
		return
	}

	for _, img := range update.Images {
		// TODO 需要Create吗？
		err = service.GoodsImage.Create(c, tx, &mysql.GoodsImage{
			GoodsId: update.Id,
			OpenId:  auth.Default(c).Id,
			Img:     img,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}

	err = service.GoodsCategoryIds.Delete(c, tx, update.Id)
	if err != nil {
		tx.Rollback()
		api.JSON(c, api.MsgErr, "create error")
		return
	}

	for _, cid := range update.CategoryIds {
		// TODO 需要Create吗？
		err = service.GoodsCategoryIds.Create(c, tx, &mysql.GoodsCategoryIds{
			GoodsId:    update.Id,
			OpenId:     auth.Default(c).Id,
			CategoryId: cid,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "create error")
			return
		}
	}
	tx.Commit()
	api.JSONOK(c)
}

func OnSale(c *gin.Context) {
	req := trans.ReqGoodsOnSale{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	tx := model.Db.Begin()
	for _, id := range req.Ids {
		// TODO  auth.Default(c).Id open_id
		err := service.Goods.Update(c, tx, id, mysql.Ups{
			"is_on_sale": 1,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "update error")
			return
		}
	}
	tx.Commit()
	api.JSONOK(c)
}

func OffSale(c *gin.Context) {
	req := trans.ReqGoodsOnSale{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	tx := model.Db.Begin()
	for _, id := range req.Ids {
		// TODO auth.Default(c).Id open_id
		err := service.Goods.Update(c, tx, id, mysql.Ups{
			"is_on_sale": 0,
		})
		if err != nil {
			tx.Rollback()
			api.JSON(c, api.MsgErr, "update error")
			return
		}
	}
	tx.Commit()
	api.JSONOK(c)
}

func filterCreateParam(req trans.ReqGoodsCreateOrUpdate) (resp createMysql, err error) {
	resp.Id = req.Id
	resp.Title = req.Title
	resp.Images = req.Images
	if len(req.Images) > 0 {
		resp.Img = req.Images[0]
	}

	// todo 需要验证分类id是否存在
	for _, value := range req.CategoryIds {
		resp.CategoryIds = append(resp.CategoryIds, cast.ToInt(value))
	}

	resp.IsOnSale = 0
	resp.BaseSaleNum = 0
	resp.FreightTemplateId = req.FreightId
	resp.Freight = req.FreightFee
	resp.SaleTime = req.SaleTime
	resp.Body = req.Body
	resp.ImageSpecId = 0
	resp.Skus = make([]mysql.GoodsSku, 0)
	resp.SpecMap = make(map[int]mysql.GoodsSpecList)
	resp.SpecList = make(mysql.GoodsSpecListJson, 0)
	if len(req.Skus) > 0 {
		staticValueExistIds := make([]int, 0)
		for key, sku := range req.Skus {
			// 初始化每个sku的图片（下面需要处理：如果选择了图片规格，自动设置为图片规格封面）
			sku.Img = resp.Img
			resp.Skus = append(resp.Skus, sku)

			// 价格
			if key == 0 {
				resp.Price = sku.Price
			} else if resp.Price > sku.Price {
				resp.Price = sku.Price
			}
			// 库存
			resp.Stock = sku.Stock

			// 规格层级集合json 不要重复，如色彩下有:xx色 xx 色
			for _, spec := range sku.Spec {
				if !funk.Contains(staticValueExistIds, spec.ValueID) {
					// 存着防止 sku 规格循环被重复记录
					staticValueExistIds = append(staticValueExistIds, spec.ValueID)
					if _, ok := resp.SpecMap[spec.ID]; !ok {
						resp.SpecMap[spec.ID] = mysql.GoodsSpecList{
							Id:        0,
							Name:      "",
							ValueList: make([]mysql.CreateSpecValue, 0),
						}
					}
					specTmp := resp.SpecMap[spec.ID]
					specTmp.Id = spec.ID
					specTmp.Name = spec.Name
					specTmp.ValueList = append(specTmp.ValueList, mysql.CreateSpecValue{
						Id:   spec.ValueID,
						Name: spec.ValueName,
					})

					// 规格图片,防止0或空进来 当默认没规格时是会为空
					if resp.ImageSpecId > 0 && resp.ImageSpecId == spec.ID {
						resp.ImageSpecImages = append(resp.ImageSpecImages, spec.ValueImg)
					}
					resp.SpecMap[spec.ID] = specTmp
				}
			}
		}
		for _, value := range resp.SpecMap {
			resp.SpecList = append(resp.SpecList, value)
		}
		if len(resp.Skus) == 0 || len(resp.SpecList) == 0 {
			err = errors.New("skulist or speclist error")
			return
		}
	}

	resp.ImageSpecImages = resp.ImageSpecImages
	resp.Body = resp.Body
	resp.SpecList = resp.SpecList
	resp.Skus = resp.Skus
	resp.CategoryIds = resp.CategoryIds
	resp.Images = resp.Images

	return
}

func SkuList(c *gin.Context) {
	var req struct {
		trans.ReqPage
		GoodsId int `form:"goods_id"`
	}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	total, res := service.GoodsSku.ListPage(c, mysql.Conds{
		"goods_id": req.GoodsId,
	}, req.ReqPage)
	api.JSONList(c, res, total)
}

func skuSign(specs mysql.GoodsSkuSpecJson) (spec string, specSign string, specValueSign string, title string) {
	ids := make([]int, 0)
	valueIds := make([]int, 0)
	valueName := make([]string, 0)
	for _, value := range specs {
		ids = append(ids, value.ID)
		valueIds = append(valueIds, value.ValueID)
		valueName = append(valueName, value.ValueName)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	sort.Slice(valueIds, func(i, j int) bool {
		return valueIds[i] < valueIds[j]
	})
	spec = util.JsonMarshal(specs)
	specSign = util.JsonMarshal(ids)
	specValueSign = util.JsonMarshal(valueIds)
	title = strings.Join(valueName, " ")
	return
}
