package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/utils"

	"github.com/gin-gonic/gin"
)

type PosInventoryHistoryController interface {
	HandleCreatePosInventoryHistoryRequest(c *gin.Context)
	HandleReadPosInventoryHistoryRequest(c *gin.Context)
	HandleUpdatePosInventoryHistoryRequest(c *gin.Context)
	HandleDeletePosInventoryHistoryRequest(c *gin.Context)
	HandleReadAllPosInventoryHistoriesRequest(c *gin.Context)
}

type posInventoryHistoryController struct {
	service pb.PosInventoryHistoryServiceClient
}

func NewPosInventoryHistoryController(service pb.PosInventoryHistoryServiceClient) PosInventoryHistoryController {
	return &posInventoryHistoryController{
		service: service,
	}
}

func (ctrl *posInventoryHistoryController) HandleCreatePosInventoryHistoryRequest(c *gin.Context) {
	var req pb.CreatePosInventoryHistoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVENTORY_HISTORY, "Jwt Payload is Empty", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token
	fmt.Println("REQUEST BODY ALL : ", &req)
	resp, err := ctrl.service.CreatePosInventoryHistory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_INVENTORY_HISTORY, resp)
	c.JSON(http.StatusOK, successResponse)
}

func (ctrl *posInventoryHistoryController) HandleReadPosInventoryHistoryRequest(c *gin.Context) {
	var req pb.ReadPosInventoryHistoryRequest

	inventoryID := c.Param("id")
	req.InventoryId = inventoryID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVENTORY_HISTORY, "Jwt Payload is Empty", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	resp, err := ctrl.service.ReadPosInventoryHistory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posInventoryHistoryController) HandleUpdatePosInventoryHistoryRequest(c *gin.Context) {
	var req pb.UpdatePosInventoryHistoryRequest

	inventoryID := c.Param("id")

	if err := c.ShouldBindJSON(&req.PosInventoryHistory); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVENTORY_HISTORY, "Jwt Payload is Empty", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token
	req.PosInventoryHistory.InventoryId = inventoryID
	resp, err := ctrl.service.UpdatePosInventoryHistory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posInventoryHistoryController) HandleDeletePosInventoryHistoryRequest(c *gin.Context) {
	var req pb.DeletePosInventoryHistoryRequest

	inventoryID := c.Param("id")
	req.InventoryId = inventoryID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVENTORY_HISTORY, "Jwt Payload is Empty", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	resp, err := ctrl.service.DeletePosInventoryHistory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posInventoryHistoryController) HandleReadAllPosInventoryHistoriesRequest(c *gin.Context) {
	limitQuery := c.Query("limit")
	pageQuery := c.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosInventoryHistoryResponse{})
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosInventoryHistoriesRequest

	if limitQuery != "" && pageQuery != "" {
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			errorResponse := utils.BuildResponseFailed("Invalid limit value", err.Error(), nil)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			errorResponse := utils.BuildResponseFailed("Invalid page value", err.Error(), nil)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosInventoryHistoriesRequest{
			Limit: int32(limit),
			Page:  int32(page),
		}
	}
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CATEGORY, "Jwt Payload is Empty", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	res, err := ctrl.service.ReadAllPosInventoryHistories(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVENTORY_HISTORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
