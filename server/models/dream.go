package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Dream struct {
	gorm.Model
	Dream string `gorm:"type:text;not null" json:"dream"`
}

// MarshalJSON implements custom JSON marshaling
func (d Dream) MarshalJSON() ([]byte, error) {
	type Alias Dream
	return json.Marshal(&struct {
		Alias
		ID        uint   `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     Alias(d),
		ID:        d.ID,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	})
}
