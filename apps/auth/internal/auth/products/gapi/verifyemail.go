package gapi

import (
	"context"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/jackc/pgx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GAPI) VerifyEmail(ctx context.Context, req *auth.VerifyEmailRequest) (*auth.VerifyEmailResponse, error) {
	txRes, err := server.services.Store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		SecretCode: &req.SecretCode,
		Token:      &req.Token,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "We don't have info about this activation message")
		}
		server.services.Logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	return &auth.VerifyEmailResponse{
		Username:    txRes.User.Username,
		Email:       txRes.User.Email,
		IsActivated: txRes.User.IsEmailVerified,
	}, nil
}
