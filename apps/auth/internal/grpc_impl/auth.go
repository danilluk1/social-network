package grpc_impl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/helpers"
	"github.com/danilluk1/social-network/apps/auth/internal/val"
	"github.com/danilluk1/social-network/libs/avro"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	types "github.com/danilluk1/social-network/libs/kafka/schemas"
	"github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/danilluk1/social-network/libs/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) VerifyEmail(ctx context.Context, req *auth.VerifyEmailRequest) (*auth.VerifyEmailResponse, error) {
	txRes, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		SecretCode: &req.SecretCode,
		Token:      &req.Token,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "We don't have info about this activation message")
		}
		server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	return &auth.VerifyEmailResponse{
		Username:    txRes.User.Username,
		Email:       txRes.User.Email,
		IsActivated: txRes.User.IsEmailVerified,
	}, nil
}

func (server *Server) LoginUser(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	mtdt := server.extractMetadata(ctx)

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Can't find user with this username")
		}
		server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Bad password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		30*time.Minute,
	)
	if err != nil {
		server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		30*24*time.Hour,
	)
	if err != nil {
		server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshPayload.ID, Valid: true},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
	})
	if err != nil {
		server.logger.Sugar().Error(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	rsp := &auth.LoginUserResponse{
		AccessToken:           accessToken,
		SessionId:             fmt.Sprintf("%x-%x-%x-%x-%x", session.ID.Bytes[0:4], session.ID.Bytes[4:6], session.ID.Bytes[6:8], session.ID.Bytes[8:10], session.ID.Bytes[10:16]),
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		User: &auth.User{
			Username:          user.FullName,
			FullName:          user.FullName,
			Email:             user.Email,
			CreatedAt:         timestamppb.New(user.CreatedAt.Time),
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		},
	}
	return rsp, nil
}

func validateLoginUserRequest(req *auth.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return
}

func (server *Server) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		server.logger.Sugar().Error(err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			Email:          req.GetEmail(),
			FullName:       req.GetFullName(),
		},
		AfterCreate: func(user db.User) error {
			verifyEmail, err := server.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
				Username:   user.Username,
				Email:      user.Email,
				SecretCode: strconv.Itoa(helpers.GenerateEasyCode()),
				Token:      utils.RandomString(32),
			})
			if err != nil {
				return err
			}

			verifyUrl := fmt.Sprintf("http://localhost:9090/verify_email?id=%d&token=%s", verifyEmail.ID, verifyEmail.Token)
			email := types.EmailMessage{
				From:    "socialnetwork@mail.ru",
				To:      []string{user.Email},
				Cc:      []string{"dan@example.com"},
				Subject: "Welcome to our Social Network!",
				Body: fmt.Sprintf(`Hello %s,<br/>
				Thank you for registering with us <br/>
				Please <a href="%s">Click here</a> to verify your email address. Or enter secret code %s<br/>
				`, user.FullName, verifyUrl, verifyEmail.SecretCode),
				Attachments: []string{},
			}

			schema, err := server.schemaRegistry.GetLatestSchema(topics.Mail)
			if err != nil {
				return err
			}

			encodedMsg, err := avro.Encode(&email, schema.Codec(), schema.ID())
			if err != nil {
				return err
			}

			err = server.kafkaWriter.WriteMessages(ctx, kafka.Message{
				Value: encodedMsg,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	tx, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				return nil, status.Error(codes.AlreadyExists, "username already exists")
			}
		}
		server.logger.Sugar().Error(err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	rsp := &auth.CreateUserResponse{
		User: convertUser(tx.User),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *auth.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
