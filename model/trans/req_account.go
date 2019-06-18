package trans

type ReqAccountList struct {
	CurrentPage int    `json:"currentPage" form:"currentPage"`
	PageSize    int    `json:"pageSize"  form:"pageSize"`
	Query       string `json:"query" form:"query"` // 模糊搜索
	Sort        string `json:"sort" form:"sort"`   // 模糊搜索
}

type ReqAccountAdd struct {
	Nickname string `json:"nickname" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
	RePwd    string `json:"repwd" binding:"required"`
}

type ReqAccountUpdate struct {
	Uid      int    `json:"uid"`
	Nickname string `json:"nickname" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
	RePwd    string `json:"repwd" binding:"required"`
}

type ReqRegister struct {
	Email string `json:"mail"`
	Pwd   string `json:"password"`
}
