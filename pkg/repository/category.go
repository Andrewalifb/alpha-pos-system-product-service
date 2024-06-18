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

type PosProductCategoryRepository interface {
	CreatePosProductCategory(posProductCategory *entity.PosProductCategory) error
	ReadPosProductCategory(categoryID string) (*pb.PosProductCategory, error)
	IsCategoryExist(categoryID string) (bool, error)
	UpdatePosProductCategory(posProductCategory *entity.PosProductCategory) (*pb.PosProductCategory, error)
	DeletePosProductCategory(categoryID string) error
	ReadAllPosProductCategories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posProductCategoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosProductCategoryRepository(db *gorm.DB, redis *redis.Client) PosProductCategoryRepository {
	return &posProductCategoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posProductCategoryRepository) CreatePosProductCategory(posProductCategory *entity.PosProductCategory) error {
	result := r.db.Create(posProductCategory)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (r *posProductCategoryRepository) ReadPosProductCategory(categoryID string) (*pb.PosProductCategory, error) {
	// Try to get the category from Redis first
	categoryData, err := r.redis.Get(context.Background(), categoryID).Result()
	if err == redis.Nil {
		// Category not found in Redis, get from PostgreSQL
		var posProductCategoryEntity entity.PosProductCategory
		if err := r.db.Where("category_id = ?", categoryID).First(&posProductCategoryEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosProductCategory to pb.PosProductCategory
		posProductCategory := &pb.PosProductCategory{
			CategoryId:   posProductCategoryEntity.CategoryID.String(),
			CategoryName: posProductCategoryEntity.CategoryName,
			CompanyId:    posProductCategoryEntity.CompanyID.String(),
			CreatedAt:    timestamppb.New(posProductCategoryEntity.CreatedAt),
			CreatedBy:    posProductCategoryEntity.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posProductCategoryEntity.UpdatedAt),
			UpdatedBy:    posProductCategoryEntity.UpdatedBy.String(),
		}

		// Store the category in Redis for future queries
		categoryData, err := json.Marshal(posProductCategoryEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), categoryID, categoryData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posProductCategory, nil
	} else if err != nil {
		return nil, err
	}

	// Category found in Redis, unmarshal the data
	var posProductCategoryEntity entity.PosProductCategory
	err = json.Unmarshal([]byte(categoryData), &posProductCategoryEntity)
	if err != nil {
		return nil, err
	}

	posProductCategory := &pb.PosProductCategory{
		CategoryId:   posProductCategoryEntity.CategoryID.String(),
		CategoryName: posProductCategoryEntity.CategoryName,
		CompanyId:    posProductCategoryEntity.CompanyID.String(),
		CreatedAt:    timestamppb.New(posProductCategoryEntity.CreatedAt),
		CreatedBy:    posProductCategoryEntity.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posProductCategoryEntity.UpdatedAt),
		UpdatedBy:    posProductCategoryEntity.UpdatedBy.String(),
	}

	return posProductCategory, nil
}

func (r *posProductCategoryRepository) IsCategoryExist(categoryID string) (bool, error) {
	var categoryEntity entity.PosProductCategory
	if err := r.db.Where("category_id = ?", categoryID).First(&categoryEntity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *posProductCategoryRepository) UpdatePosProductCategory(posProductCategory *entity.PosProductCategory) (*pb.PosProductCategory, error) {
	if err := r.db.Save(posProductCategory).Error; err != nil {
		return nil, err
	}

	protoProductCategory := &pb.PosProductCategory{
		CategoryId:   posProductCategory.CategoryID.String(),
		CategoryName: posProductCategory.CategoryName,
		CompanyId:    posProductCategory.CompanyID.String(),
		CreatedAt:    timestamppb.New(posProductCategory.CreatedAt),
		CreatedBy:    posProductCategory.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(posProductCategory.UpdatedAt),
		UpdatedBy:    posProductCategory.UpdatedBy.String(),
	}
	// Update the category in Redis
	categoryData, err := json.Marshal(posProductCategory)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), posProductCategory.CategoryID.String(), categoryData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return protoProductCategory, nil
}

func (r *posProductCategoryRepository) DeletePosProductCategory(categoryID string) error {
	if err := r.db.Where("category_id = ?", categoryID).Delete(&entity.PosProductCategory{}).Error; err != nil {
		return err
	}

	// Delete the category from Redis
	err := r.redis.Del(context.Background(), categoryID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posProductCategoryRepository) ReadAllPosProductCategories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posProductCategories []entity.PosProductCategory
	var totalRecords int64

	query := r.db.Model(&entity.PosProductCategory{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case storeRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posProductCategories)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posProductCategories,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
