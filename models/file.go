package models

import (
	"github.com/GymWorkoutApp/gwap-files/utils/uuid"
	"github.com/jinzhu/gorm"
)

// Client client model
type File struct {
	Base
	ID     		string `json:"id" gorm:"not null;primary_key;"`
	Source 		string `json:"-" gorm:"not null;"`
	Filename 	string `json:"filename" gorm:"not null;"`
}

func (u File) GetID() string {
	return u.ID
}

func (u File) SetID(id string) {
	u.ID = id
}

func (u File) GetSource() string {
	return u.Source
}

func (u File) SetSource(source string) {
	u.Source = source
}

func (u File) GetFilename() string {
	return u.Filename
}

func (u File) SetFilename(filename string) {
	u.Filename = filename
}

func (c *File) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.Must(uuid.NewRandom()).String())
	return nil
}