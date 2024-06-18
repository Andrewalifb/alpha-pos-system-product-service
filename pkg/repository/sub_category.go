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

type PosProductSubCategoryRepository interface {
	CreatePosProductSubCategory(posProductSubCategory *entity.PosProductSubCategory) error
	ReadPosProductSubCategory(subCategoryID string) (*pb.PosProductSubCategory, error)
	UpdatePosProductSubCategory(posProductSubCategory *entity.PosProductSubCategory) (*pb.PosProductSubCategory, error)
	DeletePosProductSubCategory(subCategoryID string) error
	ReadAllPosProductSubCategories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posProductSubCategoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosProductSubCategoryRepository(db *gorm.DB, redis *redis.Client) PosProductSubCategoryRepository {
	return &posProductSubCategoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posProductSubCategoryRepository) CreatePosProductSubCategory(posProductSubCategory *entity.PosProductSubCategory) error {
	result := r.db.Create(posProductSubCategory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posProductSubCategoryRepository) ReadPosProductSubCategory(subCategoryID string) (*pb.PosProductSubCategory, error) {
	subCategoryData, err := r.redis.Get(context.Background(), subCategoryID).Result()
	if err == redis.Nil {
		var posProductSubCategoryEntity entity.PosProductSubCategory
		if err := r.db.Where("sub_category_id = ?", subCategoryID).First(&posProductSubCategoryEntity).Error; err != nil {
			return nil, err
		}

		posProductSubCategory := &pb.PosProductSubCategory{
			SubCategoryId:   posProductSubCategoryEntity.SubCategoryID.String(),
			SubCategoryName: posProductSubCategoryEntity.SubCategoryName,
			CategoryId:      posProductSubCategoryEntity.CategoryID.String(),
			CompanyId:       posProductSubCategoryEntity.CompanyID.String(),
			CreatedAt:       timestamppb.New(posProductSubCategoryEntity.CreatedAt),
			CreatedBy:       posProductSubCategoryEntity.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posProductSubCategoryEntity.UpdatedAt),
			UpdatedBy:       posProductSubCategoryEntity.UpdatedBy.String(),
		}

		subCategoryData, err := json.Marshal(posProductSubCategoryEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), subCategoryID, subCategoryData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posProductSubCategory, nil
	} else if err != nil {
		return nil, err
	}

	var posProductSubCategoryEntity entity.PosProductSubCategory
	err = json.Unmarshal([]byte(subCategoryData), &posProductSubCategoryEntity)
	if err != nil {
		return nil, err
	}

	posProductSubCategory := &pb.PosProductSubCategory{
		SubCategoryId:   posProductSubCategoryEntity.SubCategoryID.String(),
		SubCategoryName: posProductSubCategoryEntity.SubCategoryName,
		CategoryId:      posProductSubCategoryEntity.CategoryID.String(),
		CompanyId:       posProductSubCategoryEntity.CompanyID.String(),
		CreatedAt:       timestamppb.New(posProductSubCategoryEntity.CreatedAt),
		CreatedBy:       posProductSubCategoryEntity.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posProductSubCategoryEntity.UpdatedAt),
		UpdatedBy:       posProductSubCategoryEntity.UpdatedBy.String(),
	}

	return posProductSubCategory, nil
}

func (r *posProductSubCategoryRepository) UpdatePosProductSubCategory(posProductSubCategory *entity.PosProductSubCategory) (*pb.PosProductSubCategory, error) {
	if err := r.db.Save(posProductSubCategory).Error; err != nil {
		return nil, err
	}

	updatedPosProductSubCategory := &pb.PosProductSubCategory{
		SubCategoryId:   posProductSubCategory.SubCategoryID.String(),
		SubCategoryName: posProductSubCategory.SubCategoryName,
		CategoryId:      posProductSubCategory.CategoryID.String(),
		CompanyId:       posProductSubCategory.CompanyID.String(),
		CreatedAt:       timestamppb.New(posProductSubCategory.CreatedAt),
		CreatedBy:       posProductSubCategory.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posProductSubCategory.UpdatedAt),
		UpdatedBy:       posProductSubCategory.UpdatedBy.String(),
	}

	subCategoryData, err := json.Marshal(posProductSubCategory)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosProductSubCategory.SubCategoryId, subCategoryData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosProductSubCategory, nil
}

func (r *posProductSubCategoryRepository) DeletePosProductSubCategory(subCategoryID string) error {
	if err := r.db.Where("sub_category_id = ?", subCategoryID).Delete(&entity.PosProductSubCategory{}).Error; err != nil {
		return err
	}

	err := r.redis.Del(context.Background(), subCategoryID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posProductSubCategoryRepository) ReadAllPosProductSubCategories(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posProductSubCategories []entity.PosProductSubCategory
	var totalRecords int64

	query := r.db.Model(&entity.PosProductSubCategory{})

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

	query.Find(&posProductSubCategories)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posProductSubCategories,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
