package trans

type RespStaticsQuantity struct {
	NoSendCount            int     `json:"no_send_count"`
	DayTotal               float64 `json:"day_total"`
	CostAverage            int     `json:"cost_average"`
	YesterdayNewUser       int     `json:"yesterday_new_user"`
	AllUser                int     `json:"all_user"`
	PositiveCount          int     `json:"positive_count"`
	YesterdayPositiveCount int     `json:"yesterday_positive_count"`
	YesterdayModerateCount int     `json:"yesterday_moderate_count"`
	YesterdayNegativeCount int     `json:"yesterday_negative_count"`
}
