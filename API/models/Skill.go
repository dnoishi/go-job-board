package models

import "github.com/jinzhu/gorm"

type Skill struct {
	gorm.Model
	SkillName string
}

type SkillsService interface {
	SkillDB
}

type skillsService struct {
	SkillDB
}

func NewSkillService(db *gorm.DB) SkillsService {
	return &skillsService{
		SkillDB: &skillsValidator{
			&skillsGorm{db},
		},
	}
}

type SkillDB interface {
	FindAll() ([]Skill, error)
	AddSkillToOwner(owner interface{}, skill Skill) error
	DeleteSkillFromOwner(owner interface{}, skill Skill) error
}

type skillsValidator struct {
	SkillDB
}

func (sv *skillsValidator) AddSkillToOwner(owner interface{}, skill Skill) error {
	err := runSkillValFuncs(
		&skill,
		sv.skillIDRequired,
	)
	if err != nil {
		return err
	}

	return sv.SkillDB.AddSkillToOwner(owner, skill)
}

func (sv *skillsValidator) skillIDRequired(s *Skill) error {
	if s.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

var _ SkillDB = &skillsGorm{}

type skillsGorm struct {
	db *gorm.DB
}

func (sg skillsGorm) FindAll() ([]Skill, error) {
	var skills []Skill

	err := sg.db.Find(&skills).Error
	if err != nil {
		return nil, err
	}

	return skills, nil
}
func (sg skillsGorm) AddSkillToOwner(owner interface{}, skill Skill) error {
	return sg.db.Model(owner).Association("Skills").Append(skill).Error
}

func (sg skillsGorm) DeleteSkillFromOwner(owner interface{}, skill Skill) error {
	return sg.db.Model(owner).Association("Skills").Delete(skill).Error
}

type skillValFunc func(skill *Skill) error

func runSkillValFuncs(skill *Skill, fns ...skillValFunc) error {
	for _, fn := range fns {
		if err := fn(skill); err != nil {
			return err
		}
	}

	return nil
}
