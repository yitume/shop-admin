package mysql

type UrlMenu struct {
	ID         uint32    `json:"id"`
	Pid        uint32    `json:"pid"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	Icon       string    `json:"icon"`
	SortedNum  uint32    `json:"sortedNum"`
	Status     uint32    `json:"status"`
	DefaultURL string    `json:"defaultUrl"`
	CreatedAt  int64     `json:"createdAt"`
	UpdatedAt  int64     `json:"updatedAt"`
	Children   []UrlMenu `gorm:"-"json:"children"`
	MenuType   int       `json:"menuType"`
}

func (UrlMenu) TableName() string {
	return "menu"
}
