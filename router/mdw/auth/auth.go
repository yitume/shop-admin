package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"git.yitum.com/saas/shop-admin/model"
)

var (
	DefaultKey = "git.yitum.com/saas/shop-admin/app/modules/auth"
	SessionKey = "AUTHUNIQUEID"
	// RedirectUrl should be the relative URL for your login route
	RedirectUrl string = "/user/login"

	// RedirectParam is the query string parameter that will be set
	// with the page the user was trying to visit before they were
	// intercepted.
	RedirectParam string = "return_url"
)

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		openId := session.Get(SessionKey)
		auth := &Auth{}
		if openId != nil {
			err := auth.get(openId)
			if err != nil {
				model.Logger.Error("login error", zap.String("err", err.Error()))
			} else {
				auth.login()
			}
		} else {
			model.Logger.Debug("login status No UserId")
		}
		c.Set(DefaultKey, auth)
		c.Next()
	}
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := Default(c)
		if a.IsAuthenticated() == false {
			c.JSON(http.StatusOK, gin.H{
				"code":   401,
				"result": "",
				"msg":    "user not login",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

type Auth struct {
	Id            int    `json:"uid,omitempty" gorm:"column:open_id"`
	Nickname      string `form:"nickname" json:"nickname,omitempty"`
	Avatar        string `form:"avatar" json:"avatar,omitempty"`
	Password      string `form:"password" json:"-"`
	Gender        int64  `json:"gender,omitempty"`
	Birthday      int64  `json:"birthday,omitempty"`
	CreatedAt     int64  `gorm:"column:created_time" json:"created_time,omitempty"`
	UpdatedAt     int64  `gorm:"column:updated_time" json:"updated_time,omitempty"`
	authenticated bool   `form:"-" db:"-" json:"-"`
}

func (Auth) TableName() string {
	return "biz"
}

func (a *Auth) get(openId interface{}) error {
	if err := model.Db.Where("open_id = ?", openId).First(&a).Error; err != nil {
		return err
	}
	return nil
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (a *Auth) login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	a.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (a *Auth) logout() {
	// Remove from logged-in user's list
	// etc ...
	a.authenticated = false
}

func (a *Auth) IsAuthenticated() bool {
	return a.authenticated
}

func (a *Auth) UniqueId() int {
	return a.Id
}

// shortcut to get Auth
func Default(c *gin.Context) *Auth {
	// return c.MustGet(DefaultKey).(auth)
	return c.MustGet(DefaultKey).(*Auth)
}

// AuthenticateSession will mark the session and user object as authenticated. Then
// the Login() user function will be called. This function should be called after
// you have validated a user.
func AuthenticateSession(s sessions.Session, a *Auth) error {
	a.login()
	return UpdateUser(s, a)
}

// UpdateUser updates the User object stored in the session. This is useful incase a change
// is made to the user model that needs to persist across requests.
func UpdateUser(s sessions.Session, a *Auth) error {
	s.Set(SessionKey, a.UniqueId())
	s.Save()
	return nil
}

// Logout will clear out the session and call the Logout() user function.
func Logout(s sessions.Session, a *Auth) {
	a.logout()
	s.Delete(SessionKey)
	s.Save()
}
