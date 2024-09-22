package auth

import (
	"context"

	pb "github.com/cstati/auth/pkg/auth"
)

func (s *Service) Ping(ctx context.Context, r *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}
