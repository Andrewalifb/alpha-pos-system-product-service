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

type PosProductCategoryController interface {
	HandleCreatePosProductCategoryRequest(c *gin.Context)
	HandleReadPosProductCategoryRequest(c *gin.Context)
	HandleUpdatePosProductCategoryRequest(c *gin.Context)
	HandleDeletePosProductCategoryRequest(c *gin.Context)
	HandleReadAllPosProductCategoriesRequest(c *gin.Context)
}

type posProductCategoryController struct {
	service pb.PosProductCategoryServiceClient
}

func NewPosProductCategoryController(service pb.PosProductCategoryServiceClient) PosProductCategoryController {
	return &posProductCategoryController{
		service: service,
	}
}

func (ctrl *posProductCategoryController) HandleCreatePosProductCategoryRequest(c *gin.Context) {
	var req pb.CreatePosProductCategoryRequest

	if err := c.ShouldBindJSON(&req.PosProductCategory); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CATEGORY, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.CreatePosProductCategory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductCategoryController) HandleReadPosProductCategoryRequest(c *gin.Context) {
	var req pb.ReadPosProductCategoryRequest

	categoryID := c.Param("id")
	req.CategoryId = categoryID

	// Get JWT Payload data from middleware
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

	resp, err := ctrl.service.ReadPosProductCategory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductCategoryController) HandleUpdatePosProductCategoryRequest(c *gin.Context) {
	var req pb.UpdatePosProductCategoryRequest

	// Get role ID from URL
	categoryID := c.Param("id")

	if err := c.ShouldBindJSON(&req.PosProductCategory); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CATEGORY, "Jwt Payload is Empty", nil)
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
	req.PosProductCategory.CategoryId = categoryID
	req.JwtToken = token

	resp, err := ctrl.service.UpdatePosProductCategory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductCategoryController) HandleDeletePosProductCategoryRequest(c *gin.Context) {
	var req pb.DeletePosProductCategoryRequest

	categoryID := c.Param("id")
	req.CategoryId = categoryID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := c.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_CATEGORY, "Jwt Payload is Empty", nil)
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

	resp, err := ctrl.service.DeletePosProductCategory(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *posProductCategoryController) HandleReadAllPosProductCategoriesRequest(c *gin.Context) {
	var req pb.ReadAllPosProductCategoriesRequest

	limitQuery := c.Query("limit")
	pageQuery := c.Query("page")

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

		req = pb.ReadAllPosProductCategoriesRequest{
			Limit: int32(limit),
			Page:  int32(page),
		}
	}

	// Get JWT Payload data from middleware
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

	resp, err := ctrl.service.ReadAllPosProductCategories(c.Request.Context(), &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CATEGORY, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, resp)
}
