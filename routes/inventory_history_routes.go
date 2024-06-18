package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"

	"github.com/gin-gonic/gin"
)

func PosInventoryHistoryRoutes(r *gin.Engine, posInventoryHistoryController controller.PosInventoryHistoryController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/inventory-histories")
	// Create New PosInventoryHistory
	routesV1.POST("/pos_inventory_history", posInventoryHistoryController.HandleCreatePosInventoryHistoryRequest)
	// Get PosInventoryHistory by ID
	routesV1.GET("/pos_inventory_history/:id", posInventoryHistoryController.HandleReadPosInventoryHistoryRequest)
	// Update Existing PosInventoryHistory
	routesV1.PUT("/pos_inventory_history/:id", posInventoryHistoryController.HandleUpdatePosInventoryHistoryRequest)
	// Delete PosInventoryHistory
	routesV1.DELETE("/pos_inventory_history/:id", posInventoryHistoryController.HandleDeletePosInventoryHistoryRequest)
	// Get All PosInventoryHistories
	routesV1.GET("/pos_inventory_histories", posInventoryHistoryController.HandleReadAllPosInventoryHistoriesRequest)
}
