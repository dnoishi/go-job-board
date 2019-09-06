package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		//db.LogMode(true)
		s.db = db
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {

	return func(s *Services) error {

		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithJobPost() ServicesConfig {

	return func(s *Services) error {
		s.JobPost = NewJobPostService(s.db)
		return nil
	}
}

func WithOAuth() ServicesConfig {
	return func(s *Services) error {

		s.OAuth = NewOAuthService(s.db)
		return nil
	}
}
func WithLogMode(logMode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(logMode)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil

}

type Services struct {
	JobPost JobPostService
	User    UserService
	OAuth   OAuthService
	db      *gorm.DB
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automatically migrate
// all tables
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(&User{}, &Role{}, &Location{}, &JobPost{}, &pwReset{}, &OAuth{}).Error
	if err != nil {
		return err
	}
	return runPopulatingFuncs(s.seedRoles, s.seedLocations)
}
func (s *Services) seedRoles() error {
	return s.db.Model(&Role{}).Create(&Role{RoleName: "Company"}).Create(&Role{RoleName: "Candidate"}).Error
}

func (s *Services) seedLocations() error {
	return s.db.Model(&Location{}).
		Create(&Location{LocationName: "USA"}).
		Create(&Location{LocationName: "Canada"}).
		Create(&Location{LocationName: "Europe"}).
		Create(&Location{LocationName: "Remote"}).Error
}

// DestructiveReset drops the all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(
		&User{},
		&Role{},
		&JobPost{},
		&Location{},
		&pwReset{},
		&OAuth{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

type populatingFunc func() error

func runPopulatingFuncs(fns ...populatingFunc) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
