package mysql

import (
	"database/sql/driver"
	"encoding/json"
)

// type GoodsBodysJson GoodsBodyJson
type GoodsBodyJson []struct {
	Type string `json:"type"`
	Val  struct {
		Content string `json:"content"`
	} `json:"value"`
}

type GoodsSpecList struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	ValueList []CreateSpecValue `json:"value_list"`
}

// type GoodsSpecListsJson []GoodsSpecListJson
type GoodsSpecListJson []struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	ValueList []CreateSpecValue `json:"value_list"`
}

type CreateSpecValue struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// type GoodsSkuSpecsJson []GoodsSkuSpecJson
type GoodsSkuSpecJson []struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ValueID   int    `json:"value_id"`
	ValueName string `json:"value_name"`
	ValueImg  string `json:"value_img"`
}

type IntsJson []int

func (c IntsJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *IntsJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type StringsJson []string

func (c StringsJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *StringsJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type GoodsSkusJson []GoodsSku

func (c GoodsSkusJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *GoodsSkusJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type FreightAreasJson []Area

// FullcutHierarchyJson TODO
type FullcutHierarchyJson []string

type GoodsImageSpecImagesJson StringsJson

type GoodsEvaluateImagesJson StringsJson

type GoodsEvaluateAdditionalImagesJson StringsJson

// MaterialMediaJson TODO
type MaterialMediaJson StringsJson

type OrderExtendOrderGoodsJson []Goods

// OrderExtendReciverInfoJson TODO
type OrderExtendReciverInfoJson struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	CombineDetail string `json:"combineDetail"`
}

// OrderExtendInvoiceInfoJson TODO
type OrderExtendInvoiceInfoJson StringsJson

type OrderGoodsGoodsSpecJson []GoodsSpecList

type OrderRefundUserImagesJson StringsJson

type OrderRefundTrackingImagesJson StringsJson

type OrderRefundGoodsSpecJson []Goods

type PageBodyJson string

type SendAreaAreaIdsJson IntsJson

type SettingConfigJson struct {
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

type SmsProviderConfigJson string

type UserOpenInfoAggregateJson struct {
	OpenID    string                             `json:"openId"`
	Nickname  string                             `json:"nickName"`
	Gender    int                                `json:"gender"`
	Province  string                             `json:"province"`
	Language  string                             `json:"language"`
	Country   string                             `json:"country"`
	City      string                             `json:"city"`
	Avatar    string                             `json:"avatarUrl"`
	UnionID   string                             `json:"unionId"`
	Watermark UserOpenInfoAggregateJsonWatermark `json:"watermark"`
}

type UserOpenInfoAggregateJsonWatermark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

type UserWechatTagidListJson struct {
}

type WechatAutoReplySubscribeReplayContentJson struct{}

type WechatAutoReplyReplyContentJson struct{}

type WechatAutoReplyKeysJson struct{}

type WechatBroadcastConditionJson struct{}
type WechatBroadcastSendContentJson struct{}
type WechatBroadcastOpenidsJson struct{}
type WechatUserTagidListJson struct{}
type AuthGroupRuleIdsJson []int
type GoodsImagesJson []string
type GoodsCategoryIdsJson []int
type GoodsSkuListJson []GoodsSku
