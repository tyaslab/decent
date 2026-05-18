package entity

import "time"

type Book struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"not null;size:200" json:"title" validate:"required,max=200"`
	Author        string    `gorm:"not null;size:100" json:"author" validate:"required,max=100"`
	Description   string    `gorm:"size:1000" json:"description" validate:"max=1000"`
	PublishedYear int       `json:"year" validate:"min=1000,max=9999"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
