package createUserCommand

import "github.com/danilluk1/social-network/libs/go/pkg/logger"

type CreateUserHandler struct {
	log    logger.Logger
	tracer tracing
}
