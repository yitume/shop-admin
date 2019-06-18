package trans

type ReqMemberPassword struct {
	OldPassword string `json:"oldpassword"`
	Password    string `json:"password"`
}

type ReqMember struct {
	Password string `json:"password"`
	Username string `json:"Username"`
	Name     string `json:"name"`
	Status   bool   `json:"status"`
}
