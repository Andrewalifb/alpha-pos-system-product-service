package dto

import "errors"

// PROMOTION Failed Messages
const (
	MESSAGE_FAILED_CREATE_PROMOTION = "failed to create promotion"
	MESSAGE_FAILED_UPDATE_PROMOTION = "failed to update promotion"
	MESSAGE_FAILED_DELETE_PROMOTION = "failed to delete promotion"
	MESSAGE_FAILED_GET_PROMOTION    = "failed to get promotion"
)

// PROMOTION Success Messages
const (
	MESSAGE_SUCCESS_CREATE_PROMOTION = "success create promotion"
	MESSAGE_SUCCESS_UPDATE_PROMOTION = "success update promotion"
	MESSAGE_SUCCESS_DELETE_PROMOTION = "success delete promotion"
	MESSAGE_SUCCESS_GET_PROMOTION    = "success get promotion"
)

// PROMOTION Custom Errors
var (
	ErrCreatePromotion = errors.New(MESSAGE_FAILED_CREATE_PROMOTION)
	ErrUpdatePromotion = errors.New(MESSAGE_FAILED_UPDATE_PROMOTION)
	ErrDeletePromotion = errors.New(MESSAGE_FAILED_DELETE_PROMOTION)
	ErrGetPromotion    = errors.New(MESSAGE_FAILED_GET_PROMOTION)
)
