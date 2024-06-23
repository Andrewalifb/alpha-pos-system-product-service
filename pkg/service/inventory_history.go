package service

import (
	"context"
	"errors"
	"math"
	"os"
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

type PosInventoryHistoryService interface {
	CreatePosInventoryHistory(ctx context.Context, req *pb.CreatePosInventoryHistoryRequest) (*pb.CreatePosInventoryHistoryResponse, error)
	ReadPosInventoryHistory(ctx context.Context, req *pb.ReadPosInventoryHistoryRequest) (*pb.ReadPosInventoryHistoryResponse, error)
	UpdatePosInventoryHistory(ctx context.Context, req *pb.UpdatePosInventoryHistoryRequest) (*pb.UpdatePosInventoryHistoryResponse, error)
	DeletePosInventoryHistory(ctx context.Context, req *pb.DeletePosInventoryHistoryRequest) (*pb.DeletePosInventoryHistoryResponse, error)
	ReadAllPosInventoryHistories(ctx context.Context, req *pb.ReadAllPosInventoryHistoriesRequest) (*pb.ReadAllPosInventoryHistoriesResponse, error)
}

type posInventoryHistoryService struct {
	pb.UnimplementedPosInventoryHistoryServiceServer
	repoInventory      repository.PosInventoryHistoryRepository
	repoProduct        repository.PosProductRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosInventoryHistoryService(repoInventory repository.PosInventoryHistoryRepository, repoProduct repository.PosProductRepository, companyServiceConn *grpc.ClientConn) *posInventoryHistoryService {
	return &posInventoryHistoryService{
		repoInventory:      repoInventory,
		repoProduct:        repoProduct,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posInventoryHistoryService) CreatePosInventoryHistory(ctx context.Context, req *pb.CreatePosInventoryHistoryRequest) (*pb.CreatePosInventoryHistoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to create new inventory history")
	}

	req.PosInventoryHistory.InventoryId = uuid.New().String() // Generate a new UUID for the inventory_id

	now := timestamppb.New(time.Now())
	req.PosInventoryHistory.CreatedAt = now
	req.PosInventoryHistory.UpdatedAt = now
	req.PosInventoryHistory.Date = now

	// Convert pb.PosInventoryHistory to entity.PosInventoryHistory
	gormInventoryHistory := &entity.PosInventoryHistory{
		InventoryID: uuid.MustParse(req.PosInventoryHistory.InventoryId), // auto
		ProductID:   uuid.MustParse(req.PosInventoryHistory.ProductId),
		StoreID:     nil,
		Date:        req.PosInventoryHistory.Date.AsTime(), // auto
		Quantity:    int(req.PosInventoryHistory.Quantity),
		BranchID:    nil,
		CompanyID:   uuid.MustParse(req.JwtPayload.CompanyId),   // auto
		CreatedAt:   req.PosInventoryHistory.CreatedAt.AsTime(), // auto
		CreatedBy:   uuid.MustParse(req.JwtPayload.UserId),      // auto
		UpdatedAt:   req.PosInventoryHistory.UpdatedAt.AsTime(), // auto
		UpdatedBy:   uuid.MustParse(req.JwtPayload.UserId),      // auto
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	// set Branch ID Store ID base in login role
	switch loginRole.PosRole.RoleName {
	case companyRole:
		gormInventoryHistory.BranchID = utils.ParseUUID(req.PosInventoryHistory.BranchId)

		if gormInventoryHistory.BranchID == nil {
			return nil, errors.New("error created inventory history, branch id could not be empty")
		}

		if gormInventoryHistory.StoreID == nil {
			return nil, errors.New("error created inventory history, store id could not be empty")
		}

	case branchRole:
		gormInventoryHistory.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)

		if gormInventoryHistory.StoreID == nil {
			return nil, errors.New("error created inventory history, store id could not be empty")
		}
	case storeRole:
		gormInventoryHistory.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
		gormInventoryHistory.StoreID = utils.ParseUUID(req.JwtPayload.StoreId)
	}

	err = s.repoInventory.CreatePosInventoryHistory(gormInventoryHistory)
	if err != nil {
		return nil, err
	}

	// Get data product
	getDataProduct, err := s.repoProduct.ReadPosProduct(gormInventoryHistory.ProductID.String())
	if err != nil {
		return nil, err
	}

	// check if decrease the quantity
	if gormInventoryHistory.Quantity < 0 {
		countQuantity := getDataProduct.StockQuantity + int32(gormInventoryHistory.Quantity)
		if countQuantity < 0 {
			return nil, errors.New("error cant decrease stock quantity, the request quantity is bigger than current avaliable stock quantity")
		}
	}

	getDataProduct.UpdatedAt = now
	// Update the products quantity
	updateDataProduct := &entity.PosProduct{
		ProductID:          uuid.MustParse(getDataProduct.ProductId),                          // auto
		ProductBarcodeID:   getDataProduct.ProductBarcodeId,                                   // auto
		ProductName:        getDataProduct.ProductName,                                        // auto
		Price:              getDataProduct.Price,                                              // auto
		CostPrice:          getDataProduct.CostPrice,                                          // auto
		CategoryID:         uuid.MustParse(getDataProduct.CategoryId),                         // auto
		SubCategoryID:      uuid.MustParse(getDataProduct.SubCategoryId),                      // auto
		StockQuantity:      int(getDataProduct.StockQuantity) + gormInventoryHistory.Quantity, // auto
		ReorderLevel:       int(getDataProduct.ReorderLevel),                                  // auto
		SupplierID:         uuid.MustParse(getDataProduct.SupplierId),                         // auto
		ProductDescription: getDataProduct.ProductDescription,                                 // auto
		Active:             getDataProduct.Active,                                             // auto
		StoreID:            uuid.MustParse(getDataProduct.StoreId),                            // auto
		BranchID:           nil,                                                               // auto
		CompanyID:          uuid.MustParse(getDataProduct.CompanyId),                          // auto
		CreatedAt:          getDataProduct.CreatedAt.AsTime(),                                 // auto
		CreatedBy:          uuid.MustParse(getDataProduct.CreatedBy),                          // auto
		UpdatedAt:          getDataProduct.UpdatedAt.AsTime(),                                 // auto
		UpdatedBy:          uuid.MustParse(req.JwtPayload.UserId),                             // auto
	}

	updateDataProduct.BranchID = utils.ParseUUID(getDataProduct.BranchId)

	err = s.repoProduct.UpdatePosProduct(updateDataProduct)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosInventoryHistoryResponse{
		PosInventoryHistory: req.PosInventoryHistory,
	}, nil
}

func (s *posInventoryHistoryService) ReadPosInventoryHistory(ctx context.Context, req *pb.ReadPosInventoryHistoryRequest) (*pb.ReadPosInventoryHistoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read inventory history")
	}

	posInventoryHistory, err := s.repoInventory.ReadPosInventoryHistory(req.InventoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posInventoryHistory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve inventory history within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posInventoryHistory.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve inventory history within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posInventoryHistory.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve inventory history within their store")
		}
	}

	return &pb.ReadPosInventoryHistoryResponse{
		PosInventoryHistory: posInventoryHistory,
	}, nil
}

