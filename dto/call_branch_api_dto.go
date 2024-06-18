package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
)

type BranchApiRequest struct {
	BranchID string `json:"branch_id"`
	*pb.JWTPayload
}

type BranchApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PosStoreBranch `json:"pos_store_branch"`
	} `json:"data"`
}

type PosStoreBranch struct {
	BranchID   string   `json:"branch_id"`
	BranchName string   `json:"branch_name"`
	CompanyID  string   `json:"company_id"`
	CreatedAt  JSONTime `json:"created_at"`
	CreatedBy  string   `json:"created_by"`
	UpdatedAt  JSONTime `json:"updated_at"`
	UpdatedBy  string   `json:"updated_by"`
}
