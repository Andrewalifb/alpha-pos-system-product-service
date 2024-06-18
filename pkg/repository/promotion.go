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

type PosPromotionRepository interface {
	CreatePosPromotion(posPromotion *entity.PosPromotion) error
	ReadPosPromotion(promotionID string) (*pb.PosPromotion, error)
	UpdatePosPromotion(posPromotion *entity.PosPromotion) (*pb.PosPromotion, error)
	DeletePosPromotion(promotionID string) error
	ReadAllPosPromotions(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
	ReadPosPromotionByProductId(productID string) (*pb.PosPromotion, error) // New method
}

type posPromotionRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosPromotionRepository(db *gorm.DB, redis *redis.Client) PosPromotionRepository {
	return &posPromotionRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posPromotionRepository) CreatePosPromotion(posPromotion *entity.PosPromotion) error {
	result := r.db.Create(posPromotion)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posPromotionRepository) ReadPosPromotion(promotionID string) (*pb.PosPromotion, error) {
	// Try to get the promotion from Redis first
	promotionData, err := r.redis.Get(context.Background(), promotionID).Result()
	if err == redis.Nil {
		// Promotion not found in Redis, get from PostgreSQL
		var posPromotionEntity entity.PosPromotion
		if err := r.db.Where("promotion_id = ?", promotionID).First(&posPromotionEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosPromotion to pb.PosPromotion
		posPromotion := &pb.PosPromotion{
			PromotionId:  posPromotionEntity.PromotionID.String(),
			ProductId:    posPromotionEntity.ProductID.String(),
			StartDate:    posPromotionEntity.StartDate.Format(time.RFC3339),
			EndDate:      posPromotionEntity.EndDate.Format(time.RFC3339),
			Active:       posPromotionEntity.Active,
			DiscountRate: posPromotionEntity.DiscountRate,
			StoreId:      posPromotionEntity.StoreID.String(),
			BranchId:     posPromotionEntity.BranchID.String(),
			CompanyId:    posPromotionEntity.CompanyID.String(),
			CreatedAt:    timestamppb.New(posPromotionEntity.CreatedAt),
			CreatedBy:    posPromotionEntity.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posPromotionEntity.UpdatedAt),
			UpdatedBy:    posPromotionEntity.UpdatedBy.String(),
		}

		// Store the promotion in Redis for future queries
		promotionData, err := json.Marshal(posPromotionEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), promotionID, promotionData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posPromotion, nil
	} else if err != nil {
		return nil, err
	}

	// Promotion found in Redis, unmarshal the data
	var posPromotionEntity entity.PosPromotion
	err = json.Unmarshal([]byte(promotionData), &posPromotionEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosPromotion to pb.PosPromotion
	posPromotion := &pb.PosPromotion{
		PromotionId:  posPromotionEntity.PromotionID.String(),
		ProductId:    posPromotionEntity.ProductID.String(),
		StartDate:    posPromotionEntity.StartDate.Format(time.RFC3339),
		EndDate:      posPromotionEntity.EndDate.Format(time.RFC3339),
		Active:       posPromotionEntity.Active,
		DiscountRate: posPromotionEntity.DiscountRate,
		StoreId:      posPromotionEntity.StoreID.String(),
		BranchId:     posPromotionEntity.BranchID.String(),
		CompanyId:    posPromotionEntity.CompanyID.String(),
		CreatedAt:    timestamppb.New(posPromotionEntity.CreatedAt),
		CreatedBy:    posPromotionEntity.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posPromotionEntity.UpdatedAt),
		UpdatedBy:    posPromotionEntity.UpdatedBy.String(),
	}

	return posPromotion, nil
}

func (r *posPromotionRepository) ReadPosPromotionByProductId(productID string) (*pb.PosPromotion, error) {
	// Try to get the promotion from Redis first
	promotionData, err := r.redis.Get(context.Background(), "product_"+productID).Result()
	if err == redis.Nil {
		// Promotion not found in Redis, get from PostgreSQL
		var posPromotionEntity entity.PosPromotion
		if err := r.db.Where("product_id = ?", productID).First(&posPromotionEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosPromotion to pb.PosPromotion
		posPromotion := &pb.PosPromotion{
			PromotionId:  posPromotionEntity.PromotionID.String(),
			ProductId:    posPromotionEntity.ProductID.String(),
			StartDate:    posPromotionEntity.StartDate.Format(time.RFC3339),
			EndDate:      posPromotionEntity.EndDate.Format(time.RFC3339),
			Active:       posPromotionEntity.Active,
			DiscountRate: posPromotionEntity.DiscountRate,
			StoreId:      posPromotionEntity.StoreID.String(),
			BranchId:     posPromotionEntity.BranchID.String(),
			CompanyId:    posPromotionEntity.CompanyID.String(),
			CreatedAt:    timestamppb.New(posPromotionEntity.CreatedAt),
			CreatedBy:    posPromotionEntity.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posPromotionEntity.UpdatedAt),
			UpdatedBy:    posPromotionEntity.UpdatedBy.String(),
		}

		// Store the promotion in Redis for future queries
		promotionData, err := json.Marshal(posPromotionEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), "product_"+productID, promotionData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posPromotion, nil
	} else if err != nil {
		return nil, err
	}

	// Promotion found in Redis, unmarshal the data
	var posPromotionEntity entity.PosPromotion
	err = json.Unmarshal([]byte(promotionData), &posPromotionEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosPromotion to pb.PosPromotion
	posPromotion := &pb.PosPromotion{
		PromotionId:  posPromotionEntity.PromotionID.String(),
		ProductId:    posPromotionEntity.ProductID.String(),
		StartDate:    posPromotionEntity.StartDate.Format(time.RFC3339),
		EndDate:      posPromotionEntity.EndDate.Format(time.RFC3339),
		Active:       posPromotionEntity.Active,
		DiscountRate: posPromotionEntity.DiscountRate,
		StoreId:      posPromotionEntity.StoreID.String(),
		BranchId:     posPromotionEntity.BranchID.String(),
		CompanyId:    posPromotionEntity.CompanyID.String(),
		CreatedAt:    timestamppb.New(posPromotionEntity.CreatedAt),
		CreatedBy:    posPromotionEntity.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posPromotionEntity.UpdatedAt),
		UpdatedBy:    posPromotionEntity.UpdatedBy.String(),
	}

	return posPromotion, nil
}

func (r *posPromotionRepository) UpdatePosPromotion(posPromotion *entity.PosPromotion) (*pb.PosPromotion, error) {
	if err := r.db.Save(posPromotion).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosPromotion back to pb.PosPromotion
	updatedPosPromotion := &pb.PosPromotion{
		PromotionId:  posPromotion.PromotionID.String(),
		ProductId:    posPromotion.ProductID.String(),
		StartDate:    posPromotion.StartDate.Format(time.RFC3339),
		EndDate:      posPromotion.EndDate.Format(time.RFC3339),
		Active:       posPromotion.Active,
		DiscountRate: posPromotion.DiscountRate,
		StoreId:      posPromotion.StoreID.String(),
		BranchId:     posPromotion.BranchID.String(),
		CompanyId:    posPromotion.CompanyID.String(),
		CreatedAt:    timestamppb.New(posPromotion.CreatedAt),
		CreatedBy:    posPromotion.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posPromotion.UpdatedAt),
		UpdatedBy:    posPromotion.UpdatedBy.String(),
	}

	// Update the promotion in Redis
	promotionData, err := json.Marshal(posPromotion)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosPromotion.PromotionId, promotionData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosPromotion, nil
}

func (r *posPromotionRepository) DeletePosPromotion(promotionID string) error {
	if err := r.db.Where("promotion_id = ?", promotionID).Delete(&entity.PosPromotion{}).Error; err != nil {
		return err
	}

	// Delete the promotion from Redis
	err := r.redis.Del(context.Background(), promotionID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posPromotionRepository) ReadAllPosPromotions(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posPromotions []entity.PosPromotion
	var totalRecords int64

	query := r.db.Model(&entity.PosPromotion{})

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

	query.Find(&posPromotions)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posPromotions,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
