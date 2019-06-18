package trans

type RespAppProfileInfo struct {
	Uid            uint32 `json:"uid"`
	Nickname       string `json:"nickname"`
	Domain         string `json:"domain"`
	PlatformNo     string `json:"platformNo"`
	PlatformSecret string `json:"platformSecret"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	LastLoginIp    string `json:"lastLoginIp"`
	Telephone      string `json:"telephone"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func (RespAppProfileInfo) TableName() string {
	return "biz"
}

type RespAppWechatInfo struct {
	Id        uint32 `json:"id"`
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (RespAppWechatInfo) TableName() string {
	return "wechat_config"
}
