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
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type PosPromotionService interface {
	CreatePosPromotion(ctx context.Context, req *pb.CreatePosPromotionRequest) (*pb.CreatePosPromotionResponse, error)
	ReadPosPromotion(ctx context.Context, req *pb.ReadPosPromotionRequest) (*pb.ReadPosPromotionResponse, error)
	UpdatePosPromotion(ctx context.Context, req *pb.UpdatePosPromotionRequest) (*pb.UpdatePosPromotionResponse, error)
	DeletePosPromotion(ctx context.Context, req *pb.DeletePosPromotionRequest) (*pb.DeletePosPromotionResponse, error)
	ReadAllPosPromotions(ctx context.Context, req *pb.ReadAllPosPromotionsRequest) (*pb.ReadAllPosPromotionsResponse, error)
	ReadPosPromotionByProductId(ctx context.Context, req *pb.ReadPosPromotionByProductIdRequest) (*pb.ReadPosPromotionByProductIdResponse, error)
}

type posPromotionService struct {
	pb.UnimplementedPosPromotionServiceServer
	repo               repository.PosPromotionRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosPromotionService(repo repository.PosPromotionRepository, companyServiceConn *grpc.ClientConn) *posPromotionService {
	return &posPromotionService{
		repo:               repo,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posPromotionService) CreatePosPromotion(ctx context.Context, req *pb.CreatePosPromotionRequest) (*pb.CreatePosPromotionResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to create new promotion")
	}

	req.PosPromotion.PromotionId = uuid.New().String()

	now := timestamppb.New(time.Now())
	req.PosPromotion.CreatedAt = now
	req.PosPromotion.UpdatedAt = now

	startDate, err := time.Parse("2006-01-02", req.PosPromotion.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.PosPromotion.EndDate)
	if err != nil {
		return nil, err
	}

	// Convert pb.PosPromotion to entity.PosPromotion
	entityPromotion := &entity.PosPromotion{
		PromotionID:  uuid.MustParse(req.PosPromotion.PromotionId), // auto
		ProductID:    uuid.MustParse(req.PosPromotion.ProductId),
		StartDate:    startDate,
		EndDate:      endDate,
		Active:       req.PosPromotion.Active,
		DiscountRate: req.PosPromotion.DiscountRate,
		StoreID:      uuid.MustParse(req.PosPromotion.StoreId),
		BranchID:     nil,
		CompanyID:    uuid.MustParse(req.JwtPayload.CompanyId), // auto
		CreatedAt:    time.Now(),
		CreatedBy:    uuid.MustParse(req.JwtPayload.UserId),
		UpdatedAt:    time.Now(),
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	// set Branch ID base in login role
	switch loginRole.PosRole.RoleName {
	case companyRole:
		entityPromotion.BranchID = utils.ParseUUID(req.PosPromotion.BranchId)

		if entityPromotion.BranchID == nil {
			return nil, errors.New("error created promotions, branch id could not be empty")
		}
	case branchRole:
		entityPromotion.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
	}

	err = s.repo.CreatePosPromotion(entityPromotion)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosPromotionResponse{
		PosPromotion: req.PosPromotion,
	}, nil
}

func (s *posPromotionService) ReadPosPromotion(ctx context.Context, req *pb.ReadPosPromotionRequest) (*pb.ReadPosPromotionResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read promotion")
	}

	posPromotion, err := s.repo.ReadPosPromotion(req.PromotionId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posPromotion.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve promotion within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posPromotion.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve promotion within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posPromotion.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve promotion within their store")
		}
	}

	return &pb.ReadPosPromotionResponse{
		PosPromotion: posPromotion,
	}, nil
}

func (s *posPromotionService) ReadPosPromotionByProductId(ctx context.Context, req *pb.ReadPosPromotionByProductIdRequest) (*pb.ReadPosPromotionByProductIdResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read promotion")
	}

	posPromotion, err := s.repo.ReadPosPromotionByProductId(req.ProductId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ReadPosPromotionByProductIdResponse{
				PosPromotion: nil,
			}, nil
		}
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posPromotion.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve promotion within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posPromotion.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve promotion within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posPromotion.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve promotion within their store")
		}
	}

	return &pb.ReadPosPromotionByProductIdResponse{
		PosPromotion: posPromotion,
	}, nil
}

