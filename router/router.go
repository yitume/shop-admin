package router

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/pkg/bootstrap"
	"git.yitum.com/saas/shop-admin/router/api/area"
	"git.yitum.com/saas/shop-admin/router/api/coupon"
	"git.yitum.com/saas/shop-admin/router/api/express"
	"git.yitum.com/saas/shop-admin/router/api/freight"
	"git.yitum.com/saas/shop-admin/router/api/goods"
	"git.yitum.com/saas/shop-admin/router/api/goodsevaluate"
	"git.yitum.com/saas/shop-admin/router/api/group"
	"git.yitum.com/saas/shop-admin/router/api/image"
	"git.yitum.com/saas/shop-admin/router/api/member"
	"git.yitum.com/saas/shop-admin/router/api/order"
	"git.yitum.com/saas/shop-admin/router/api/orderrefund"
	"git.yitum.com/saas/shop-admin/router/api/page"
	"git.yitum.com/saas/shop-admin/router/api/setting"
	"git.yitum.com/saas/shop-admin/router/api/shipper"
	"git.yitum.com/saas/shop-admin/router/api/shop"
	"git.yitum.com/saas/shop-admin/router/api/statics"
	"git.yitum.com/saas/shop-admin/router/api/user"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"
)

func InitRouter() *gin.Engine {
	gin.SetMode(bootstrap.Conf.App.Mode)
	r := gin.New()
	r.Use(ginzap.Ginzap(model.Logger.Logger, time.RFC3339, true))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(model.Session)
	r.Use(auth.New())

	r.POST("/admin/member/login", member.Login)
	r.POST("/admin/member/add", member.Add)
	r.POST("/admin/member/register", member.Register)
	r.GET("/admin/member/confirm", member.Confirm)

	memberGrp := r.Group("/admin/member")
	memberGrp.Use(auth.LoginRequired())
	{
		memberGrp.GET("/self", member.Self)
		memberGrp.POST("/selfPassword", member.SelfPassword)
		memberGrp.GET("/logout", member.Logout)
	}

	areaGrp := r.Group("/admin/area")
	areaGrp.Use(auth.LoginRequired())
	{
		areaGrp.GET("/list", area.List)
	}

	userGrp := r.Group("/admin/user")
	userGrp.Use(auth.LoginRequired())
	{
		userGrp.GET("/list", user.List)
		userGrp.GET("/info", user.Info)
		userGrp.GET("/statistics", user.Statistic)
		userGrp.GET("/address", user.Address)
	}

	shopGrp := r.Group("/admin/shop")
	shopGrp.Use(auth.LoginRequired())
	{
		shopGrp.GET("/info", shop.Info)
		shopGrp.POST("/setBaseInfo", shop.Update)
		shopGrp.POST("/setGoodsCategoryStyle", shop.SetGoodsCategoryStyle)
		shopGrp.POST("/setOrderExpires", shop.SetOrderExpires)
	}

	pageGrp := r.Group("/admin/page")
	pageGrp.Use(auth.LoginRequired())
	{
		pageGrp.GET("/list", page.List)
		pageGrp.GET("/info", page.Info)
		pageGrp.POST("/add", page.Create)
		pageGrp.POST("/edit", page.Update)
		pageGrp.POST("/setPortal", page.SetPortal)
	}

	orderGrp := r.Group("/admin/order")
	orderGrp.Use(auth.LoginRequired())
	{
		orderGrp.GET("/list", order.List)
		orderGrp.GET("/info", order.Info)
		orderGrp.POST("/changePrice", order.ChangePrice)
		orderGrp.POST("/setSend", order.SetSend)
	}

	orderrefundGrp := r.Group("/admin/orderrefund")
	orderrefundGrp.Use(auth.LoginRequired())
	{
		orderrefundGrp.GET("/list", orderrefund.List)
		orderrefundGrp.GET("/info", orderrefund.Info)
	}
	goodsevaluateGrp := r.Group("/admin/goodsevaluate")
	goodsevaluateGrp.Use(auth.LoginRequired())
	{
		goodsevaluateGrp.GET("/list", goodsevaluate.List)
		goodsevaluateGrp.POST("/display", goodsevaluate.Display)
	}

	imageGrp := r.Group("/admin/image")
	imageGrp.Use(auth.LoginRequired())
	{
		imageGrp.GET("/list", image.List)
		imageGrp.GET("/goodsImageList", image.GoodsList)
		imageGrp.POST("/add", image.Add)
	}

	goodsGrp := r.Group("/admin/goods")
	goodsGrp.Use(auth.LoginRequired())
	{
		goodsGrp.GET("/info", goods.Info)
		goodsGrp.GET("/list", goods.List)
		goodsGrp.POST("/add", goods.Create)
		goodsGrp.POST("/edit", goods.Update)
		goodsGrp.POST("/onSale", goods.OnSale)
		goodsGrp.POST("/offSale", goods.OffSale)
		goodsGrp.GET("/skuList", goods.SkuList)
	}

	goodscategoryGrp := r.Group("/admin/goodscategory")
	goodscategoryGrp.Use(auth.LoginRequired())
	{
		goodscategoryGrp.GET("/list", goods.CategoryList)
		goodscategoryGrp.GET("/info", goods.CategoryInfo)
		goodscategoryGrp.POST("/add", goods.CategoryCreate)
		goodscategoryGrp.POST("/edit", goods.CategoryUpdate)
		goodscategoryGrp.POST("/del", goods.CategoryDel)
	}

	goodsspecGrp := r.Group("/admin/goodsspec")
	goodsspecGrp.Use(auth.LoginRequired())
	{
		goodsspecGrp.GET("/list", goods.SpecList)
		goodsspecGrp.POST("/add", goods.SpecCreate)
	}

	goodsspecvalueGrp := r.Group("/admin/goodsspecvalue")
	goodsspecvalueGrp.Use(auth.LoginRequired())
	{
		goodsspecvalueGrp.POST("/add", goods.SpecvalueCreate)
	}

	freightGrp := r.Group("/admin/freight")
	freightGrp.Use(auth.LoginRequired())
	{
		freightGrp.POST("/list", freight.List)
	}

	groupGrp := r.Group("/admin/group")
	groupGrp.Use(auth.LoginRequired())
	{
		groupGrp.GET("/list", group.List)
		groupGrp.GET("/pageGoods", group.PageGoodsList)
		groupGrp.GET("/selectableGoods", group.SelectableGoods)
	}

	shipperGrp := r.Group("/admin/shipper")
	shipperGrp.Use(auth.LoginRequired())
	{
		shipperGrp.GET("/list", shipper.List)
		shipperGrp.POST("/add", shipper.Add)
		shipperGrp.GET("/info", shipper.Info)
		shipperGrp.POST("/edit", shipper.Edit)
		shipperGrp.POST("/del", shipper.Del)
	}

	expressGrp := r.Group("/admin/express")
	expressGrp.Use(auth.LoginRequired())
	{
		expressGrp.GET("/list", express.List)
	}

	settingGrp := r.Group("/admin/setting")
	settingGrp.Use(auth.LoginRequired())
	{
		settingGrp.GET("/info", setting.Info)
		settingGrp.POST("/edit", setting.Update)
	}

	couponGrp := r.Group("/admin/coupon")
	couponGrp.Use(auth.LoginRequired())
	{
		couponGrp.GET("/list", coupon.List)
	}

	staticsGrp := r.Group("/admin/Statistics")
	staticsGrp.Use(auth.LoginRequired())
	{
		staticsGrp.GET("/quantity", statics.Quantity)
		staticsGrp.GET("/monthSalesHistogram", statics.MonthSalesHistogram)
		staticsGrp.GET("/monthOrderCountHistogram", statics.MonthOrderCountHistogram)
		staticsGrp.GET("/monthUserAddCountHistogram", statics.MonthUserAddCountHistogram)
		staticsGrp.GET("/monthNewUserSalesHistogram", statics.MonthNewUserSalesHistogram)
	}

	return r
}
