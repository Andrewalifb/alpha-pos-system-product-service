package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-product-service/api/midlleware"

	"github.com/gin-gonic/gin"
)

func PosSupplierRoutes(r *gin.Engine, posSupplierController controller.PosSupplierController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/suppliers")
	// Create New PosSupplier
	routesV1.POST("/pos_supplier", posSupplierController.HandleCreatePosSupplierRequest)
	// Get PosSupplier by ID
	routesV1.GET("/pos_supplier/:id", posSupplierController.HandleReadPosSupplierRequest)
	// Update Existing PosSupplier
	routesV1.PUT("/pos_supplier/:id", posSupplierController.HandleUpdatePosSupplierRequest)
	// Delete PosSupplier
	routesV1.DELETE("/pos_supplier/:id", posSupplierController.HandleDeletePosSupplierRequest)
	// Get All PosSuppliers
	routesV1.GET("/pos_suppliers", posSupplierController.HandleReadAllPosSuppliersRequest)
}
