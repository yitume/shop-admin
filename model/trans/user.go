package trans

type RespSelf struct {
	Info struct {
		ID         int           `json:"id"`
		Username   string        `json:"username"`
		Phone      interface{}   `json:"phone"`
		Email      interface{}   `json:"email"`
		State      int           `json:"state"`
		Salt       string        `json:"salt"`
		IsDiscard  int           `json:"is_discard"`
		CreateTime int           `json:"create_time"`
		DeleteTime interface{}   `json:"delete_time"`
		Profile    []interface{} `json:"profile"`
		Assets     []interface{} `json:"assets"`
	} `json:"info"`
	Group []interface{} `json:"group"`
	Rules []string      `json:"rules"`
}

type UserInfoResp struct {
	Uid        uint32 `json:"uid"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Avatar     string `json:"avatar"`
	Country    string `json:"country"`
	Geographic struct {
		City struct {
			Key   string `json:"key"`
			Label string `json:"label"`
		} `json:"city"`
		Province struct {
			Key   string `json:"key"`
			Label string `json:"label"`
		} `json:"province"`
	} `json:"geographic"`
	Group       string `json:"group"`
	NotifyCount int    `json:"notifyCount"`
	Phone       string `json:"phone"`
	Signature   string `json:"signature"`
	Tags        []struct {
		Key   string `json:"key"`
		Label string `json:"label"`
	} `json:"tags"`
}
