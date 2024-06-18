package dto

import "errors"

// INVENTORY_HISTORY Failed Messages
const (
	MESSAGE_FAILED_CREATE_INVENTORY_HISTORY = "failed to create inventory history"
	MESSAGE_FAILED_UPDATE_INVENTORY_HISTORY = "failed to update inventory history"
	MESSAGE_FAILED_DELETE_INVENTORY_HISTORY = "failed to delete inventory history"
	MESSAGE_FAILED_GET_INVENTORY_HISTORY    = "failed to get inventory history"
)

// INVENTORY_HISTORY Success Messages
const (
	MESSAGE_SUCCESS_CREATE_INVENTORY_HISTORY = "success create inventory history"
	MESSAGE_SUCCESS_UPDATE_INVENTORY_HISTORY = "success update inventory history"
	MESSAGE_SUCCESS_DELETE_INVENTORY_HISTORY = "success delete inventory history"
	MESSAGE_SUCCESS_GET_INVENTORY_HISTORY    = "success get inventory history"
)

// INVENTORY_HISTORY Custom Errors
var (
	ErrCreateInventoryHistory = errors.New(MESSAGE_FAILED_CREATE_INVENTORY_HISTORY)
	ErrUpdateInventoryHistory = errors.New(MESSAGE_FAILED_UPDATE_INVENTORY_HISTORY)
	ErrDeleteInventoryHistory = errors.New(MESSAGE_FAILED_DELETE_INVENTORY_HISTORY)
	ErrGetInventoryHistory    = errors.New(MESSAGE_FAILED_GET_INVENTORY_HISTORY)
)