func (s *posInventoryHistoryService) UpdatePosInventoryHistory(ctx context.Context, req *pb.UpdatePosInventoryHistoryRequest) (*pb.UpdatePosInventoryHistoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read inventory history")
	}

	// Get the inventory to be updated
	posInventory, err := s.repoInventory.ReadPosInventoryHistory(req.PosInventoryHistory.InventoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posInventory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update inventory history within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posInventory.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update inventory history within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosInventoryHistory.UpdatedAt = now

	// Convert pb.PosInventoryHistory to entity.PosInventoryHistory
	gormInventoryHistory := &entity.PosInventoryHistory{
		InventoryID: uuid.MustParse(posInventory.InventoryId), // auto
		ProductID:   uuid.MustParse(posInventory.ProductId),   // auto
		StoreID:     nil,
		Date:        posInventory.Date.AsTime(), // auto
		Quantity:    int(req.PosInventoryHistory.Quantity),
		BranchID:    nil,                                        // auto
		CompanyID:   uuid.MustParse(posInventory.CompanyId),     // auto
		CreatedAt:   posInventory.CreatedAt.AsTime(),            // auto
		CreatedBy:   uuid.MustParse(posInventory.CreatedBy),     // auto
		UpdatedAt:   req.PosInventoryHistory.UpdatedAt.AsTime(), // auto
		UpdatedBy:   uuid.MustParse(req.JwtPayload.UserId),      // auto
	}

	gormInventoryHistory.BranchID = utils.ParseUUID(posInventory.BranchId)
	gormInventoryHistory.StoreID = utils.ParseUUID(posInventory.StoreId)

	err = s.repoInventory.UpdatePosInventoryHistory(gormInventoryHistory)
	if err != nil {
		return nil, err
	}

	// Get data product
	getDataProduct, err := s.repoProduct.ReadPosProduct(gormInventoryHistory.ProductID.String())
	if err != nil {
		return nil, err
	}

	getDataProduct.UpdatedAt = now
	// Update the products quantity
	updateDataProduct := &entity.PosProduct{
		ProductID:          uuid.MustParse(getDataProduct.ProductId),                          // auto
		ProductBarcodeID:   getDataProduct.ProductBarcodeId,                                   // auto
		ProductName:        getDataProduct.ProductName,                                        // auto
		Price:              getDataProduct.Price,                                              // auto
		CostPrice:          getDataProduct.CostPrice,                                          // auto
		CategoryID:         uuid.MustParse(getDataProduct.CategoryId),                         // auto
		SubCategoryID:      uuid.MustParse(getDataProduct.SubCategoryId),                      // auto
		StockQuantity:      int(getDataProduct.StockQuantity) + gormInventoryHistory.Quantity, // auto
		ReorderLevel:       int(getDataProduct.ReorderLevel),                                  // auto
		SupplierID:         uuid.MustParse(getDataProduct.SupplierId),                         // auto
		ProductDescription: getDataProduct.ProductDescription,                                 // auto
		Active:             getDataProduct.Active,                                             // auto
		StoreID:            uuid.MustParse(getDataProduct.StoreId),                            // auto
		BranchID:           nil,                                                               // auto
		CompanyID:          uuid.MustParse(getDataProduct.CompanyId),                          // auto
		CreatedAt:          getDataProduct.CreatedAt.AsTime(),                                 // auto
		CreatedBy:          uuid.MustParse(getDataProduct.CreatedBy),                          // auto
		UpdatedAt:          getDataProduct.UpdatedAt.AsTime(),                                 // auto
		UpdatedBy:          uuid.MustParse(req.JwtPayload.UserId),                             // auto
	}
	// set Branch ID base in login role
	switch loginRole.PosRole.RoleName {
	case companyRole:
		// fmt.Println("Product DAta :", getDataProduct.BranchId)
		updateDataProduct.BranchID = utils.ParseUUID(getDataProduct.BranchId)

		if updateDataProduct.BranchID == nil {
			return nil, errors.New("error update inventory history, branch id could not be empty")
		}
	case branchRole:
		updateDataProduct.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
	}

	// if update req is increase
	if req.PosInventoryHistory.Quantity >= 0 && getDataProduct.StockQuantity > req.PosInventoryHistory.Quantity {
		// Check if current inventory quantity is +
		if posInventory.Quantity > 0 {
			undoCurrentQuantity := getDataProduct.StockQuantity - posInventory.Quantity
			repairedQuantity := undoCurrentQuantity + req.PosInventoryHistory.Quantity
			updateDataProduct.StockQuantity = int(repairedQuantity)
			// Check if current inventory quantity is -
		} else if posInventory.Quantity < 0 {
			undoCurrentQuantity := getDataProduct.StockQuantity + int32(math.Abs(float64(posInventory.Quantity)))
			repairedQuantity := undoCurrentQuantity + req.PosInventoryHistory.Quantity
			updateDataProduct.StockQuantity = int(repairedQuantity)
		}

	}

	// if update req is decrease
	countCurrentQuantity := getDataProduct.StockQuantity + req.PosInventoryHistory.Quantity
	if req.PosInventoryHistory.Quantity < 0 && countCurrentQuantity > 0 {
		// Check if current inventory quantity is +
		if posInventory.Quantity > 0 {
			undoCurrentQuantity := getDataProduct.StockQuantity - posInventory.Quantity
			repairedQuantity := undoCurrentQuantity + req.PosInventoryHistory.Quantity
			updateDataProduct.StockQuantity = int(repairedQuantity)
			// Check if current inventory quantity is -
		} else if posInventory.Quantity < 0 {
			undoCurrentQuantity := getDataProduct.StockQuantity + int32(math.Abs(float64(posInventory.Quantity)))
			repairedQuantity := undoCurrentQuantity + req.PosInventoryHistory.Quantity
			updateDataProduct.StockQuantity = int(repairedQuantity)
		}
	}

	err = s.repoProduct.UpdatePosProduct(updateDataProduct)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosInventoryHistoryResponse{
		PosInventoryHistory: req.PosInventoryHistory,
	}, nil
}

