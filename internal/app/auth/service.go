package auth

import (
	"github.com/cstati/auth/internal/pkg/storage/db"
	"github.com/cstati/auth/internal/pkg/storage/google"
	pb "github.com/cstati/auth/pkg/auth"
	"github.com/hse-experiments-platform/library/pkg/utils/token"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	google        google.Storage
	db            db.Storage
	tokenProvider token.Maker
}

func NewService(google google.Storage, db db.Storage, tokenProvider token.Maker) *Service {
	return &Service{
		google:        google,
		db:            db,
		tokenProvider: tokenProvider,
	}
}
