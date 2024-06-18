package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosSupplier struct {
	SupplierID   uuid.UUID  `gorm:"type:uuid;primary_key" json:"supplier_id"`
	SupplierName string     `gorm:"type:varchar(255);not null" json:"supplier_name"`
	ContactName  string     `gorm:"type:varchar(255)" json:"contact_name"`
	ContactEmail string     `gorm:"type:varchar(255)" json:"contact_email"`
	ContactPhone string     `gorm:"type:varchar(20)" json:"contact_phone"`
	BranchID     *uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt    time.Time  `gorm:"type:timestamp" json:"created_at"`
	CreatedBy    uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt    time.Time  `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy    uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
}
