package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosStoreBranch makes a gRPC request to the specified service and returns the response
func GetPosStoreBranchById(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload) (*pos.ReadPosStoreBranchResponse, error) {

	// Create a new PosStoreBranchService client
	client := pos.NewPosStoreBranchServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosStoreBranchRequest{
		BranchId:   id,
		JwtPayload: jwtPayload,
	}

	// Call the gRPC method
	resp, err := client.ReadPosStoreBranch(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
