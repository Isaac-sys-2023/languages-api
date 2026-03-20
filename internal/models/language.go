package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Language struct {
	gorm.Model

	Name          string         `json:"name" gorm:"not null;unique"`
	Description   string         `json:"description" gorm:"type:varchar(255);not null"`
	ReleaseYear   int            `json:"year"`
	Icon          string         `json:"icon"`
	PlainIcon     *string        `json:"plain_icon,omitempty"`
	LineIcon      *string        `json:"line_icon,omitempty"`
	Creator       string         `json:"creator" gorm:"not null"`
	LatestVersion string         `json:"latest_version" gorm:"not null"`
	Color         string         `json:"color" gorm:"not null"`
	HabitualUses  pq.StringArray `json:"habitual_uses" gorm:"type:text[]"`
	Website       *string        `json:"website,omitempty"`
}
