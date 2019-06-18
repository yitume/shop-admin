package trans

type ReqAppList struct {
	CurrentPage int    `json:"currentPage" form:"currentPage"`
	PageSize    int    `json:"pageSize"  form:"pageSize"`
	Query       string `json:"query" form:"query"` // 模糊搜索
	Sort        string `json:"sort" form:"sort"`   // 模糊搜索
}

type ReqAppAdd struct {
	Name        string `json:"name" form:"name"`               // 模糊搜索
	RedirectUri string `json:"redirectUri" form:"redirectUri"` // 模糊搜索
}

type ReqAppUpdate struct {
	Aid         int    `json:"aid"`
	Name        string `json:"name"`
	RedirectUri string `json:"redirectUri"`
	Status      int    `json:"status"`
}

type ReqAppBannerList struct {
	Current  int    `form:"current"`
	PageSize int    `form:"pageSize"`
	Status   string `form:"status"`
	Name     string `form:"name"` // 模糊搜索
	Sorter   string `form:"sorter"`
}

//
// name: fields.name,
// linkUrl: fields.linkUrl,
// picUrl: fields.img.fileList[0].url,
// status: parseInt(fields.status),
// sortedNum: parseInt(fields.sortedNum),
// remark: fields.remark,
type ReqAppBannerCreate struct {
	Name      string `json:"name"`
	LinkUrl   string `json:"linkUrl"`
	PicUrl    string `json:"picUrl"`
	Status    uint32 `json:"status"`
	SortedNum uint32 `json:"sortedNum"`
	Remark    string `json:"remark"`
}

type ReqAppBannerUpdate struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	LinkUrl   string `json:"linkUrl"`
	PicUrl    string `json:"picUrl"`
	Status    uint32 `json:"status"`
	SortedNum uint32 `json:"sortedNum"`
	Remark    string `json:"remark"`
}

type ReqAppWechatUpdate struct {
	Id        int    `json:"id"`
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type ReqAppProfileUpdate struct {
	Nickname  string `json:"nickname"`
	Domain    string `json:"domain"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}
