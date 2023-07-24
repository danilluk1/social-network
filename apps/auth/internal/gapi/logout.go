package gapi

import (
	"context"
	"errors"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GAPI) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	_, err := server.services.TokenMaker.VerifyToken(req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	var uuid pgtype.UUID
	err = uuid.Scan(req.GetSessionId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "SessionID must be UUID string")
	}

	_, err = server.services.Store.UpdateSession(ctx, db.UpdateSessionParams{
		ID:           uuid,
		RefreshToken: pgtype.Text{String: "", Valid: true},
	})
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return &auth.LogoutResponse{}, nil
}
