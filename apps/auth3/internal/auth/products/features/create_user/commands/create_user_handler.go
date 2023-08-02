package createUserCommand

import (
	"context"

	"github.com/danilluk1/social-network/apps/auth/internal/auth/products/features/create_user/dtos"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"github.com/danilluk1/social-network/libs/go/pkg/otel/tracing"
	attribute2 "github.com/danilluk1/social-network/libs/go/pkg/otel/tracing/attribute"
	pass_utils "github.com/danilluk1/social-network/libs/go/pkg/utils/password"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserHandler struct {
	log    logger.Logger
	tracer tracing.AppTracer
}

func NewCreateUserHandler(
	log logger.Logger,
	tracer tracing.AppTracer,
) *CreateUserHandler {
	return &CreateUserHandler{
		tracer: tracer,
		log:    log,
	}
}

func (c *CreateUserHandler) Handle(
	ctx context.Context,
	command *CreateUser,
) (*dtos.CreateUserResponse, error) {
	ctx, span := c.tracer.Start(ctx, "CreateUserHandler.Handle")
	span.SetAttributes(attribute.String("Username", command.Username))
	span.SetAttributes(attribute.String("Email", command.Email))
	span.SetAttributes(attribute2.Object("Command", command))
	defer span.End()

	var createProductResult *dtos.CreateUserResponse

	hashedPassword, err := pass_utils.HashPassword(command.Password)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
		)
		server.services.Logger.Sugar().Error(err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}
}