func (s *posInventoryHistoryService) DeletePosInventoryHistory(ctx context.Context, req *pb.DeletePosInventoryHistoryRequest) (*pb.DeletePosInventoryHistoryResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read inventory history")
	}

	// Get the inventory to be deleted
	posInventory, err := s.repoInventory.ReadPosInventoryHistory(req.InventoryId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posInventory.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete inventory history within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posInventory.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete inventory history within their branch")
		}
	}

	// Get data product
	getDataProduct, err := s.repoProduct.ReadPosProduct(posInventory.ProductId)
	if err != nil {
		return nil, err
	}

	var setQuantity int
	// If current inventory history quantity is decreased the product quantity
	if posInventory.Quantity >= 0 {
		repairedQuantity := getDataProduct.StockQuantity - posInventory.Quantity
		if repairedQuantity >= 0 {
			setQuantity = int(repairedQuantity)
		} else if repairedQuantity < 0 {
			return nil, errors.New("cant delete this inventory history, because the current amount of product quantity not sufficient to decrease")
		}
	}
	// if current inventory history quantity is increased the product quantity
	if posInventory.Quantity < 0 {
		repairedQuantity := getDataProduct.StockQuantity + int32(math.Abs(float64(posInventory.Quantity)))
		setQuantity = int(repairedQuantity)
	}

	now := timestamppb.New(time.Now())
	getDataProduct.UpdatedAt = now
	// Convert pb.PosProduct to entity.PosProduct
	gormProduct := &entity.PosProduct{
		ProductID:          uuid.MustParse(getDataProduct.ProductId),
		ProductBarcodeID:   getDataProduct.ProductBarcodeId,
		ProductName:        getDataProduct.ProductName,
		Price:              getDataProduct.Price,
		CostPrice:          getDataProduct.CostPrice,
		CategoryID:         uuid.MustParse(getDataProduct.CategoryId),
		SubCategoryID:      uuid.MustParse(getDataProduct.SubCategoryId),
		StockQuantity:      setQuantity,
		ReorderLevel:       int(getDataProduct.ReorderLevel),
		SupplierID:         uuid.MustParse(getDataProduct.SupplierId),
		ProductDescription: getDataProduct.ProductDescription,
		Active:             getDataProduct.Active,
		StoreID:            uuid.MustParse(getDataProduct.StoreId),
		BranchID:           utils.ParseUUID(getDataProduct.BranchId),
		CompanyID:          uuid.MustParse(getDataProduct.CompanyId),
		CreatedAt:          getDataProduct.CreatedAt.AsTime(),
		CreatedBy:          uuid.MustParse(getDataProduct.CreatedBy),
		UpdatedAt:          getDataProduct.UpdatedAt.AsTime(),
		UpdatedBy:          uuid.MustParse(req.JwtPayload.UserId),
	}

	err = s.repoProduct.UpdatePosProduct(gormProduct)
	if err != nil {
		return nil, err
	}

	err = s.repoInventory.DeletePosInventoryHistory(req.InventoryId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosInventoryHistoryResponse{
		Success: true,
	}, nil
}

