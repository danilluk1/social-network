package gapi

import (
	"context"
	"fmt"
	"strconv"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/utilities"
	"github.com/danilluk1/social-network/apps/auth/internal/val"
	"github.com/danilluk1/social-network/libs/avro"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	types "github.com/danilluk1/social-network/libs/kafka/schemas"
	"github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/danilluk1/social-network/libs/utils"
	"github.com/jackc/pgconn"
	"github.com/segmentio/kafka-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (server *GAPI) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		server.services.Logger.Sugar().Error(err)
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
			verifyEmail, err := server.services.Store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
				Username:   user.Username,
				Email:      user.Email,
				SecretCode: strconv.Itoa(utilities.GenerateEasyCode()),
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

			schema, err := server.services.SchemaClient.GetLatestSchema(topics.Mail)
			if err != nil {
				return err
			}

			encodedMsg, err := avro.Encode(&email, schema.Codec(), schema.ID())
			if err != nil {
				return err
			}

			err = server.services.KafkaWriter.WriteMessages(ctx, kafka.Message{
				Value: encodedMsg,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	tx, err := server.services.Store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				return nil, status.Error(codes.AlreadyExists, "username already exists")
			}
		}
		server.services.Logger.Sugar().Error(err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	rsp := &auth.CreateUserResponse{
		User: utilities.ConvertUser(tx.User),
	}
	return rsp, nil
}
