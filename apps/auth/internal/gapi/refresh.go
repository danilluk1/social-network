package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GAPI) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	// mtdt := server.extractMetadata(ctx)

	payload, err := server.services.TokenMaker.VerifyToken(req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	user, err := server.services.Store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "We don't have info about this user")
		}
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	accessToken, _, err := server.services.TokenMaker.CreateToken(
		user.Username,
		30*time.Minute,
	)
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	refreshToken, refreshPayload, err := server.services.TokenMaker.CreateToken(
		user.Username,
		30*24*time.Hour,
	)
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	currSession, err := server.services.Store.GetSessionByRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	// if currSession.ClientIp != mtdt.ClientIP || currSession.UserAgent != mtdt.UserAgent {
	// 	return nil, status.Errorf(codes.Unauthenticated, "Client IP or UserAgent changed")
	// }

	session, err := server.services.Store.UpdateSession(ctx, db.UpdateSessionParams{
		ID:           currSession.ID,
		RefreshToken: pgtype.Text{String: refreshToken, Valid: true},
		UserAgent:    pgtype.Text{String: "", Valid: false},
		ClientIp:     pgtype.Text{String: "", Valid: false},
		IsBlocked:    pgtype.Bool{Bool: currSession.IsBlocked, Valid: false},
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
	})
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return &auth.RefreshResponse{
		RefreshToken: session.RefreshToken,
		AccessToken:  accessToken,
	}, nil
}
