package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosProductSubCategory struct {
	SubCategoryID   uuid.UUID `gorm:"type:uuid;primary_key" json:"sub_category_id"`
	SubCategoryName string    `gorm:"type:varchar(255);not null" json:"sub_category_name"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	CompanyID       uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt       time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt       time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy       uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}
