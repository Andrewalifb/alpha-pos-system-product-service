package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosInventoryHistory struct {
	InventoryID uuid.UUID  `gorm:"type:uuid;primary_key" json:"inventory_id"`
	ProductID   uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	StoreID     *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	Date        time.Time  `gorm:"type:timestamp;not null" json:"date"`
	Quantity    int        `gorm:"type:int;not null" json:"quantity"`
	BranchID    *uuid.UUID `gorm:"type:uuid" json:"branch_id"`
	CompanyID   uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt   time.Time  `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
}
