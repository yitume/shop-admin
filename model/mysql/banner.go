package mysql

type Banner struct {
	ID      uint32 `gorm:"primary_key" json:"id"`
	ReferId uint32 `json:"referId"`
	OpenId  uint32 `json:"-"`
	LinkUrl string `json:"linkUrl"`
	PicUrl  string `json:"picUrl"`
	Remark  string `json:"remark"`
	Name    string `json:"name"`

	Status     uint32 `json:"status"`
	BannerType uint32 `json:"bannerType"`
	SortedNum  uint32 `json:"sortedNum"`
	CreatedBy  uint32 `json:"createdBy"`
	UpdatedBy  uint32 `json:"updatedBy"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func (Banner) TableName() string {
	return "banner"
}
