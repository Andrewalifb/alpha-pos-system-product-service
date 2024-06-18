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

type PosInventoryHistoryRepository interface {
	CreatePosInventoryHistory(posInventoryHistory *entity.PosInventoryHistory) error
	ReadPosInventoryHistory(inventoryID string) (*pb.PosInventoryHistory, error)
	UpdatePosInventoryHistory(posInventoryHistory *entity.PosInventoryHistory) error
	DeletePosInventoryHistory(inventoryID string) error
	ReadAllPosInventoryHistories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posInventoryHistoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosInventoryHistoryRepository(db *gorm.DB, redis *redis.Client) PosInventoryHistoryRepository {
	return &posInventoryHistoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posInventoryHistoryRepository) CreatePosInventoryHistory(posInventoryHistory *entity.PosInventoryHistory) error {
	result := r.db.Create(posInventoryHistory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posInventoryHistoryRepository) ReadPosInventoryHistory(inventoryID string) (*pb.PosInventoryHistory, error) {
	// Try to get the inventory history from Redis first
	inventoryHistoryData, err := r.redis.Get(context.Background(), inventoryID).Result()
	if err == redis.Nil {
		// Inventory history not found in Redis, get from PostgreSQL
		var posInventoryHistoryEntity entity.PosInventoryHistory
		if err := r.db.Where("inventory_id = ?", inventoryID).First(&posInventoryHistoryEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosInventoryHistory to pb.PosInventoryHistory
		posInventoryHistory := &pb.PosInventoryHistory{
			InventoryId: posInventoryHistoryEntity.InventoryID.String(),
			ProductId:   posInventoryHistoryEntity.ProductID.String(),
			StoreId:     posInventoryHistoryEntity.StoreID.String(),
			Date:        timestamppb.New(posInventoryHistoryEntity.Date),
			Quantity:    int32(posInventoryHistoryEntity.Quantity),
			BranchId:    posInventoryHistoryEntity.BranchID.String(),
			CompanyId:   posInventoryHistoryEntity.CompanyID.String(),
			CreatedAt:   timestamppb.New(posInventoryHistoryEntity.CreatedAt),
			CreatedBy:   posInventoryHistoryEntity.CreatedBy.String(),
			UpdatedAt:   timestamppb.New(posInventoryHistoryEntity.UpdatedAt),
			UpdatedBy:   posInventoryHistoryEntity.UpdatedBy.String(),
		}

		// Store the inventory history in Redis for future queries
		inventoryHistoryData, err := json.Marshal(posInventoryHistoryEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), inventoryID, inventoryHistoryData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posInventoryHistory, nil
	} else if err != nil {
		return nil, err
	}

	// Inventory history found in Redis, unmarshal the data
	var posInventoryHistoryEntity entity.PosInventoryHistory
	err = json.Unmarshal([]byte(inventoryHistoryData), &posInventoryHistoryEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosInventoryHistory to pb.PosInventoryHistory
	posInventoryHistory := &pb.PosInventoryHistory{
		InventoryId: posInventoryHistoryEntity.InventoryID.String(),
		ProductId:   posInventoryHistoryEntity.ProductID.String(),
		StoreId:     posInventoryHistoryEntity.StoreID.String(),
		Date:        timestamppb.New(posInventoryHistoryEntity.Date),
		Quantity:    int32(posInventoryHistoryEntity.Quantity),
		BranchId:    posInventoryHistoryEntity.BranchID.String(),
		CompanyId:   posInventoryHistoryEntity.CompanyID.String(),
		CreatedAt:   timestamppb.New(posInventoryHistoryEntity.CreatedAt),
		CreatedBy:   posInventoryHistoryEntity.CreatedBy.String(),
		UpdatedAt:   timestamppb.New(posInventoryHistoryEntity.UpdatedAt),
		UpdatedBy:   posInventoryHistoryEntity.UpdatedBy.String(),
	}

	return posInventoryHistory, nil
}

func (r *posInventoryHistoryRepository) UpdatePosInventoryHistory(posInventoryHistory *entity.PosInventoryHistory) error {
	if err := r.db.Save(posInventoryHistory).Error; err != nil {
		return err
	}

	// Update the inventory history in Redis
	inventoryHistoryData, err := json.Marshal(posInventoryHistory)
	if err != nil {
		return err
	}
	err = r.redis.Set(context.Background(), posInventoryHistory.InventoryID.String(), inventoryHistoryData, 7*24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posInventoryHistoryRepository) DeletePosInventoryHistory(inventoryID string) error {
	if err := r.db.Where("inventory_id = ?", inventoryID).Delete(&entity.PosInventoryHistory{}).Error; err != nil {
		return err
	}

	// Delete the inventory history from Redis
	err := r.redis.Del(context.Background(), inventoryID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posInventoryHistoryRepository) ReadAllPosInventoryHistories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posInventoryHistories []entity.PosInventoryHistory
	var totalRecords int64

	query := r.db.Model(&entity.PosInventoryHistory{})

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

	query.Find(&posInventoryHistories)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posInventoryHistories,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
