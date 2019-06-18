package mysql

// 2
// depthStr
// :
// "城市"
// firstLetter
// :
// "d"
// id
// :
// 110101
// jianpin
// :
// "dcq"
// name
// :
// "东城区"
// parentId
// :
// 110000
// pinyin
// :
// "dongchengqu"
type Region struct {
	Id          uint32 `json:"id"`
	Pid         uint32 `json:"parentId"`
	Depth       uint32 `json:"depth"`
	DepthName   string `json:"depthStr"`
	FirstLetter string `json:"firstLetter"`
	Jianpin     string `json:"jianpin"`
	Name        string `json:"name"`
	Pinyin      string `json:"pinyin"`
}

func (Region) TableName() string {
	return "region"
}
