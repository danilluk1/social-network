package grpc_impl

import (
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *auth.User {
	return &auth.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}
}
