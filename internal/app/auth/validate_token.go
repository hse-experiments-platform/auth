package auth

import (
	"context"
	"strings"

	pb "github.com/cstati/auth/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserInfo struct {
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

func (s *Service) ValidateToken(ctx context.Context, r *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	var token string

	if r.Token != "" {
		token = r.Token
	} else if md, ok := metadata.FromIncomingContext(ctx); ok && len(md.Get("authorization")) == 1 {
		authHeader := md.Get("authorization")[0]
		token, _ = strings.CutPrefix(authHeader, "Bearer ")
	}
	if token == "" {
		return nil, status.New(codes.Unauthenticated, "must provide token in body or in Authorization header").Err()
	}

	payload, err := s.tokenProvider.VerifyToken(token)
	if err != nil {
		return nil, status.New(codes.PermissionDenied, err.Error()).Err()
	}

	return &pb.ValidateTokenResponse{
		UserId: payload.UserID,
		Roles:  payload.Roles,
	}, nil
}
