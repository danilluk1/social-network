package createUserCommand

import (
	"github.com/danilluk1/social-network/apps/auth/internal/shared/validators"
	grpcErrors "github.com/danilluk1/social-network/libs/go/pkg/grpc/errors"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/go-playground/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type CreateUser struct {
	Username string
	Password string
	FullName string
	Email    string
}

func validateCreateUserRequest(req *auth.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, grpcErrors.FieldViolation("username", err))
	}

	if err := validators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, grpcErrors.FieldViolation("password", err))
	}

	if err := validators.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, grpcErrors.FieldViolation("full_name", err))
	}

	if err := validators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, grpcErrors.FieldViolation("email", err))
	}

	return violations
}

func NewCreateUser(username, password, fullName, email string) (*CreateUser, error) {
	command := &CreateUser{Username: username, Password: password, FullName: fullName, Email: email}
	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
