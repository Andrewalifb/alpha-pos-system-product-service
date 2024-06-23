package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-product-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-product-service/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosProductService interface {
	CreatePosProduct(ctx context.Context, req *pb.CreatePosProductRequest) (*pb.CreatePosProductResponse, error)
	ReadPosProduct(ctx context.Context, req *pb.ReadPosProductRequest) (*pb.ReadPosProductResponse, error)
	ReadPosProductByBarcode(ctx context.Context, req *pb.ReadPosProductByBarcodeRequest) (*pb.ReadPosProductByBarcodeResponse, error)
	UpdatePosProduct(ctx context.Context, req *pb.UpdatePosProductRequest) (*pb.UpdatePosProductResponse, error)
	DeletePosProduct(ctx context.Context, req *pb.DeletePosProductRequest) (*pb.DeletePosProductResponse, error)
	ReadAllPosProducts(ctx context.Context, req *pb.ReadAllPosProductsRequest) (*pb.ReadAllPosProductsResponse, error)
}

type posProductService struct {
	pb.UnimplementedPosProductServiceServer
	productRepo        repository.PosProductRepository
	supplierRepo       repository.PosSupplierRepository
	categoryRepo       repository.PosProductCategoryRepository
	subCategory        repository.PosProductSubCategoryRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosProductService(productRepo repository.PosProductRepository, supplierRepo repository.PosSupplierRepository, categoryRepo repository.PosProductCategoryRepository, subCategory repository.PosProductSubCategoryRepository, companyServiceConn *grpc.ClientConn) *posProductService {
	return &posProductService{
		productRepo:        productRepo,
		supplierRepo:       supplierRepo,
		categoryRepo:       categoryRepo,
		subCategory:        subCategory,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posProductService) CreatePosProduct(ctx context.Context, req *pb.CreatePosProductRequest) (*pb.CreatePosProductResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to create new product")
	}

	req.PosProduct.ProductId = uuid.New().String() // Generate a new UUID for the product_id

	now := timestamppb.New(time.Now())
	req.PosProduct.CreatedAt = now
	req.PosProduct.UpdatedAt = now

	// Convert pb.PosProduct to entity.PosProduct
	gormProduct := &entity.PosProduct{
		ProductID:          uuid.MustParse(req.PosProduct.ProductId),
		ProductBarcodeID:   strings.ToLower(req.PosProduct.ProductBarcodeId),
		ProductName:        req.PosProduct.ProductName,
		Price:              req.PosProduct.Price,
		CostPrice:          req.PosProduct.CostPrice,
		CategoryID:         uuid.MustParse(req.PosProduct.CategoryId),
		SubCategoryID:      uuid.MustParse(req.PosProduct.SubCategoryId),
		StockQuantity:      0, // auto default 0
		ReorderLevel:       int(req.PosProduct.ReorderLevel),
		SupplierID:         uuid.MustParse(req.PosProduct.SupplierId),
		ProductDescription: req.PosProduct.ProductDescription,
		Active:             req.PosProduct.Active,
		StoreID:            uuid.MustParse(req.PosProduct.StoreId),
		BranchID:           nil,
		CompanyID:          uuid.MustParse(req.JwtPayload.CompanyId), // auto
		CreatedAt:          req.PosProduct.CreatedAt.AsTime(),        // auto
		CreatedBy:          uuid.MustParse(req.JwtPayload.UserId),    // auto
		UpdatedAt:          req.PosProduct.UpdatedAt.AsTime(),        // auto
		UpdatedBy:          uuid.MustParse(req.JwtPayload.UserId),    // auto
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	// set Branch ID base in login role
	switch loginRole.PosRole.RoleName {
	case companyRole:
		gormProduct.BranchID = utils.ParseUUID(req.PosProduct.BranchId)

		if gormProduct.BranchID == nil {
			return nil, errors.New("error created product, branch id could not be empty")
		}
	case branchRole:
		gormProduct.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
	}

	// Check if category ID is correct
	_, err = s.categoryRepo.ReadPosProductCategory(gormProduct.CategoryID.String())
	if err != nil {
		return nil, err
	}

	// Check if sub category ID is correct
	_, err = s.subCategory.ReadPosProductSubCategory(gormProduct.SubCategoryID.String())
	if err != nil {
		return nil, err
	}

	// Check if supplier ID is correct
	_, err = s.supplierRepo.ReadPosSupplier(gormProduct.SupplierID.String())
	if err != nil {
		return nil, err
	}

	err = s.productRepo.CreatePosProduct(gormProduct)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosProductResponse{
		PosProduct: req.PosProduct,
	}, nil
}

func (s *posProductService) ReadPosProduct(ctx context.Context, req *pb.ReadPosProductRequest) (*pb.ReadPosProductResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read product")
	}

	posProduct, err := s.productRepo.ReadPosProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posProduct.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve product within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posProduct.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve product within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posProduct.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve product within their store")
		}
	}

	return &pb.ReadPosProductResponse{
		PosProduct: posProduct,
	}, nil
}