func (s *posInventoryHistoryService) ReadAllPosInventoryHistories(ctx context.Context, req *pb.ReadAllPosInventoryHistoriesRequest) (*pb.ReadAllPosInventoryHistoriesResponse, error) {
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
		return nil, errors.New("users are not allowed to read all inventory history")
	}

	paginationResult, err := s.repoInventory.ReadAllPosInventoryHistories(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posInventoryHistories := paginationResult.Records.([]entity.PosInventoryHistory)
	pbPosInventoryHistories := make([]*pb.PosInventoryHistory, len(posInventoryHistories))

	for i, posInventoryHistory := range posInventoryHistories {
		pbPosInventoryHistories[i] = &pb.PosInventoryHistory{
			InventoryId: posInventoryHistory.InventoryID.String(),
			ProductId:   posInventoryHistory.ProductID.String(),
			StoreId:     posInventoryHistory.StoreID.String(),
			Date:        timestamppb.New(posInventoryHistory.Date),
			Quantity:    int32(posInventoryHistory.Quantity),
			BranchId:    posInventoryHistory.BranchID.String(),
			CompanyId:   posInventoryHistory.CompanyID.String(),
			CreatedAt:   timestamppb.New(posInventoryHistory.CreatedAt),
			CreatedBy:   posInventoryHistory.CreatedBy.String(),
			UpdatedAt:   timestamppb.New(posInventoryHistory.UpdatedAt),
			UpdatedBy:   posInventoryHistory.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosInventoryHistoriesResponse{
		PosInventoryHistories: pbPosInventoryHistories,
		Limit:                 int32(pagination.Limit),
		Page:                  int32(pagination.Page),
		MaxPage:               int32(paginationResult.TotalPages),
		Count:                 paginationResult.TotalRecords,
	}, nil
}
