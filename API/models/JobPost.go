package models

import (
	"github.com/jinzhu/gorm"
)

// JobPost represents a job post
type JobPost struct {
	gorm.Model
	UserID           uint   `gorm:"not_null;index"`
	Title            string `gorm:"not_null"`
	Location         string `gorm:"not_null"`
	Category         string `gorm:"not_null"`
	Description      string `gorm:"not_null"`
	ApplicationEmail string `gorm:"not_null"`
}

type JobPostService interface {
	JobPostDB
}

type jobPostService struct {
	JobPostDB
}

func NewJobPostService(db *gorm.DB) JobPostService {
	return &jobPostService{
		JobPostDB: &jobPostValidator{
			&jobPostGorm{db},
		},
	}
}

type JobPostDB interface {
	ByUserID(id uint) ([]JobPost, error)
	ByID(id uint) (*JobPost, error)
	Create(jobPost *JobPost) error
	Update(jobPost *JobPost) error
	Delete(id uint) error
}

type jobPostValidator struct {
	JobPostDB
}

func (gv *jobPostValidator) Create(jobPost *JobPost) error {

	err := runJobPostValFuncs(
		jobPost, gv.userIDRequired, gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.JobPostDB.Create(jobPost)
}

func (gv *jobPostValidator) Update(jobPost *JobPost) error {

	err := runJobPostValFuncs(
		jobPost, gv.userIDRequired, gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.JobPostDB.Update(jobPost)
}

func (gv *jobPostValidator) Delete(id uint) error {

	if id <= 0 {
		return ErrIDInvalid
	}

	return gv.JobPostDB.Delete(id)
}

func (gv *jobPostValidator) userIDRequired(jp *JobPost) error {
	if jp.UserID <= 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (gv *jobPostValidator) titleRequired(jp *JobPost) error {
	if jp.Title == "" {
		return ErrTitleRequired
	}

	return nil
}

var _ JobPostDB = &jobPostGorm{}

type jobPostGorm struct {
	db *gorm.DB
}

// Create will create the provided jobPost and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (jpg *jobPostGorm) Create(jobPost *JobPost) error {
	return jpg.db.Create(jobPost).Error
}

func (jpg *jobPostGorm) Update(jobPost *JobPost) error {
	return jpg.db.Save(jobPost).Error
}

func (jpg *jobPostGorm) Delete(id uint) error {
	jobPost := JobPost{Model: gorm.Model{ID: id}}
	return jpg.db.Delete(&jobPost).Error
}

func (jpg *jobPostGorm) ByID(id uint) (*JobPost, error) {
	var jobPost JobPost
	db := jpg.db.Where("id = ?", id)
	err := first(db, &jobPost)

	return &jobPost, err

}
func (jpg *jobPostGorm) ByUserID(id uint) ([]JobPost, error) {
	var galleries []JobPost
	err := jpg.db.Where("user_id = ?", id).Find(&galleries).Error
	if err != nil {
		return nil, err
	}
	return galleries, nil

}

type jobPostValFunc func(*JobPost) error

func runJobPostValFuncs(jobPost *JobPost, fns ...jobPostValFunc) error {
	for _, fn := range fns {
		if err := fn(jobPost); err != nil {
			return err
		}
	}

	return nil
}
