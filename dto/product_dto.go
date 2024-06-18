package dto

import "errors"

// PRODUCT Failed Messages
const (
	MESSAGE_FAILED_CREATE_PRODUCT = "failed to create product"
	MESSAGE_FAILED_UPDATE_PRODUCT = "failed to update product"
	MESSAGE_FAILED_DELETE_PRODUCT = "failed to delete product"
	MESSAGE_FAILED_GET_PRODUCT    = "failed to get product"
)

// PRODUCT Success Messages
const (
	MESSAGE_SUCCESS_CREATE_PRODUCT = "success create product"
	MESSAGE_SUCCESS_UPDATE_PRODUCT = "success update product"
	MESSAGE_SUCCESS_DELETE_PRODUCT = "success delete product"
	MESSAGE_SUCCESS_GET_PRODUCT    = "success get product"
)

// PRODUCT Custom Errors
var (
	ErrCreateProduct = errors.New(MESSAGE_FAILED_CREATE_PRODUCT)
	ErrUpdateProduct = errors.New(MESSAGE_FAILED_UPDATE_PRODUCT)
	ErrDeleteProduct = errors.New(MESSAGE_FAILED_DELETE_PRODUCT)
	ErrGetProduct    = errors.New(MESSAGE_FAILED_GET_PRODUCT)
)
