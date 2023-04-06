package auth

import (
	"context"

	"github.com/danilluk1/social-network/apps/gateway/internal/di"
	"github.com/danilluk1/social-network/apps/gateway/internal/helpers"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

func createUserService(dto *createUserDto) (*userResponse, error) {
	authService := do.MustInvoke[auth.AuthClient](di.Provider)

	res, err := authService.CreateUser(context.Background(), &auth.CreateUserRequest{
		Username: dto.Username,
		FullName: dto.FullName,
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot create user")
		return nil, helpers.GetFiberErrorFromGrpcError(err)
	}

	return &userResponse{
		Username:          res.User.Username,
		FullName:          res.User.FullName,
		Email:             res.User.Email,
		PasswordChangedAt: res.User.PasswordChangedAt.AsTime(),
		CreatedAt:         res.User.CreatedAt.AsTime(),
	}, nil
}
