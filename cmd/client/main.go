package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	clientPort := os.Getenv("CLIENT_PORT")
	serverPort := os.Getenv("SERVER_PORT")

	addr := fmt.Sprintf("localhost:%s", serverPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC clients for each service
	productCategoryClient := pb.NewPosProductCategoryServiceClient(conn)
	inventoryHistoryClient := pb.NewPosInventoryHistoryServiceClient(conn)
	productClient := pb.NewPosProductServiceClient(conn)
	promotionClient := pb.NewPosPromotionServiceClient(conn)
	productSubCategoryClient := pb.NewPosProductSubCategoryServiceClient(conn)
	supplierClient := pb.NewPosSupplierServiceClient(conn)

	// Initialize the controllers with the gRPC clients
	productCategoryCtrl := controller.NewPosProductCategoryController(productCategoryClient)
	inventoryHistoryCtrl := controller.NewPosInventoryHistoryController(inventoryHistoryClient)
	productCtrl := controller.NewPosProductController(productClient)
	promotionCtrl := controller.NewPosPromotionController(promotionClient)
	productSubCategoryCtrl := controller.NewPosProductSubCategoryController(productSubCategoryClient)
	supplierCtrl := controller.NewPosSupplierController(supplierClient)

	// Create a new router
	r := gin.Default()

	// Define your routes
	routes.PosProductCategoryRoutes(r, productCategoryCtrl)
	routes.PosInventoryHistoryRoutes(r, inventoryHistoryCtrl)
	routes.PosProductRoutes(r, productCtrl)
	routes.PosPromotionRoutes(r, promotionCtrl)
	routes.PosProductSubCategoryRoutes(r, productSubCategoryCtrl)
	routes.PosSupplierRoutes(r, supplierCtrl)

	// Start the server
	r.Run(":" + clientPort)
}
