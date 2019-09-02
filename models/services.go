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

func WithCompany(pepper, hmacKey string) ServicesConfig {

	return func(s *Services) error {

		s.Company = NewCompanyService(s.db, pepper, hmacKey)
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
	Company CompanyService
	OAuth   OAuthService
	db      *gorm.DB
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automatically migrate the
// all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&Company{}, &JobPost{}, &pwReset{}, &OAuth{}).Error
}

// DestructiveReset drops the all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&Company{}, &JobPost{}, &pwReset{}, &OAuth{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
