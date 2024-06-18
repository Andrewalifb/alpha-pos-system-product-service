package dto

import "errors"

// CATEGORY Failed Messages
const (
	MESSAGE_FAILED_CREATE_CATEGORY = "failed to create category"
	MESSAGE_FAILED_UPDATE_CATEGORY = "failed to update category"
	MESSAGE_FAILED_DELETE_CATEGORY = "failed to delete category"
	MESSAGE_FAILED_GET_CATEGORY    = "failed to get category"
)

// CATEGORY Success Messages
const (
	MESSAGE_SUCCESS_CREATE_CATEGORY = "success create category"
	MESSAGE_SUCCESS_UPDATE_CATEGORY = "success update category"
	MESSAGE_SUCCESS_DELETE_CATEGORY = "success delete category"
	MESSAGE_SUCCESS_GET_CATEGORY    = "success get category"
)

// CATEGORY Custom Errors
var (
	ErrCreateCategory = errors.New(MESSAGE_FAILED_CREATE_CATEGORY)
	ErrUpdateCategory = errors.New(MESSAGE_FAILED_UPDATE_CATEGORY)
	ErrDeleteCategory = errors.New(MESSAGE_FAILED_DELETE_CATEGORY)
	ErrGetCategory    = errors.New(MESSAGE_FAILED_GET_CATEGORY)
)
