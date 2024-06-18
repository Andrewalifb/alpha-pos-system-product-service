package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosProductSubCategoryRoutes(r *gin.Engine, posProductSubCategoryController controller.PosProductSubCategoryController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/product-sub-categories")
	// Create New PosProductSubCategory
	routesV1.POST("/pos_product_sub_category", posProductSubCategoryController.HandleCreatePosProductSubCategoryRequest)
	// Get PosProductSubCategory by ID
	routesV1.GET("/pos_product_sub_category/:id", posProductSubCategoryController.HandleReadPosProductSubCategoryRequest)
	// Update Existing PosProductSubCategory
	routesV1.PUT("/pos_product_sub_category/:id", posProductSubCategoryController.HandleUpdatePosProductSubCategoryRequest)
	// Delete PosProductSubCategory
	routesV1.DELETE("/pos_product_sub_category/:id", posProductSubCategoryController.HandleDeletePosProductSubCategoryRequest)
	// Get All PosProductSubCategories
	routesV1.GET("/pos_product_sub_categories", posProductSubCategoryController.HandleReadAllPosProductSubCategoriesRequest)
}
