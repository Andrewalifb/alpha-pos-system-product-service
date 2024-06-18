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

type PosSupplierRepository interface {
	CreatePosSupplier(posSupplier *entity.PosSupplier) error
	ReadPosSupplier(supplierID string) (*pb.PosSupplier, error)
	UpdatePosSupplier(posSupplier *entity.PosSupplier) (*pb.PosSupplier, error)
	DeletePosSupplier(supplierID string) error
	ReadAllPosSuppliers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posSupplierRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosSupplierRepository(db *gorm.DB, redis *redis.Client) PosSupplierRepository {
	return &posSupplierRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posSupplierRepository) CreatePosSupplier(posSupplier *entity.PosSupplier) error {
	result := r.db.Create(posSupplier)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posSupplierRepository) ReadPosSupplier(supplierID string) (*pb.PosSupplier, error) {
	// Try to get the supplier from Redis first
	supplierData, err := r.redis.Get(context.Background(), supplierID).Result()
	if err == redis.Nil {
		// Supplier not found in Redis, get from PostgreSQL
		var posSupplierEntity entity.PosSupplier
		if err := r.db.Where("supplier_id = ?", supplierID).First(&posSupplierEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosSupplier to pb.PosSupplier
		posSupplier := &pb.PosSupplier{
			SupplierId:   posSupplierEntity.SupplierID.String(),
			SupplierName: posSupplierEntity.SupplierName,
			ContactName:  posSupplierEntity.ContactName,
			ContactEmail: posSupplierEntity.ContactEmail,
			ContactPhone: posSupplierEntity.ContactPhone,
			BranchId:     posSupplierEntity.BranchID.String(),
			CompanyId:    posSupplierEntity.CompanyID.String(),
			CreatedAt:    timestamppb.New(posSupplierEntity.CreatedAt),
			CreatedBy:    posSupplierEntity.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posSupplierEntity.UpdatedAt),
			UpdatedBy:    posSupplierEntity.UpdatedBy.String(),
		}

		// Store the supplier in Redis for future queries
		supplierData, err := json.Marshal(posSupplierEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), supplierID, supplierData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posSupplier, nil
	} else if err != nil {
		return nil, err
	}

	// Supplier found in Redis, unmarshal the data
	var posSupplierEntity entity.PosSupplier
	err = json.Unmarshal([]byte(supplierData), &posSupplierEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosSupplier to pb.PosSupplier
	posSupplier := &pb.PosSupplier{
		SupplierId:   posSupplierEntity.SupplierID.String(),
		SupplierName: posSupplierEntity.SupplierName,
		ContactName:  posSupplierEntity.ContactName,
		ContactEmail: posSupplierEntity.ContactEmail,
		ContactPhone: posSupplierEntity.ContactPhone,
		BranchId:     posSupplierEntity.BranchID.String(),
		CompanyId:    posSupplierEntity.CompanyID.String(),
		CreatedAt:    timestamppb.New(posSupplierEntity.CreatedAt),
		CreatedBy:    posSupplierEntity.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posSupplierEntity.UpdatedAt),
		UpdatedBy:    posSupplierEntity.UpdatedBy.String(),
	}

	return posSupplier, nil
}

func (r *posSupplierRepository) UpdatePosSupplier(posSupplier *entity.PosSupplier) (*pb.PosSupplier, error) {
	if err := r.db.Save(posSupplier).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosSupplier back to pb.PosSupplier
	updatedPosSupplier := &pb.PosSupplier{
		SupplierId:   posSupplier.SupplierID.String(),
		SupplierName: posSupplier.SupplierName,
		ContactName:  posSupplier.ContactName,
		ContactEmail: posSupplier.ContactEmail,
		ContactPhone: posSupplier.ContactPhone,
		BranchId:     posSupplier.BranchID.String(),
		CompanyId:    posSupplier.CompanyID.String(),
		CreatedAt:    timestamppb.New(posSupplier.CreatedAt),
		CreatedBy:    posSupplier.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posSupplier.UpdatedAt),
		UpdatedBy:    posSupplier.UpdatedBy.String(),
	}

	// Update the supplier in Redis
	supplierData, err := json.Marshal(posSupplier)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosSupplier.SupplierId, supplierData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosSupplier, nil
}

func (r *posSupplierRepository) DeletePosSupplier(supplierID string) error {
	if err := r.db.Where("supplier_id = ?", supplierID).Delete(&entity.PosSupplier{}).Error; err != nil {
		return err
	}

	// Delete the supplier from Redis
	err := r.redis.Del(context.Background(), supplierID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posSupplierRepository) ReadAllPosSuppliers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posSuppliers []entity.PosSupplier
	var totalRecords int64

	query := r.db.Model(&entity.PosSupplier{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	case storeRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posSuppliers)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posSuppliers,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
