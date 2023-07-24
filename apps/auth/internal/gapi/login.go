package gapi

import (
	"context"
	"fmt"
	"time"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/val"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/danilluk1/social-network/libs/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func validateLoginUserRequest(req *auth.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return
}

func (server *GAPI) LoginUser(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	mtdt := server.extractMetadata(ctx)

	user, err := server.services.Store.GetUser(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Can't find user with this username")
		}

		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Bad password")
	}

	accessToken, accessPayload, err := server.services.TokenMaker.CreateToken(
		user.Username,
		30*time.Minute,
	)
	if err != nil {
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

	session, err := server.services.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshPayload.ID, Valid: true},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
	})
	if err != nil {
		// server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	rsp := &auth.LoginUserResponse{
		AccessToken:           accessToken,
		SessionId:             fmt.Sprintf("%x-%x-%x-%x-%x", session.ID.Bytes[0:4], session.ID.Bytes[4:6], session.ID.Bytes[6:8], session.ID.Bytes[8:10], session.ID.Bytes[10:16]),
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		User: &auth.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			CreatedAt:         timestamppb.New(user.CreatedAt.Time),
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		},
	}
	return rsp, nil
}
