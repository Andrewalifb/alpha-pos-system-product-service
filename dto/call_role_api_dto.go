package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
)

type RoleApiRequest struct {
	RoleID string `json:"role_id"`
	*pb.JWTPayload
}

type RoleApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PosRole `json:"pos_role"`
	} `json:"data"`
}

type PosRole struct {
	RoleID    string   `json:"role_id"`
	RoleName  string   `json:"role_name"`
	CreatedAt JSONTime `json:"created_at"`
	CreatedBy string   `json:"created_by"`
	UpdatedAt JSONTime `json:"updated_at"`
	UpdatedBy string   `json:"updated_by"`
}
