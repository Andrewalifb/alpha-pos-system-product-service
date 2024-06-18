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
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
)

type PosSupplierService interface {
	CreatePosSupplier(ctx context.Context, req *pb.CreatePosSupplierRequest) (*pb.CreatePosSupplierResponse, error)
	ReadPosSupplier(ctx context.Context, req *pb.ReadPosSupplierRequest) (*pb.ReadPosSupplierResponse, error)
	UpdatePosSupplier(ctx context.Context, req *pb.UpdatePosSupplierRequest) (*pb.UpdatePosSupplierResponse, error)
	DeletePosSupplier(ctx context.Context, req *pb.DeletePosSupplierRequest) (*pb.DeletePosSupplierResponse, error)
	ReadAllPosSuppliers(ctx context.Context, req *pb.ReadAllPosSuppliersRequest) (*pb.ReadAllPosSuppliersResponse, error)
}

type posSupplierService struct {
	pb.UnimplementedPosSupplierServiceServer
	repo repository.PosSupplierRepository
}

func NewPosSupplierService(repo repository.PosSupplierRepository) *posSupplierService {
	return &posSupplierService{
		repo: repo,
	}
}

func (s *posSupplierService) CreatePosSupplier(ctx context.Context, req *pb.CreatePosSupplierRequest) (*pb.CreatePosSupplierResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to create new supplier")
	}

	now := timestamppb.New(time.Now())
	req.PosSupplier.SupplierId = uuid.New().String()
	req.PosSupplier.CreatedAt = now
	req.PosSupplier.UpdatedAt = now

	entitySupplier := &entity.PosSupplier{
		SupplierID:   uuid.MustParse(req.PosSupplier.SupplierId), // auto
		SupplierName: req.PosSupplier.SupplierName,
		BranchID:     nil,                                      // auto
		CompanyID:    uuid.MustParse(req.JwtPayload.CompanyId), // auto
		CreatedAt:    req.PosSupplier.CreatedAt.AsTime(),       // auto
		CreatedBy:    uuid.MustParse(req.JwtPayload.UserId),    // auto
		UpdatedAt:    req.PosSupplier.UpdatedAt.AsTime(),       // auto
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),    // auto
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	switch loginRole.Data.RoleName {
	case companyRole:
		entitySupplier.BranchID = utils.ParseUUID(req.PosSupplier.BranchId)

		if entitySupplier.BranchID == nil {
			return nil, errors.New("error created supplier, branch id could not be empty")
		}
	case branchRole:
		entitySupplier.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
	}

	// Check if Branch ID is correct
	_, err = utils.GetPosStoreBranch(entitySupplier.BranchID.String(), req.JwtToken)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreatePosSupplier(entitySupplier)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosSupplierResponse{
		PosSupplier: req.PosSupplier,
	}, nil
}

func (s *posSupplierService) ReadPosSupplier(ctx context.Context, req *pb.ReadPosSupplierRequest) (*pb.ReadPosSupplierResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to read promotion")
	}

	posSupplier, err := s.repo.ReadPosSupplier(req.SupplierId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posSupplier.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve supplier within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSupplier.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve supplier within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posSupplier.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("store users can only retrieve supplier within their branch")
		}
	}

	return &pb.ReadPosSupplierResponse{
		PosSupplier: posSupplier,
	}, nil
}

func (s *posSupplierService) UpdatePosSupplier(ctx context.Context, req *pb.UpdatePosSupplierRequest) (*pb.UpdatePosSupplierResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to update promotion")
	}

	// Get the supplier to be updated
	posSupplierData, err := s.repo.ReadPosSupplier(req.PosSupplier.SupplierId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posSupplierData.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update promotion within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSupplierData.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update promotion within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosSupplier.UpdatedAt = now
	entitySupplier := &entity.PosSupplier{
		SupplierID:   uuid.MustParse(posSupplierData.SupplierId), // auto
		SupplierName: req.PosSupplier.SupplierName,
		ContactName:  req.PosSupplier.ContactName,
		ContactEmail: req.PosSupplier.ContactEmail,
		ContactPhone: req.PosSupplier.ContactPhone,
		BranchID:     nil,                                       // auto
		CompanyID:    uuid.MustParse(posSupplierData.CompanyId), // auto
		CreatedAt:    posSupplierData.CreatedAt.AsTime(),        // auto
		CreatedBy:    uuid.MustParse(posSupplierData.CreatedBy), // auto
		UpdatedAt:    req.PosSupplier.UpdatedAt.AsTime(),        // auto
		UpdatedBy:    uuid.MustParse(req.JwtPayload.UserId),     // auto
	}

	// Set Branch ID from database
	entitySupplier.BranchID = utils.ParseUUID(posSupplierData.BranchId)

	branchData, err := utils.GetPosStoreBranch(entitySupplier.BranchID.String(), req.JwtToken)
	if err != nil {
		return nil, err
	}
	// set branch id base on data
	entitySupplier.BranchID = utils.ParseUUID(branchData.Data.BranchID)

	posSupplier, err := s.repo.UpdatePosSupplier(entitySupplier)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosSupplierResponse{
		PosSupplier: posSupplier,
	}, nil
}

func (s *posSupplierService) DeletePosSupplier(ctx context.Context, req *pb.DeletePosSupplierRequest) (*pb.DeletePosSupplierResponse, error) {

	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to update promotion")
	}

	// Get the supplier to be updated
	posSupplierData, err := s.repo.ReadPosSupplier(req.SupplierId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posSupplierData.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update promotion within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSupplierData.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update promotion within their branch")
		}
	}

	err = s.repo.DeletePosSupplier(req.SupplierId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosSupplierResponse{
		Success: true,
	}, nil
}

func (s *posSupplierService) ReadAllPosSuppliers(ctx context.Context, req *pb.ReadAllPosSuppliersRequest) (*pb.ReadAllPosSuppliersResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}

	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to read all promotion")
	}

	paginationResult, err := s.repo.ReadAllPosSuppliers(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posSuppliers := paginationResult.Records.([]entity.PosSupplier)
	pbPosSuppliers := make([]*pb.PosSupplier, len(posSuppliers))

	for i, posSupplier := range posSuppliers {
		pbPosSuppliers[i] = &pb.PosSupplier{
			SupplierId:   posSupplier.SupplierID.String(),
			SupplierName: posSupplier.SupplierName,
			BranchId:     posSupplier.BranchID.String(),
			CompanyId:    posSupplier.CompanyID.String(),
			CreatedAt:    timestamppb.New(posSupplier.CreatedAt),
			CreatedBy:    posSupplier.CreatedBy.String(),
			UpdatedAt:    timestamppb.New(posSupplier.UpdatedAt),
			UpdatedBy:    posSupplier.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosSuppliersResponse{
		PosSuppliers: pbPosSuppliers,
		Limit:        int32(pagination.Limit),
		Page:         int32(pagination.Page),
		MaxPage:      int32(paginationResult.TotalPages),
		Count:        paginationResult.TotalRecords,
	}, nil
}
