package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosRole makes a gRPC request to the specified service and returns the response
func GetPosRoleById(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload) (*pos.ReadPosRoleResponse, error) {

	// Create a new PosRoleService client
	client := pos.NewPosRoleServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosRoleRequest{
		RoleId:     id,
		JwtPayload: jwtPayload,
	}

	// Call the gRPC method
	resp, err := client.ReadPosRole(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
