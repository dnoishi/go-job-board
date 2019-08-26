package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

const (
	OauthDropbox = "dropbox"
)

type OAuth struct {
	gorm.Model
	UserID  uint   `gorm:not null;unique_index:user_id_service`
	Service string `gorm:not null;unique_index:user_id_service`
	oauth2.Token
}

func NewOAuthService(db *gorm.DB) OAuthService {
	return &oauthValidator{
		&oauthGorm{db},
	}
}

type OAuthService interface {
	OAuthDB
}

type oauthValidator struct {
	OAuthDB
}

func (ov *oauthValidator) Create(oauth *OAuth) error {

	err := runOauthValFuncs(
		oauth, ov.userIDRequired, ov.serviceRequired)
	if err != nil {
		return err
	}
	return ov.OAuthDB.Create(oauth)
}

func (ov *oauthValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}

	return ov.OAuthDB.Delete(id)
}

func (ov *oauthValidator) userIDRequired(o *OAuth) error {
	if o.UserID <= 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (ov *oauthValidator) serviceRequired(o *OAuth) error {
	if o.Service == "" {
		return ErrServiceRequired

	}

	return nil
}

type OAuthDB interface {
	Find(userID uint, service string) (*OAuth, error)
	Create(oauth *OAuth) error
	Delete(id uint) error
}

var _ OAuthDB = &oauthGorm{}

type oauthGorm struct {
	db *gorm.DB
}

func (og *oauthGorm) Create(oauth *OAuth) error {
	return og.db.Create(oauth).Error
}

func (og *oauthGorm) Delete(id uint) error {
	oauth := OAuth{Model: gorm.Model{ID: id}}
	return og.db.Unscoped().Delete(&oauth).Error
}

func (og *oauthGorm) Find(userID uint, service string) (*OAuth, error) {
	var oauth OAuth
	db := og.db.Where("user_id = ?", userID).Where("service = ?", service)
	err := first(db, &oauth)

	return &oauth, err

}

type oAuthValFunc func(*OAuth) error

func runOauthValFuncs(oauth *OAuth, fns ...oAuthValFunc) error {
	for _, fn := range fns {
		if err := fn(oauth); err != nil {
			return err
		}
	}

	return nil
}
