package gapi

import (
	"context"

	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
)

func (server *GAPI) ValidateUser(ctx context.Context, req *auth.ValidateUserRequest) (*auth.ValidateUserResponse, error) {
	payload, err := server.services.TokenMaker.VerifyToken(req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return &auth.ValidateUserResponse{
		Id:        payload.ID.String(),
		Username:  payload.Username,
		ExpiresAt: payload.ExpiredAt.Unix(),
		IssuedAt:  payload.IssuedAt.Unix(),
	}, nil
}
