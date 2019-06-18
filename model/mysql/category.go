package mysql

type Category struct {
	Id        uint32     `json:"id"`
	Icon      string     `json:"icon"`
	Name      string     `json:"name"`
	Pid       uint32     `json:"pid"`
	Status    uint32     `json:"status"`
	SortedNum uint32     `json:"sortedNum"`
	CreatedAt int64      `json:"createdAt"`
	UpdatedAt int64      `json:"updatedAt"`
	CreatedBy uint32     `json:"createdBy"`
	UpdatedBy uint32     `json:"updatedBy"`
	OpenId    uint32     `json:"-"`
	Children  []Category `gorm:"-"json:"children"`
}

func (Category) TableName() string {
	return "category"
}

type CategorySelect struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

func (CategorySelect) TableName() string {
	return "category"
}
