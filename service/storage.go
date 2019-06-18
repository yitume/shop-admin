package service

import (
	"fmt"
	"time"

	"github.com/RangelReale/osin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
)

var schemas = []string{`CREATE TABLE IF NOT EXISTS app (
  id int(255) UNSIGNED NOT NULL AUTO_INCREMENT,
  name varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  secret varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  redirect_uri varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  extra varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  create_time int(11) NOT NULL,
  status int(11) NOT NULL,
  update_time int(11) NOT NULL,
  PRIMARY KEY (id) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`, `CREATE TABLE IF NOT EXISTS authorize (
	client       varchar(255) BINARY NOT NULL,
	code         varchar(255) BINARY NOT NULL PRIMARY KEY,
	expires_in   int(10) NOT NULL,
	scope        varchar(255) NOT NULL,
	redirect_uri varchar(255) NOT NULL,
	state        varchar(255) NOT NULL,
	extra 		 varchar(255) NOT NULL,
	created_at   int(10) NOT NULL
)`, `CREATE TABLE IF NOT EXISTS access (
	client        varchar(255) BINARY NOT NULL,
	authorize     varchar(255) BINARY NOT NULL,
	previous      varchar(255) BINARY NOT NULL,
	access_token  varchar(255) BINARY NOT NULL PRIMARY KEY,
	refresh_token varchar(255) BINARY NOT NULL,
	expires_in    int(10) NOT NULL,
	scope         varchar(255) NOT NULL,
	redirect_uri  varchar(255) NOT NULL,
	extra 		  varchar(255) NOT NULL,
	created_at    int(10) NOT NULL
)`, `CREATE TABLE IF NOT EXISTS refresh (
	token         varchar(255) BINARY NOT NULL PRIMARY KEY,
	access        varchar(255) BINARY NOT NULL
)`, `CREATE TABLE IF NOT EXISTS expires (
	id 		int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	token		varchar(255) BINARY NOT NULL,
	expires_at	int(10) NOT NULL,
	INDEX expires_index (expires_at),
	INDEX token_expires_index (token)
)`,
}

// storage implements interface "github.com/RangelReale/osin".storage and interface "github.com/felipeweb/osin-mysql/storage".storage
type storage struct {
}

// New returns a new mysql storage instance.
func NewStorage() *storage {
	return &storage{}
}

// CreateSchemas creates the schemata, if they do not exist yet in the dataapi. Returns an error if something went wrong.
func (s *storage) CreateSchemas() error {
	for k, schema := range schemas {
		if dbHandler := model.Db.Exec(schema); dbHandler.Error != nil {
			model.Logger.Error(fmt.Sprintf("Error creating schema %d: %s", k, schema))
			return dbHandler.Error
		}
	}
	return nil
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *storage) Clone() osin.Storage {
	return s
}

// Close the resources the storage potentially holds (using Clone for example)
func (s *storage) Close() {
}

// GetClient loads the client by id
func (s *storage) GetClient(aid string) (osin.Client, error) {
	data := mysql.Client{}
	model.Db.Select("aid, secret, redirect_uri").Where("aid=?", aid).Find(&data)
	c := osin.DefaultClient{
		Id:          data.Aid,
		Secret:      data.Secret,
		RedirectUri: data.RedirectUri,
		UserData:    data.Extra,
	}
	return &c, nil
}

// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
func (s *storage) UpdateClient(c osin.Client) error {
	// data := conv.ToStr(c.GetUserData())
	model.Db.Table("app").Where("id=?", c.GetId()).Update(mysql.Ups{
		"secret":       c.GetSecret(),
		"redirect_uri": c.GetRedirectUri(),
		"extra":        c.GetUserData(),
	})
	return nil
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (s *storage) CreateClient(c osin.Client) error {
	create := mysql.Client{
		Aid:         c.GetId(),
		Secret:      c.GetSecret(),
		RedirectUri: c.GetRedirectUri(),
		Extra:       cast.ToString(c.GetUserData()),
	}
	model.Db.Create(&create)
	return nil
}

// RemoveClient removes a client (identified by id) from the dataapi. Returns an error if something went wrong.
func (s *storage) RemoveClient(aid string) (err error) {
	obj := mysql.Client{}
	model.Db.Where("aid=?", aid).Delete(&obj)
	return nil
}

// SaveAuthorize saves authorize data.
func (s *storage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	obj := mysql.Authorize{
		Client:      data.Client.GetId(),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt.Unix(),
		Extra:       cast.ToString(data.UserData),
	}
	model.Db.Create(&obj)
	if err = s.AddExpireAtData(data.Code, data.ExpireAt()); err != nil {
		return err
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var data osin.AuthorizeData
	obj := mysql.Authorize{}
	model.Db.Select("client, code, expires_in, scope, redirect_uri, state, created_at, extra").Where("code=?", code).Find(&obj)
	data = osin.AuthorizeData{
		Code:        obj.Code,
		ExpiresIn:   obj.ExpiresIn,
		Scope:       obj.Scope,
		RedirectUri: obj.RedirectUri,
		State:       obj.State,
		CreatedAt:   time.Unix(obj.CreatedAt, 0),
		UserData:    obj.Extra,
	}
	c, err := s.GetClient(obj.Client)
	if err != nil {
		return nil, err
	}
	if data.ExpireAt().Before(time.Now()) {
		return nil, errors.New(fmt.Sprintf("Token expired at %s.", data.ExpireAt().String()))
	}
	data.Client = c
	return &data, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *storage) RemoveAuthorize(code string) (err error) {
	obj := mysql.Authorize{}
	if err = model.Db.Where("code=?", code).Delete(&obj).Error; err != nil {
		return err
	}
	if err = s.RemoveExpireAtData(code); err != nil {
		return err
	}
	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *storage) SaveAccess(data *osin.AccessData) (err error) {
	prev := ""
	authorizeData := &osin.AuthorizeData{}
	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}
	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}
	extra := cast.ToString(data.UserData)
	tx := model.Db.Begin()
	if data.RefreshToken != "" {
		if err := s.saveRefresh(tx, data.RefreshToken, data.AccessToken); err != nil {
			return err
		}
	}
	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}
	obj := mysql.Access{
		Client:       data.Client.GetId(),
		Authorize:    authorizeData.Code,
		Previous:     prev,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt.Unix(),
		Extra:        extra,
	}
	err = model.Db.Create(&obj).Error
	if err != nil {
		if rbe := tx.Rollback().Error; rbe != nil {
			return rbe
		}
		return err
	}
	if err = s.AddExpireAtData(data.AccessToken, data.ExpireAt()); err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *storage) LoadAccess(code string) (*osin.AccessData, error) {
	var result osin.AccessData
	obj := mysql.Access{}
	err := model.Db.Where("access_token=?", code).Find(&obj).Error
	if err != nil {
		return nil, err
	}
	result.AccessToken = obj.AccessToken
	result.RefreshToken = obj.RefreshToken
	result.ExpiresIn = obj.ExpiresIn
	result.Scope = obj.Scope
	result.RedirectUri = obj.RedirectUri
	result.CreatedAt = time.Unix(obj.CreatedAt, 0)
	result.UserData = obj.Extra
	client, err := s.GetClient(obj.Client)
	if err != nil {
		return nil, err
	}
	result.Client = client
	result.AuthorizeData, _ = s.LoadAuthorize(obj.Authorize)
	prevAccess, _ := s.LoadAccess(obj.Previous)
	result.AccessData = prevAccess
	return &result, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *storage) RemoveAccess(code string) (err error) {
	obj := mysql.Access{}
	err = model.Db.Where("access_token=?", code).Delete(&obj).Error
	if err != nil {
		return err
	}
	if err = s.RemoveExpireAtData(code); err != nil {
		return err
	}
	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *storage) LoadRefresh(code string) (*osin.AccessData, error) {
	obj := mysql.Refresh{}
	err := model.Db.Where("token=?", code).Find(&obj).Error
	if err != nil {
		return nil, err
	}
	return s.LoadAccess(obj.Access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *storage) RemoveRefresh(code string) error {
	obj := mysql.Refresh{}
	err := model.Db.Where("token=?", code).Delete(&obj).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateClientWithInformation Makes easy to create a osin.DefaultClient
func (s *storage) CreateClientWithInformation(id string, secret string, redirectURI string, userData interface{}) osin.Client {
	return &osin.DefaultClient{
		Id:          id,
		Secret:      secret,
		RedirectUri: redirectURI,
		UserData:    userData,
	}
}

func (s *storage) saveRefresh(tx *gorm.DB, refresh, access string) (err error) {
	obj := mysql.Refresh{
		Token:  refresh,
		Access: access,
	}
	err = tx.Create(&obj).Error
	if err != nil {
		if rbe := tx.Rollback().Error; rbe != nil {
			return rbe
		}
		return err
	}
	return
}

// AddExpireAtData add info in expires table
func (s *storage) AddExpireAtData(code string, expireAt time.Time) error {
	obj := mysql.Expires{
		Token:     code,
		ExpiresAt: expireAt.Unix(),
	}
	err := model.Db.Create(&obj).Error
	if err != nil {
		return err
	}
	return nil
}

// RemoveExpireAtData remove info in expires table
func (s *storage) RemoveExpireAtData(code string) error {
	obj := mysql.Expires{}
	err := model.Db.Where("token=?", code).Delete(&obj).Error
	if err != nil {
		return err
	}
	return nil
}