func (s *posProductService) ReadPosProductByBarcode(ctx context.Context, req *pb.ReadPosProductByBarcodeRequest) (*pb.ReadPosProductByBarcodeResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read product")
	}

	posProduct, err := s.productRepo.ReadPosProductBarcode(strings.ToLower(req.ProductBarcodeId))
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posProduct.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve product within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posProduct.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve product within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posProduct.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("store users can only retrieve product within their store")
		}
	}

	return &pb.ReadPosProductByBarcodeResponse{
		PosProduct: posProduct,
	}, nil
}

func (s *posProductService) UpdatePosProduct(ctx context.Context, req *pb.UpdatePosProductRequest) (*pb.UpdatePosProductResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to update product")
	}

	// Get the product to be updated
	posProduct, err := s.productRepo.ReadPosProduct(req.PosProduct.ProductId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posProduct.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update product within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posProduct.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update product within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosProduct.UpdatedAt = now

	// Convert pb.PosProduct to entity.PosProduct
	gormProduct := &entity.PosProduct{
		ProductID:          uuid.MustParse(posProduct.ProductId),         // auto
		ProductBarcodeID:   strings.ToLower(posProduct.ProductBarcodeId), // auto
		ProductName:        req.PosProduct.ProductName,
		Price:              req.PosProduct.Price,
		CostPrice:          req.PosProduct.CostPrice,
		CategoryID:         uuid.MustParse(posProduct.CategoryId),    // auto
		SubCategoryID:      uuid.MustParse(posProduct.SubCategoryId), // auto
		StockQuantity:      int(posProduct.StockQuantity),            // auto
		ReorderLevel:       int(req.PosProduct.ReorderLevel),
		SupplierID:         uuid.MustParse(req.PosProduct.SupplierId),
		ProductDescription: req.PosProduct.ProductDescription,
		Active:             req.PosProduct.Active,
		StoreID:            uuid.MustParse(req.PosProduct.StoreId),
		BranchID:           nil,                                   // auto
		CompanyID:          uuid.MustParse(posProduct.CompanyId),  // auto
		CreatedAt:          posProduct.CreatedAt.AsTime(),         // auto
		CreatedBy:          uuid.MustParse(posProduct.CreatedBy),  // auto
		UpdatedAt:          req.PosProduct.UpdatedAt.AsTime(),     // auto
		UpdatedBy:          uuid.MustParse(req.JwtPayload.UserId), // auto
	}

	// Set Branch ID From Databse
	gormProduct.BranchID = utils.ParseUUID(posProduct.BranchId)
	err = s.productRepo.UpdatePosProduct(gormProduct)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosProductResponse{
		PosProduct: req.PosProduct,
	}, nil
}

func (s *posProductService) DeletePosProduct(ctx context.Context, req *pb.DeletePosProductRequest) (*pb.DeletePosProductResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to update product")
	}

	// Get the supplier to be updated
	posProduct, err := s.productRepo.ReadPosProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posProduct.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete product within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posProduct.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete product within their branch")
		}
	}

	err = s.productRepo.DeletePosProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosProductResponse{
		Success: true,
	}, nil
}

func (s *posProductService) ReadAllPosProducts(ctx context.Context, req *pb.ReadAllPosProductsRequest) (*pb.ReadAllPosProductsResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}

	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read all product")
	}

	paginationResult, err := s.productRepo.ReadAllPosProducts(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posProducts := paginationResult.Records.([]entity.PosProduct)
	pbPosProducts := make([]*pb.PosProduct, len(posProducts))

	for i, posProduct := range posProducts {
		pbPosProducts[i] = &pb.PosProduct{
			ProductId:          posProduct.ProductID.String(),
			ProductBarcodeId:   posProduct.ProductBarcodeID,
			ProductName:        posProduct.ProductName,
			Price:              posProduct.Price,
			CostPrice:          posProduct.CostPrice,
			CategoryId:         posProduct.CategoryID.String(),
			SubCategoryId:      posProduct.SubCategoryID.String(),
			StockQuantity:      int32(posProduct.StockQuantity),
			ReorderLevel:       int32(posProduct.ReorderLevel),
			SupplierId:         posProduct.SupplierID.String(),
			ProductDescription: posProduct.ProductDescription,
			Active:             posProduct.Active,
			StoreId:            posProduct.StoreID.String(),
			BranchId:           posProduct.BranchID.String(),
			CompanyId:          posProduct.CompanyID.String(),
			CreatedAt:          timestamppb.New(posProduct.CreatedAt),
			CreatedBy:          posProduct.CreatedBy.String(),
			UpdatedAt:          timestamppb.New(posProduct.UpdatedAt),
			UpdatedBy:          posProduct.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosProductsResponse{
		PosProducts: pbPosProducts,
		Limit:       int32(pagination.Limit),
		Page:        int32(pagination.Page),
		MaxPage:     int32(paginationResult.TotalPages),
		Count:       paginationResult.TotalRecords,
	}, nil
}
