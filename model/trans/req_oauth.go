package trans

type ReqOauthLogin struct {
	Identity string `json:"identity" binding:"required"`
	Nn       string `json:"username"`
	Pwd      string `json:"password" binding:"required"`
	Params   struct {
		Redirect     string `json:"redirect"`     // redirect by ant design
		RedirectUri  string `json:"redirect_uri"` // redirect by backend
		ClientId     string `json:"client_id"`
		ResponseType string `json:"response_type"`
	} `form:"params" binding:"required"`
}
