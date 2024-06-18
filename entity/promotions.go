package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosPromotion struct {
	PromotionID  uuid.UUID  `gorm:"type:uuid;primary_key" json:"promotion_id"`
	ProductID    uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	StartDate    time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time  `gorm:"type:date;not null" json:"end_date"`
	Active       bool       `gorm:"type:boolean;default:true" json:"active"`
	DiscountRate float64    `gorm:"type:decimal(5,2);not null" json:"discount_rate"`
	StoreID      uuid.UUID  `gorm:"type:uuid" json:"store_id"`
	BranchID     *uuid.UUID `gorm:"type:uuid" json:"branch_id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt    time.Time  `gorm:"type:timestamp" json:"created_at"`
	CreatedBy    uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt    time.Time  `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy    uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
}
