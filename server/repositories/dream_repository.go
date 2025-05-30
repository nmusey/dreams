package repositories

import (
	"dreams/models"

	"gorm.io/gorm"
)

type DreamRepository struct {
	db *gorm.DB
}

func NewDreamRepository(db *gorm.DB) *DreamRepository {
	return &DreamRepository{db: db}
}

func (r *DreamRepository) Create(content string) error {
	dream := models.Dream{
		Content: content,
	}
	return r.db.Create(&dream).Error
}

func (r *DreamRepository) FindAll() ([]models.Dream, error) {
	var dreams []models.Dream
	err := r.db.Find(&dreams).Error
	return dreams, err
}
