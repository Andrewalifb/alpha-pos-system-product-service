package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"

	"github.com/gin-gonic/gin"
)

func PosProductCategoryRoutes(r *gin.Engine, posProductCategoryController controller.PosProductCategoryController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/product-categories")
	// Create New PosProductCategory
	routesV1.POST("/pos_product_category", posProductCategoryController.HandleCreatePosProductCategoryRequest)
	// Get PosProductCategory by ID
	routesV1.GET("/pos_product_category/:id", posProductCategoryController.HandleReadPosProductCategoryRequest)
	// Update Existing PosProductCategory
	routesV1.PUT("/pos_product_category/:id", posProductCategoryController.HandleUpdatePosProductCategoryRequest)
	// Delete PosProductCategory
	routesV1.DELETE("/pos_product_category/:id", posProductCategoryController.HandleDeletePosProductCategoryRequest)
	// Get All PosProductCategories
	routesV1.GET("/pos_product_categories", posProductCategoryController.HandleReadAllPosProductCategoriesRequest)
}
