package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/config"
	"github.com/Andrewalifb/alpha-pos-system-product-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-product-service/pkg/service"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	// Initialize the database
	dbConfig := config.NewConfig()
	grpcConfig := config.NewGRPCConfig()

	// Initialize the repositories
	productCategoryRepo := repository.NewPosProductCategoryRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	inventoryHistoryRepo := repository.NewPosInventoryHistoryRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	productRepo := repository.NewPosProductRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	promotionRepo := repository.NewPosPromotionRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	productSubCategoryRepo := repository.NewPosProductSubCategoryRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	supplierRepo := repository.NewPosSupplierRepository(dbConfig.SQLDB, dbConfig.RedisDB)

	// Initialize the services
	productCategorySvc := service.NewPosProductCategoryService(productCategoryRepo, grpcConfig.CompanyServiceConn)
	inventoryHistorySvc := service.NewPosInventoryHistoryService(inventoryHistoryRepo, productRepo, grpcConfig.CompanyServiceConn)
	productSvc := service.NewPosProductService(productRepo, supplierRepo, productCategoryRepo, productSubCategoryRepo, grpcConfig.CompanyServiceConn)
	promotionSvc := service.NewPosPromotionService(promotionRepo, grpcConfig.CompanyServiceConn)
	productSubCategorySvc := service.NewPosProductSubCategoryService(productSubCategoryRepo, productCategoryRepo, grpcConfig.CompanyServiceConn)
	supplierSvc := service.NewPosSupplierService(supplierRepo, grpcConfig.CompanyServiceConn)

	// Create a gRPC server
	s := grpc.NewServer()

	// Register the services with the gRPC server
	pb.RegisterPosProductCategoryServiceServer(s, productCategorySvc)
	pb.RegisterPosInventoryHistoryServiceServer(s, inventoryHistorySvc)
	pb.RegisterPosProductServiceServer(s, productSvc)
	pb.RegisterPosPromotionServiceServer(s, promotionSvc)
	pb.RegisterPosProductSubCategoryServiceServer(s, productSubCategorySvc)
	pb.RegisterPosSupplierServiceServer(s, supplierSvc)

	// Start the gRPC server
	serverPort := os.Getenv("SERVER_PORT")
	lis, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
