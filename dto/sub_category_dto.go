package dto

import "errors"

// SUB_CATEGORY Failed Messages
const (
	MESSAGE_FAILED_CREATE_SUB_CATEGORY = "failed to create sub category"
	MESSAGE_FAILED_UPDATE_SUB_CATEGORY = "failed to update sub category"
	MESSAGE_FAILED_DELETE_SUB_CATEGORY = "failed to delete sub category"
	MESSAGE_FAILED_GET_SUB_CATEGORY    = "failed to get sub category"
)

// SUB_CATEGORY Success Messages
const (
	MESSAGE_SUCCESS_CREATE_SUB_CATEGORY = "success create sub category"
	MESSAGE_SUCCESS_UPDATE_SUB_CATEGORY = "success update sub category"
	MESSAGE_SUCCESS_DELETE_SUB_CATEGORY = "success delete sub category"
	MESSAGE_SUCCESS_GET_SUB_CATEGORY    = "success get sub category"
)

// SUB_CATEGORY Custom Errors
var (
	ErrCreateSubCategory = errors.New(MESSAGE_FAILED_CREATE_SUB_CATEGORY)
	ErrUpdateSubCategory = errors.New(MESSAGE_FAILED_UPDATE_SUB_CATEGORY)
	ErrDeleteSubCategory = errors.New(MESSAGE_FAILED_DELETE_SUB_CATEGORY)
	ErrGetSubCategory    = errors.New(MESSAGE_FAILED_GET_SUB_CATEGORY)
)
