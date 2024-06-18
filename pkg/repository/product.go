package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/entity"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosProductRepository interface {
	CreatePosProduct(posProduct *entity.PosProduct) error
	ReadPosProduct(productID string) (*pb.PosProduct, error)
	ReadPosProductBarcode(productID string) (*pb.PosProduct, error)
	UpdatePosProduct(posProduct *entity.PosProduct) error
	DeletePosProduct(productID string) error
	ReadAllPosProducts(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posProductRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosProductRepository(db *gorm.DB, redis *redis.Client) PosProductRepository {
	return &posProductRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posProductRepository) CreatePosProduct(posProduct *entity.PosProduct) error {
	result := r.db.Create(posProduct)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posProductRepository) ReadPosProduct(productID string) (*pb.PosProduct, error) {
	// Try to get the product from Redis first
	productData, err := r.redis.Get(context.Background(), productID).Result()
	if err == redis.Nil {
		// Product not found in Redis, get from PostgreSQL
		var posProductEntity entity.PosProduct
		if err := r.db.Where("product_id = ?", productID).First(&posProductEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosProduct to pb.PosProduct
		posProduct := &pb.PosProduct{
			ProductId:          posProductEntity.ProductID.String(),
			ProductBarcodeId:   posProductEntity.ProductBarcodeID,
			ProductName:        posProductEntity.ProductName,
			Price:              posProductEntity.Price,
			CostPrice:          posProductEntity.CostPrice,
			CategoryId:         posProductEntity.CategoryID.String(),
			SubCategoryId:      posProductEntity.SubCategoryID.String(),
			StockQuantity:      int32(posProductEntity.StockQuantity),
			ReorderLevel:       int32(posProductEntity.ReorderLevel),
			SupplierId:         posProductEntity.SupplierID.String(),
			ProductDescription: posProductEntity.ProductDescription,
			Active:             posProductEntity.Active,
			StoreId:            posProductEntity.StoreID.String(),
			BranchId:           posProductEntity.BranchID.String(),
			CompanyId:          posProductEntity.CompanyID.String(),
			CreatedAt:          timestamppb.New(posProductEntity.CreatedAt),
			CreatedBy:          posProductEntity.CreatedBy.String(),
			UpdatedAt:          timestamppb.New(posProductEntity.UpdatedAt),
			UpdatedBy:          posProductEntity.UpdatedBy.String(),
		}

		// Store the product in Redis for future queries
		productData, err := json.Marshal(posProductEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), productID, productData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posProduct, nil
	} else if err != nil {
		return nil, err
	}

	// Product found in Redis, unmarshal the data
	var posProductEntity entity.PosProduct
	err = json.Unmarshal([]byte(productData), &posProductEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosProduct to pb.PosProduct
	posProduct := &pb.PosProduct{
		ProductId:          posProductEntity.ProductID.String(),
		ProductBarcodeId:   posProductEntity.ProductBarcodeID,
		ProductName:        posProductEntity.ProductName,
		Price:              posProductEntity.Price,
		CostPrice:          posProductEntity.CostPrice,
		CategoryId:         posProductEntity.CategoryID.String(),
		SubCategoryId:      posProductEntity.SubCategoryID.String(),
		StockQuantity:      int32(posProductEntity.StockQuantity),
		ReorderLevel:       int32(posProductEntity.ReorderLevel),
		SupplierId:         posProductEntity.SupplierID.String(),
		ProductDescription: posProductEntity.ProductDescription,
		Active:             posProductEntity.Active,
		StoreId:            posProductEntity.StoreID.String(),
		BranchId:           posProductEntity.BranchID.String(),
		CompanyId:          posProductEntity.CompanyID.String(),
		CreatedAt:          timestamppb.New(posProductEntity.CreatedAt),
		CreatedBy:          posProductEntity.CreatedBy.String(),
		UpdatedAt:          timestamppb.New(posProductEntity.UpdatedAt),
		UpdatedBy:          posProductEntity.UpdatedBy.String(),
	}

	return posProduct, nil
}

func (r *posProductRepository) ReadPosProductBarcode(productBarcodeID string) (*pb.PosProduct, error) {
	// Try to get the product from Redis first
	productData, err := r.redis.Get(context.Background(), productBarcodeID).Result()
	if err == redis.Nil {
		// Product not found in Redis, get from PostgreSQL
		var posProductEntity entity.PosProduct
		if err := r.db.Where("product_barcode_id = ?", productBarcodeID).First(&posProductEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosProduct to pb.PosProduct
		posProduct := &pb.PosProduct{
			ProductId:          posProductEntity.ProductID.String(),
			ProductBarcodeId:   posProductEntity.ProductBarcodeID,
			ProductName:        posProductEntity.ProductName,
			Price:              posProductEntity.Price,
			CostPrice:          posProductEntity.CostPrice,
			CategoryId:         posProductEntity.CategoryID.String(),
			SubCategoryId:      posProductEntity.SubCategoryID.String(),
			StockQuantity:      int32(posProductEntity.StockQuantity),
			ReorderLevel:       int32(posProductEntity.ReorderLevel),
			SupplierId:         posProductEntity.SupplierID.String(),
			ProductDescription: posProductEntity.ProductDescription,
			Active:             posProductEntity.Active,
			StoreId:            posProductEntity.StoreID.String(),
			BranchId:           posProductEntity.BranchID.String(),
			CompanyId:          posProductEntity.CompanyID.String(),
			CreatedAt:          timestamppb.New(posProductEntity.CreatedAt),
			CreatedBy:          posProductEntity.CreatedBy.String(),
			UpdatedAt:          timestamppb.New(posProductEntity.UpdatedAt),
			UpdatedBy:          posProductEntity.UpdatedBy.String(),
		}

		// Store the product in Redis for future queries
		productData, err := json.Marshal(posProductEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), productBarcodeID, productData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posProduct, nil
	} else if err != nil {
		return nil, err
	}

	// Product found in Redis, unmarshal the data
	var posProductEntity entity.PosProduct
	err = json.Unmarshal([]byte(productData), &posProductEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosProduct to pb.PosProduct
	posProduct := &pb.PosProduct{
		ProductId:          posProductEntity.ProductID.String(),
		ProductBarcodeId:   posProductEntity.ProductBarcodeID,
		ProductName:        posProductEntity.ProductName,
		Price:              posProductEntity.Price,
		CostPrice:          posProductEntity.CostPrice,
		CategoryId:         posProductEntity.CategoryID.String(),
		SubCategoryId:      posProductEntity.SubCategoryID.String(),
		StockQuantity:      int32(posProductEntity.StockQuantity),
		ReorderLevel:       int32(posProductEntity.ReorderLevel),
		SupplierId:         posProductEntity.SupplierID.String(),
		ProductDescription: posProductEntity.ProductDescription,
		Active:             posProductEntity.Active,
		StoreId:            posProductEntity.StoreID.String(),
		BranchId:           posProductEntity.BranchID.String(),
		CompanyId:          posProductEntity.CompanyID.String(),
		CreatedAt:          timestamppb.New(posProductEntity.CreatedAt),
		CreatedBy:          posProductEntity.CreatedBy.String(),
		UpdatedAt:          timestamppb.New(posProductEntity.UpdatedAt),
		UpdatedBy:          posProductEntity.UpdatedBy.String(),
	}

	return posProduct, nil
}

func (r *posProductRepository) UpdatePosProduct(posProduct *entity.PosProduct) error {
	if err := r.db.Save(posProduct).Error; err != nil {
		return err
	}

	// Update the product in Redis
	productData, err := json.Marshal(posProduct)
	if err != nil {
		return err
	}
	err = r.redis.Set(context.Background(), posProduct.ProductID.String(), productData, 7*24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posProductRepository) DeletePosProduct(productID string) error {
	if err := r.db.Where("product_id = ?", productID).Delete(&entity.PosProduct{}).Error; err != nil {
		return err
	}

	// Delete the product from Redis
	err := r.redis.Del(context.Background(), productID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posProductRepository) ReadAllPosProducts(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posProducts []entity.PosProduct
	var totalRecords int64

	query := r.db.Model(&entity.PosProduct{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	case storeRole:
		query = query.Where("store_id = ?", jwtPayload.StoreId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posProducts)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posProducts,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
