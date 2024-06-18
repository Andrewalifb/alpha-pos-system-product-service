package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosProduct struct {
	ProductID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"product_id"`
	ProductBarcodeID   string     `gorm:"type:varchar(255);not null" json:"product_barcode_id"`
	ProductName        string     `gorm:"type:varchar(255);not null" json:"product_name"`
	Price              float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	CostPrice          float64    `gorm:"type:decimal(10,2)" json:"cost_price"`
	CategoryID         uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	SubCategoryID      uuid.UUID  `gorm:"type:uuid;not null" json:"sub_category_id"`
	StockQuantity      int        `gorm:"type:int;not null" json:"stock_quantity"`
	ReorderLevel       int        `gorm:"type:int" json:"reorder_level"`
	SupplierID         uuid.UUID  `gorm:"type:uuid" json:"supplier_id"`
	ProductDescription string     `gorm:"type:text" json:"product_description"`
	Active             bool       `gorm:"type:boolean;default:true" json:"active"`
	StoreID            uuid.UUID  `gorm:"type:uuid" json:"store_id"`
	BranchID           *uuid.UUID `gorm:"type:uuid" json:"branch_id"`
	CompanyID          uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt          time.Time  `gorm:"type:timestamp" json:"created_at"`
	CreatedBy          uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt          time.Time  `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy          uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
}
