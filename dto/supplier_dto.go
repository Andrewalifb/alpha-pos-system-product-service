package dto

import "errors"

// SUPPLIER Failed Messages
const (
	MESSAGE_FAILED_CREATE_SUPPLIER = "failed to create supplier"
	MESSAGE_FAILED_UPDATE_SUPPLIER = "failed to update supplier"
	MESSAGE_FAILED_DELETE_SUPPLIER = "failed to delete supplier"
	MESSAGE_FAILED_GET_SUPPLIER    = "failed to get supplier"
)

// SUPPLIER Success Messages
const (
	MESSAGE_SUCCESS_CREATE_SUPPLIER = "success create supplier"
	MESSAGE_SUCCESS_UPDATE_SUPPLIER = "success update supplier"
	MESSAGE_SUCCESS_DELETE_SUPPLIER = "success delete supplier"
	MESSAGE_SUCCESS_GET_SUPPLIER    = "success get supplier"
)

// SUPPLIER Custom Errors
var (
	ErrCreateSupplier = errors.New(MESSAGE_FAILED_CREATE_SUPPLIER)
	ErrUpdateSupplier = errors.New(MESSAGE_FAILED_UPDATE_SUPPLIER)
	ErrDeleteSupplier = errors.New(MESSAGE_FAILED_DELETE_SUPPLIER)
	ErrGetSupplier    = errors.New(MESSAGE_FAILED_GET_SUPPLIER)
)
