package statics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
)

func Quantity(c *gin.Context) {
	resp := trans.RespStaticsQuantity{
		NoSendCount: noSendCount(),
		DayTotal:    dayTotal(),
	}
	api.JSONOK(c, gin.H{
		"info": resp,
	})
}

const (
	DaySeconds = 3600 * 24
)

func getMonthData(date string) (beginMonth int64, endMonth int64, monthDays int64) {
	var t = now.New(time.Now())
	if date != "" {
		tm, e := time.Parse("2006-01", date)
		if e == nil {
			t = now.New(tm)
		}
	}
	beginMonth = t.BeginningOfMonth().Unix()
	endMonth = t.EndOfMonth().Unix()
	monthDays = (endMonth-beginMonth)/DaySeconds + 1
	return
}

func MonthSalesHistogram(c *gin.Context) {
	beginMonth, endMonth, monthDays := getMonthData(c.Query("date"))
	sql := `
select 
	sum(t2.goods_pay_price) as sales,t1.days
from
	(select id,payment_time,date_format(FROM_UNIXTIME(payment_time, '%Y-%m-%d %H:%i:%S'), '%Y-%m-%d') as days from .order where state >=20 and payment_time between ? and ? order by payment_time asc) t1
	left join 
	(select goods_pay_price,order_id from .order_goods where lock_state = 0) t2
	on t1.id = t2.order_id 
group by
	t1.days
`
	res := []struct {
		Sales float32 `json:"sales"`
		Days  string  `json:"days"`
	}{}
	if err := model.Db.Raw(sql, beginMonth, endMonth).Scan(&res).Error; err != nil {
		api.JSONErr(c)
		return
	}

	var list = make([]map[string]interface{}, 0)
	for i := 0; i < int(monthDays); i++ {
		m := make(map[string]interface{})
		curDay := time.Unix(beginMonth+int64(i*DaySeconds), 0).Format("2006-01-02")
		m["day"] = i + 1
		m["sale_number"] = 0
		list = append(list, m)
		for _, v := range res {
			if v.Days == curDay {
				m["sale_number"] = v.Sales
			}
		}
	}
	api.JSONList(c, list, len(list))
}

func MonthOrderCountHistogram(c *gin.Context) {
	beginMonth, endMonth, monthDays := getMonthData(c.Query("date"))
	sql := `
select 
	count(t2.order_id) as number,t1.days 
from
	(select id,payment_time,date_format(FROM_UNIXTIME(payment_time, '%Y-%m-%d %H:%i:%S'),'%Y-%m-%d') as days from .order where state >= 20 and payment_time between ? and ? order by payment_time asc) t1
	left join
	(select order_id from .order_goods where lock_state = 0) t2
	on t1.id = t2.order_id
group by
	t1.days,t2.order_id
`
	res := []struct {
		Number int
		Days   string
	}{}
	if err := model.Db.Raw(sql, beginMonth, endMonth).Scan(&res).Error; err != nil {
		api.JSONErr(c)
		return
	}

	var list = make([]map[string]interface{}, 0)
	for i := 0; i < int(monthDays); i++ {
		m := make(map[string]interface{})
		curDay := time.Unix(beginMonth+int64(i*DaySeconds), 0).Format("2006-01-02")
		m["day"] = i + 1
		m["order_number"] = 0
		list = append(list, m)
		for _, v := range res {
			if v.Days == curDay {
				m["order_number"] = v.Number
			}
		}
	}
	api.JSONList(c, list, len(list))
}

func MonthUserAddCountHistogram(c *gin.Context) {
	beginMonth, endMonth, monthDays := getMonthData(c.Query("date"))
	sql := `
select 
	count(uid) as number,date_format(FROM_UNIXTIME(create_time, '%Y-%m-%d %H:%i:%S'),'%Y-%m-%d') as days
from
	.user_open
where
	state = 1 and delete_time <> 0 and uid > 1 and create_time between ? and ?
group by
	days
`
	var res []struct {
		Number int
		Days   string
	}
	if err := model.Db.Raw(sql, beginMonth, endMonth).Scan(&res).Error; err != nil {
		api.JSONErr(c)
		return
	}

	var list = make([]map[string]interface{}, 0)
	for i := 0; i < int(monthDays); i++ {
		m := make(map[string]interface{})
		curDay := time.Unix(beginMonth+int64(i*DaySeconds), 0).Format("2006-01-02")
		m["day"] = i + 1
		m["customer_number"] = 0
		list = append(list, m)
		for _, v := range res {
			if v.Days == curDay {
				m["customer_number"] = v.Number
			}
		}
	}
	api.JSONList(c, list, len(list))
}

func MonthNewUserSalesHistogram(c *gin.Context) {
	beginMonth, endMonth, monthDays := getMonthData(c.Query("date"))
	sql := `
select 
    sum(t2.goods_pay_price) as sales,t1.days 
from
	(select id,uid,payment_time,date_format(FROM_UNIXTIME(payment_time, '%Y-%m-%d %H:%i:%S'),'%Y-%m-%d') as days from .order where state >= 20 and payment_time between '1554076800' and '1556668799' order by payment_time asc) t1
left join
	(select order_id,goods_pay_price from .order_goods where lock_state = 0) t2
	on t1.id = t2.order_id
left join
	(select uid from .user_open where create_time between ? and ?) t3
	on t1.uid = t3.uid
group by
    t1.days  
`
	res := []struct {
		Sales float32
		Days  string
	}{}
	if err := model.Db.Raw(sql, beginMonth, endMonth).Scan(&res).Error; err != nil {
		api.JSONErr(c)
		return
	}

	var list = make([]map[string]interface{}, 0)
	for i := 0; i < int(monthDays); i++ {
		m := make(map[string]interface{})
		curDay := time.Unix(beginMonth+int64(i*DaySeconds), 0).Format("2006-01-02")
		m["day"] = i + 1
		m["cost"] = 0
		list = append(list, m)
		for _, v := range res {
			if v.Days == curDay {
				m["cost"] = v.Sales
			}
		}
	}
	api.JSONList(c, list, len(list))
}

func noSendCount() (cnt int) {
	model.Db.Table("order").Where("state = ? AND refund_state = ? AND lock_state = ?", 20, 0, 0).Count(&cnt)
	return
}

func dayTotal() (total float64) {
	priceSum := struct {
		PriceSum float64
	}{}
	model.Db.Table("order_goods").Select("sum(order_goods.goods_pay_price) as price_sum").
		Joins("LEFT JOIN .order ON order_goods.order_id = order.id").
		Where("order.state >= ? AND order_goods.lock_state = ?", 20, 0).Find(&priceSum)
	total = priceSum.PriceSum
	return
}
