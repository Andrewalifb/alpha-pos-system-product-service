package entity

import "github.com/dgrijalva/jwt-go"

type JWTPayload struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	CompanyID string `json:"companyId"`
	BranchID  string `json:"branchId"`
	StoreID   string `json:"storeId"`
	UserID    string `json:"userId"`
	jwt.StandardClaims
}
