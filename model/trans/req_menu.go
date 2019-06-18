package trans

type ReqMenuList struct {
	CurrentPage int    `json:"currentPage"`
	Search      string `json:"search"`
}

type ReqMenuTreeNode struct {
	Id uint32 `form:"id" binding:"required"`
}

type ReqMenuCreate struct {
	Name       string `json:"name"`
	Pid        uint32 `json:"pid"`
	Url        string `json:"url"`
	SortedNum  uint32 `json:"sortedNum"`
	Icon       string `json:"icon"`
	DefaultUrl string `json:"defaultUrl"`
}

type ReqMenuUpdate struct {
	Id         uint32 `json:"id"`
	Name       string `json:"name"`
	Pid        uint32 `json:"pid"`
	Url        string `json:"url"`
	SortedNum  uint32 `json:"sortedNum"`
	Icon       string `json:"icon"`
	DefaultUrl string `json:"defaultUrl"`
}

type ReqMenuDel struct {
	Id uint32 `json:"id"`
}
