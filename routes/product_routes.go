package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosProductRoutes(r *gin.Engine, posProductController controller.PosProductController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/products")
	// Create New PosProduct
	routesV1.POST("/pos_product", posProductController.HandleCreatePosProductRequest)
	// Get PosProduct by ID
	routesV1.GET("/pos_product/:id", posProductController.HandleReadPosProductRequest)
	// get PosProduct by Barcode ID
	routesV1.GET("/pos_product_barcode/:id", posProductController.HandleReadPosProductBarcodeRequest)
	// Update Existing PosProduct
	routesV1.PUT("/pos_product/:id", posProductController.HandleUpdatePosProductRequest)
	// Delete PosProduct
	routesV1.DELETE("/pos_product/:id", posProductController.HandleDeletePosProductRequest)
	// Get All PosProducts
	routesV1.GET("/pos_products", posProductController.HandleReadAllPosProductsRequest)
}
