package service

var (
	Invoice                 *invoice
	System                  *system
	Plugin                  *plugin
	SaleNum                 *saleNum
	Article                 *article
	Sms                     *sms
	OrderGoods              *orderGoods
	Shipper                 *shipper
	Express                 *express
	Area                    *area
	GoodsCart               *goodsCart
	Cron                    *cron
	OrderStatis             *orderStatis
	OrderExtend             *orderExtend
	Order                   *order
	GroupConfig             *groupConfig
	DiscountGoods           *discountGoods
	Goods                   *goods
	Coupon                  *coupon
	OrderPay                *orderPay
	Upload                  *upload
	SmsProvider             *smsProvider
	Biz                     *biz
	PayLog                  *payLog
	Message                 *message
	Visit                   *visit
	GoodsSpecValue          *goodsSpecValue
	WechatBroadcast         *wechatBroadcast
	AuthGroup               *authGroup
	Wechat                  *wechat
	Freight                 *freight
	GoodsEvaluate           *goodsEvaluate
	Image                   *image
	UserOpen                *userOpen
	Dispatch                *dispatch
	OrderRefundLog          *orderRefundLog
	MessageType             *messageType
	Page                    *page
	GoodsCategoryIds        *goodsCategoryIds
	VerifyCode              *verifyCode
	UserProfile             *userProfile
	UserAssets              *userAssets
	SendArea                *sendArea
	Fd                      *fd
	SmsScene                *smsScene
	UserAccount             *userAccount
	Group                   *group
	InfoCategory            *infoCategory
	Shop                    *shop
	Queue                   *queue
	PdCash                  *pdCash
	GoodsSpec               *goodsSpec
	WechatAutoReply         *wechatAutoReply
	GoodsImage              *goodsImage
	OrderRefund             *orderRefund
	UserLevel               *userLevel
	GoodsCategory           *goodsCategory
	TransportExtend         *transportExtend
	Extend                  *extend
	Material                *material
	AuthGroupAccess         *authGroupAccess
	OffpayArea              *offpayArea
	Address                 *address
	Discount                *discount
	AccessToken             *accessToken
	GoodsSku                *goodsSku
	GroupGoods              *groupGoods
	Version                 *version
	UserPointsLog           *userPointsLog
	UserAlias               *userAlias
	UserTemp                *userTemp
	CouponGoods             *couponGoods
	OrderRefundReason       *orderRefundReason
	PdRecharge              *pdRecharge
	OrderLog                *orderLog
	AuthRule                *authRule
	MessageState            *messageState
	Setting                 *setting
	Transport               *transport
	GoodsCollect            *goodsCollect
	WechatUser              *wechatUser
	WechatAutoReplyKeywords *wechatAutoReplyKeywords
	FullcutGoods            *fullcutGoods
	Cart                    *cart
	Fullcut                 *fullcut
	PdLog                   *pdLog
	UserVisit               *userVisit
)

func InitGen() {
	FullcutGoods = InitFullcutGoods()
	Cart = InitCart()
	Fullcut = InitFullcut()
	PdLog = InitPdLog()
	Transport = InitTransport()
	GoodsCollect = InitGoodsCollect()
	WechatUser = InitWechatUser()
	WechatAutoReplyKeywords = InitWechatAutoReplyKeywords()
	UserVisit = InitUserVisit()
	Invoice = InitInvoice()
	System = InitSystem()
	Sms = InitSms()
	OrderGoods = InitOrderGoods()
	Shipper = InitShipper()
	Plugin = InitPlugin()
	SaleNum = InitSaleNum()
	Article = InitArticle()
	Cron = InitCron()
	OrderStatis = InitOrderStatis()
	OrderExtend = InitOrderExtend()
	Order = InitOrder()
	Express = InitExpress()
	Area = InitArea()
	GoodsCart = InitGoodsCart()
	OrderPay = InitOrderPay()
	Upload = InitUpload()
	SmsProvider = InitSmsProvider()
	Biz = InitBiz()
	GroupConfig = InitGroupConfig()
	DiscountGoods = InitDiscountGoods()
	Goods = InitGoods()
	Coupon = InitCoupon()
	PayLog = InitPayLog()
	Message = InitMessage()
	Visit = InitVisit()
	GoodsSpecValue = InitGoodsSpecValue()
	WechatBroadcast = InitWechatBroadcast()
	Image = InitImage()
	UserOpen = InitUserOpen()
	Dispatch = InitDispatch()
	OrderRefundLog = InitOrderRefundLog()
	AuthGroup = InitAuthGroup()
	Wechat = InitWechat()
	Freight = InitFreight()
	GoodsEvaluate = InitGoodsEvaluate()
	MessageType = InitMessageType()
	Page = InitPage()
	UserAssets = InitUserAssets()
	SendArea = InitSendArea()
	GoodsCategoryIds = InitGoodsCategoryIds()
	VerifyCode = InitVerifyCode()
	UserProfile = InitUserProfile()
	Fd = InitFd()
	SmsScene = InitSmsScene()
	Shop = InitShop()
	UserAccount = InitUserAccount()
	Group = InitGroup()
	InfoCategory = InitInfoCategory()
	WechatAutoReply = InitWechatAutoReply()
	GoodsImage = InitGoodsImage()
	OrderRefund = InitOrderRefund()
	UserLevel = InitUserLevel()
	Queue = InitQueue()
	PdCash = InitPdCash()
	GoodsSpec = InitGoodsSpec()
	Material = InitMaterial()
	AuthGroupAccess = InitAuthGroupAccess()
	GoodsCategory = InitGoodsCategory()
	TransportExtend = InitTransportExtend()
	Extend = InitExtend()
	AccessToken = InitAccessToken()
	GoodsSku = InitGoodsSku()
	GroupGoods = InitGroupGoods()
	OffpayArea = InitOffpayArea()
	Address = InitAddress()
	Discount = InitDiscount()
	Version = InitVersion()
	UserPointsLog = InitUserPointsLog()
	UserAlias = InitUserAlias()
	OrderLog = InitOrderLog()
	AuthRule = InitAuthRule()
	MessageState = InitMessageState()
	Setting = InitSetting()
	UserTemp = InitUserTemp()
	CouponGoods = InitCouponGoods()
	OrderRefundReason = InitOrderRefundReason()
	PdRecharge = InitPdRecharge()

}
