package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosPromotionRoutes(r *gin.Engine, posPromotionController controller.PosPromotionController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/promotions")
	// Create New PosPromotion
	routesV1.POST("/pos_promotion", posPromotionController.HandleCreatePosPromotionRequest)
	// Get PosPromotion by ID
	routesV1.GET("/pos_promotion/:id", posPromotionController.HandleReadPosPromotionRequest)
	// Update Existing PosPromotion
	routesV1.PUT("/pos_promotion/:id", posPromotionController.HandleUpdatePosPromotionRequest)
	// Delete PosPromotion
	routesV1.DELETE("/pos_promotion/:id", posPromotionController.HandleDeletePosPromotionRequest)
	// Get All PosPromotions
	routesV1.GET("/pos_promotions", posPromotionController.HandleReadAllPosPromotionsRequest)
	// Get PosPromotion by Product ID
	routesV1.GET("/pos_promotion/by_product/:product_id", posPromotionController.HandleReadPosPromotionByProductIdRequest)
}
