package service

import (
	"context"
	"errors"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-product-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-product-service/utils"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosProductCategoryService interface {
	CreatePosProductCategory(ctx context.Context, req *pb.CreatePosProductCategoryRequest) (*pb.CreatePosProductCategoryResponse, error)
	ReadPosProductCategory(ctx context.Context, req *pb.ReadPosProductCategoryRequest) (*pb.ReadPosProductCategoryResponse, error)
	UpdatePosProductCategory(ctx context.Context, req *pb.UpdatePosProductCategoryRequest) (*pb.UpdatePosProductCategoryResponse, error)
	DeletePosProductCategory(ctx context.Context, req *pb.DeletePosProductCategoryRequest) (*pb.DeletePosProductCategoryResponse, error)
	ReadAllPosProductCategories(ctx context.Context, req *pb.ReadAllPosProductCategoriesRequest) (*pb.ReadAllPosProductCategoriesResponse, error)
}

type posProductCategoryService struct {
	pb.UnimplementedPosProductCategoryServiceServer
	repo repository.PosProductCategoryRepository
}

func NewPosProductCategoryService(repo repository.PosProductCategoryRepository) *posProductCategoryService {
	return &posProductCategoryService{
		repo: repo,
	}
}

func (s *posProductCategoryService) CreatePosProductCategory(ctx context.Context, req *pb.CreatePosProductCategoryRequest) (*pb.CreatePosProductCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to create new product category")
	}

	req.PosProductCategory.CategoryId = uuid.New().String() // Generate a new UUID for the category_id

	now := timestamppb.New(time.Now())
	req.PosProductCategory.CreatedAt = now
	req.PosProductCategory.UpdatedAt = now

	// Convert pb.PosProductCategory to entity.PosProductCategory
	gormCategory := &entity.PosProductCategory{
		CategoryID:   uuid.MustParse(req.PosProductCategory.CategoryId), // auto
		CategoryName: req.PosProductCategory.CategoryName,
		CompanyID:    uuid.MustParse(req.JwtPayload.CompanyId),  // auto
		CreatedAt:    req.PosProductCategory.CreatedAt.AsTime(), // auto
		CreatedBy:    uuid.MustParse(req.JwtPayload.UserId),     // auto
		UpdatedAt:    req.PosProductCategory.UpdatedAt.AsTime(), // auto
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),     // auto
	}

	err = s.repo.CreatePosProductCategory(gormCategory)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosProductCategoryResponse{
		PosProductCategory: req.PosProductCategory,
	}, nil
}

func (s *posProductCategoryService) ReadPosProductCategory(ctx context.Context, req *pb.ReadPosProductCategoryRequest) (*pb.ReadPosProductCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to read product category")
	}

	posProductCategory, err := s.repo.ReadPosProductCategory(req.CategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posProductCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve product category within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posProductCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("branch users can only retrieve product category within their company")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posProductCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("store users can only retrieve product category within their company")
		}
	}

	return &pb.ReadPosProductCategoryResponse{
		PosProductCategory: posProductCategory,
	}, nil
}

func (s *posProductCategoryService) UpdatePosProductCategory(ctx context.Context, req *pb.UpdatePosProductCategoryRequest) (*pb.UpdatePosProductCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to update product category")
	}

	posCategory, err := s.repo.ReadPosProductCategory(req.PosProductCategory.CategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update product category within their company")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosProductCategory.UpdatedAt = now

	// Convert pb.PosProductCategory to entity.PosProductCategory
	gormCategory := &entity.PosProductCategory{
		CategoryID:   uuid.MustParse(posCategory.CategoryId), // auto
		CategoryName: req.PosProductCategory.CategoryName,
		CompanyID:    uuid.MustParse(posCategory.CompanyId),     // auto
		CreatedAt:    posCategory.CreatedAt.AsTime(),            // auto
		CreatedBy:    uuid.MustParse(posCategory.CreatedBy),     // auto
		UpdatedAt:    req.PosProductCategory.UpdatedAt.AsTime(), // auto
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),     // auto
	}

	posProductCategory, err := s.repo.UpdatePosProductCategory(gormCategory)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosProductCategoryResponse{
		PosProductCategory: posProductCategory,
	}, nil
}

func (s *posProductCategoryService) DeletePosProductCategory(ctx context.Context, req *pb.DeletePosProductCategoryRequest) (*pb.DeletePosProductCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to delete product category")
	}

	posCategory, err := s.repo.ReadPosProductCategory(req.CategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete product category within their company")
		}
	}

	err = s.repo.DeletePosProductCategory(req.CategoryId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosProductCategoryResponse{
		Success: true,
	}, nil
}

func (s *posProductCategoryService) ReadAllPosProductCategories(ctx context.Context, req *pb.ReadAllPosProductCategoriesRequest) (*pb.ReadAllPosProductCategoriesResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to read all product category")
	}

	paginationResult, err := s.repo.ReadAllPosProductCategories(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posProductCategories := paginationResult.Records.([]entity.PosProductCategory)
	pbPosProductCategories := make([]*pb.PosProductCategory, len(posProductCategories))

	for i, posProductCategory := range posProductCategories {
		pbPosProductCategories[i] = &pb.PosProductCategory{
			CategoryId:   posProductCategory.CategoryID.String(),
			CategoryName: posProductCategory.CategoryName,
			CompanyId:    posProductCategory.CompanyID.String(),
			CreatedAt:    timestamppb.New(posProductCategory.CreatedAt),
			CreatedBy:    posProductCategory.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posProductCategory.UpdatedAt),
			UpdatedBy:    posProductCategory.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosProductCategoriesResponse{
		PosProductCategories: pbPosProductCategories,
		Limit:                int32(pagination.Limit),
		Page:                 int32(pagination.Page),
		MaxPage:              int32(paginationResult.TotalPages),
		Count:                paginationResult.TotalRecords,
	}, nil
}
