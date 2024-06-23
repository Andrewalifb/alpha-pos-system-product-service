package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-product-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-product-service/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosProductSubCategoryService interface {
	CreatePosProductSubCategory(ctx context.Context, req *pb.CreatePosProductSubCategoryRequest) (*pb.CreatePosProductSubCategoryResponse, error)
	ReadPosProductSubCategory(ctx context.Context, req *pb.ReadPosProductSubCategoryRequest) (*pb.ReadPosProductSubCategoryResponse, error)
	UpdatePosProductSubCategory(ctx context.Context, req *pb.UpdatePosProductSubCategoryRequest) (*pb.UpdatePosProductSubCategoryResponse, error)
	DeletePosProductSubCategory(ctx context.Context, req *pb.DeletePosProductSubCategoryRequest) (*pb.DeletePosProductSubCategoryResponse, error)
	ReadAllPosProductSubCategories(ctx context.Context, req *pb.ReadAllPosProductSubCategoriesRequest) (*pb.ReadAllPosProductSubCategoriesResponse, error)
}

type posProductSubCategoryService struct {
	pb.UnimplementedPosProductSubCategoryServiceServer
	subCategoryRepo    repository.PosProductSubCategoryRepository
	categoryRepo       repository.PosProductCategoryRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosProductSubCategoryService(subCategoryRepo repository.PosProductSubCategoryRepository, categoryRepo repository.PosProductCategoryRepository, companyServiceConn *grpc.ClientConn) *posProductSubCategoryService {
	return &posProductSubCategoryService{
		subCategoryRepo:    subCategoryRepo,
		categoryRepo:       categoryRepo,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posProductSubCategoryService) CreatePosProductSubCategory(ctx context.Context, req *pb.CreatePosProductSubCategoryRequest) (*pb.CreatePosProductSubCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	if req.PosProductSubCategory.CategoryId == "" {
		return nil, errors.New("error created sub category, category id could not be empty")
	}

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to create new product sub category")
	}

	// Check if category is avaliable for sub category
	isCategoryExist, err := s.categoryRepo.IsCategoryExist(req.PosProductSubCategory.CategoryId)
	if err != nil {
		// Log the error and return a user-friendly message
		log.Printf("Error checking if category exists: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !isCategoryExist {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("category id %s not avaliable", req.PosProductSubCategory.CategoryId))
	}

	req.PosProductSubCategory.SubCategoryId = uuid.New().String()

	now := timestamppb.New(time.Now())
	req.PosProductSubCategory.CreatedAt = now
	req.PosProductSubCategory.UpdatedAt = now

	entitySubCategory := &entity.PosProductSubCategory{
		SubCategoryID:   uuid.MustParse(req.PosProductSubCategory.SubCategoryId), // auto
		SubCategoryName: req.PosProductSubCategory.SubCategoryName,
		CategoryID:      uuid.MustParse(req.PosProductSubCategory.CategoryId),
		CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId),     // auto
		CreatedAt:       req.PosProductSubCategory.CreatedAt.AsTime(), // auto
		CreatedBy:       uuid.MustParse(req.JwtPayload.UserId),        // auto
		UpdatedAt:       req.PosProductSubCategory.UpdatedAt.AsTime(), // auto
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),        // auto
	}

	err = s.subCategoryRepo.CreatePosProductSubCategory(entitySubCategory)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosProductSubCategoryResponse{
		PosProductSubCategory: req.PosProductSubCategory,
	}, nil
}

func (s *posProductSubCategoryService) ReadPosProductSubCategory(ctx context.Context, req *pb.ReadPosProductSubCategoryRequest) (*pb.ReadPosProductSubCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read product sub category")
	}

	posSubCategory, err := s.subCategoryRepo.ReadPosProductSubCategory(req.SubCategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posSubCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve product sub category within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posSubCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("branch users can only retrieve product sub category within their company")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posSubCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("store users can only retrieve product sub category within their company")
		}
	}

	return &pb.ReadPosProductSubCategoryResponse{
		PosProductSubCategory: posSubCategory,
	}, nil
}

func (s *posProductSubCategoryService) UpdatePosProductSubCategory(ctx context.Context, req *pb.UpdatePosProductSubCategoryRequest) (*pb.UpdatePosProductSubCategoryResponse, error) {

	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to update product sub category")
	}

	posSubCategory, err := s.subCategoryRepo.ReadPosProductSubCategory(req.PosProductSubCategory.SubCategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posSubCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update product sub category within their company")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosProductSubCategory.UpdatedAt = now
	entitySubCategory := &entity.PosProductSubCategory{
		SubCategoryID:   uuid.MustParse(posSubCategory.SubCategoryId), // auto
		SubCategoryName: req.PosProductSubCategory.SubCategoryName,
		CategoryID:      uuid.MustParse(posSubCategory.CategoryId),    // auto
		CompanyID:       uuid.MustParse(posSubCategory.CompanyId),     // autp
		CreatedAt:       posSubCategory.CreatedAt.AsTime(),            // auto
		CreatedBy:       uuid.MustParse(posSubCategory.CreatedBy),     // autp
		UpdatedAt:       req.PosProductSubCategory.UpdatedAt.AsTime(), // auto
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),        // auto
	}

	posProductSubCategory, err := s.subCategoryRepo.UpdatePosProductSubCategory(entitySubCategory)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosProductSubCategoryResponse{
		PosProductSubCategory: posProductSubCategory,
	}, nil
}

func (s *posProductSubCategoryService) DeletePosProductSubCategory(ctx context.Context, req *pb.DeletePosProductSubCategoryRequest) (*pb.DeletePosProductSubCategoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to delete product sub category")
	}

	posSubCategory, err := s.subCategoryRepo.ReadPosProductSubCategory(req.SubCategoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posSubCategory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete product sub category within their company")
		}
	}

	err = s.subCategoryRepo.DeletePosProductSubCategory(req.SubCategoryId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosProductSubCategoryResponse{
		Success: true,
	}, nil
}

func (s *posProductSubCategoryService) ReadAllPosProductSubCategories(ctx context.Context, req *pb.ReadAllPosProductSubCategoriesRequest) (*pb.ReadAllPosProductSubCategoriesResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}

	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read all product category")
	}

	paginationResult, err := s.subCategoryRepo.ReadAllPosProductSubCategories(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posProductSubCategories := paginationResult.Records.([]entity.PosProductSubCategory)
	pbPosProductSubCategories := make([]*pb.PosProductSubCategory, len(posProductSubCategories))

	for i, posProductSubCategory := range posProductSubCategories {
		pbPosProductSubCategories[i] = &pb.PosProductSubCategory{
			SubCategoryId:   posProductSubCategory.SubCategoryID.String(),
			SubCategoryName: posProductSubCategory.SubCategoryName,
			CategoryId:      posProductSubCategory.CategoryID.String(),
			CompanyId:       posProductSubCategory.CompanyID.String(),
			CreatedAt:       timestamppb.New(posProductSubCategory.CreatedAt),
			CreatedBy:       posProductSubCategory.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posProductSubCategory.UpdatedAt),
			UpdatedBy:       posProductSubCategory.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosProductSubCategoriesResponse{
		PosProductSubCategories: pbPosProductSubCategories,
		Limit:                   int32(pagination.Limit),
		Page:                    int32(pagination.Page),
		MaxPage:                 int32(paginationResult.TotalPages),
		Count:                   paginationResult.TotalRecords,
	}, nil
}
