package controller

import (
	"net/http"
	"strconv"
	"strings"

	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-product-service/utils"

	"github.com/gin-gonic/gin"
)

type PosProductController interface {
	HandleCreatePosProductRequest(c *gin.Context)
	HandleReadPosProductRequest(c *gin.Context)
	HandleReadPosProductBarcodeRequest(c *gin.Context)
	HandleUpdatePosProductRequest(c *gin.Context)
	HandleDeletePosProductRequest(c *gin.Context)
	HandleReadAllPosProductsRequest(c *gin.Context)
}

type posProductController struct {
	service pb.PosProductServiceClient
}

func NewPosProductController(service pb.PosProductServiceClient) PosProductController {
	return &posProductController{
		service: service,
	}
}

func (ctrl *posProductController) HandleCreatePosProductRequest(c *gin.Context) {
	var req pb.CreatePosProductRequest

	if err := c.ShouldBindJSON(&req.PosProduct); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRODUCT, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.CreatePosProduct(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductController) HandleReadPosProductRequest(c *gin.Context) {
	var req pb.ReadPosProductRequest

	productID := c.Param("id")
	req.ProductId = productID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODUCT, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.ReadPosProduct(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODUCT, resp)
	c.JSON(http.StatusOK, successResponse)
}

func (ctrl *posProductController) HandleReadPosProductBarcodeRequest(c *gin.Context) {
	var req pb.ReadPosProductByBarcodeRequest

	barcodeID := c.Param("id")
	req.ProductBarcodeId = barcodeID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODUCT, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.ReadPosProductByBarcode(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODUCT, resp)
	c.JSON(http.StatusOK, successResponse)
}

func (ctrl *posProductController) HandleUpdatePosProductRequest(c *gin.Context) {
	var req pb.UpdatePosProductRequest
	productID := c.Param("id")
	if err := c.ShouldBindJSON(&req.PosProduct); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODUCT, "Jwt Payload is Empty", nil)
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
	req.PosProduct.ProductId = productID
	resp, err := ctrl.service.UpdatePosProduct(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductController) HandleDeletePosProductRequest(c *gin.Context) {
	var req pb.DeletePosProductRequest

	productID := c.Param("id")
	req.ProductId = productID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PRODUCT, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.DeletePosProduct(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductController) HandleReadAllPosProductsRequest(c *gin.Context) {
	limitQuery := c.Query("limit")
	pageQuery := c.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosProductResponse{})
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosProductsRequest

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

		req = pb.ReadAllPosProductsRequest{
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
	res, err := ctrl.service.ReadAllPosProducts(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODUCT, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
