package utils

import "os"

func IsSuperUser(userRole string) bool {
	return userRole == os.Getenv("SUPER_USER_ROLE")
}

func IsCompanyUser(userRole string) bool {
	return userRole == os.Getenv("COMPANY_USER_ROLE")
}

func IsBranchUser(userRole string) bool {
	return userRole == os.Getenv("BRANCH_USER_ROLE")
}

func IsStoreUser(userRole string) bool {
	return userRole == os.Getenv("STORE_USER_ROLE")
}

func IsCompanyOrBranchUser(userRole string) bool {
	return userRole == os.Getenv("COMPANY_USER_ROLE") || userRole == os.Getenv("BRANCH_USER_ROLE")
}

func IsBranchOrStoreUser(userRole string) bool {
	return userRole == os.Getenv("BRANCH_USER_ROLE") || userRole == os.Getenv("STORE_USER_ROLE")
}

func IsSuperUserOrCompany(userRole string) bool {
	return userRole == os.Getenv("SUPER_USER_ROLE") || userRole == os.Getenv("COMPANY_USER_ROLE")
}

func IsCompanyOrBranchOrStoreUser(userRole string) bool {
	return userRole == os.Getenv("COMPANY_USER_ROLE") || userRole == os.Getenv("BRANCH_USER_ROLE") || userRole == os.Getenv("STORE_USER_ROLE")
}

func VerifyStoreUserAccess(loginRole, branchId, jwtBranchId string) bool {
	return loginRole == os.Getenv("STORE_USER_ROLE") && branchId == jwtBranchId
}

func VerifyBranchUserAccess(loginRole, branchId, jwtBranchId string) bool {
	return loginRole == os.Getenv("BRANCH_USER_ROLE") && branchId == jwtBranchId
}

func VerifyCompanyUserAccess(loginRole, companyId, jwtCompanyId string) bool {
	return loginRole == os.Getenv("COMPANY_USER_ROLE") && companyId == jwtCompanyId
}