func (s *posPromotionService) UpdatePosPromotion(ctx context.Context, req *pb.UpdatePosPromotionRequest) (*pb.UpdatePosPromotionResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to update promotion")
	}

	posPromotion, err := s.repo.ReadPosPromotion(req.PosPromotion.PromotionId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posPromotion.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update promotion within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posPromotion.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update promotion within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosPromotion.UpdatedAt = now

	startDate, err := time.Parse("2006-01-02", req.PosPromotion.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.PosPromotion.EndDate)
	if err != nil {
		return nil, err
	}

	// Convert pb.PosPromotion to entity.PosPromotion
	entityPromotion := &entity.PosPromotion{
		PromotionID:  uuid.MustParse(posPromotion.PromotionId), // auto
		ProductID:    uuid.MustParse(posPromotion.ProductId),   // auto
		StartDate:    startDate,
		EndDate:      endDate,
		Active:       req.PosPromotion.Active,
		DiscountRate: req.PosPromotion.DiscountRate,
		StoreID:      uuid.MustParse(posPromotion.StoreId),   // auto
		BranchID:     nil,                                    // auto
		CompanyID:    uuid.MustParse(posPromotion.CompanyId), // auto
		CreatedAt:    posPromotion.CreatedAt.AsTime(),        // auto
		CreatedBy:    uuid.MustParse(posPromotion.CreatedBy), // auto
		UpdatedAt:    req.PosPromotion.UpdatedAt.AsTime(),    // auto
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),  // auto
	}

	// set Branch ID base in login role
	switch loginRole.PosRole.RoleName {
	case companyRole:
		entityPromotion.BranchID = utils.ParseUUID(posPromotion.BranchId)

		if entityPromotion.BranchID == nil {
			return nil, errors.New("error update promotions, branch id could not be empty")
		}

	case branchRole:
		entityPromotion.BranchID = utils.ParseUUID(posPromotion.BranchId)
	}

	posPromotion, err = s.repo.UpdatePosPromotion(entityPromotion)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosPromotionResponse{
		PosPromotion: posPromotion,
	}, nil
}

func (s *posPromotionService) DeletePosPromotion(ctx context.Context, req *pb.DeletePosPromotionRequest) (*pb.DeletePosPromotionResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to delete promotion")
	}

	posPromotion, err := s.repo.ReadPosPromotion(req.PromotionId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posPromotion.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete promotion within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posPromotion.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete promotion within their branch")
		}
	}

	err = s.repo.DeletePosPromotion(req.PromotionId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosPromotionResponse{
		Success: true,
	}, nil
}

func (s *posPromotionService) ReadAllPosPromotions(ctx context.Context, req *pb.ReadAllPosPromotionsRequest) (*pb.ReadAllPosPromotionsResponse, error) {
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
		return nil, errors.New("users are not allowed to read all promotion")
	}

	paginationResult, err := s.repo.ReadAllPosPromotions(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posPromotions := paginationResult.Records.([]entity.PosPromotion)
	pbPosPromotions := make([]*pb.PosPromotion, len(posPromotions))

	for i, posPromotion := range posPromotions {
		pbPosPromotions[i] = &pb.PosPromotion{
			PromotionId:  posPromotion.PromotionID.String(),
			ProductId:    posPromotion.ProductID.String(),
			StartDate:    posPromotion.StartDate.Format(time.RFC3339),
			EndDate:      posPromotion.EndDate.Format(time.RFC3339),
			DiscountRate: posPromotion.DiscountRate,
			BranchId:     posPromotion.BranchID.String(),
			CompanyId:    posPromotion.CompanyID.String(),
			CreatedAt:    timestamppb.New(posPromotion.CreatedAt),
			CreatedBy:    posPromotion.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posPromotion.UpdatedAt),
			UpdatedBy:    posPromotion.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosPromotionsResponse{
		PosPromotions: pbPosPromotions,
		Limit:         int32(pagination.Limit),
		Page:          int32(pagination.Page),
		MaxPage:       int32(paginationResult.TotalPages),
		Count:         paginationResult.TotalRecords,
	}, nil
}
