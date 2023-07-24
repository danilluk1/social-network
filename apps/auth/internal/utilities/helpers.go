package utilities

import (
	"math/rand"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GenerateEasyCode() int {
	easyDigits := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	code := 0
	for i := 0; i < 4; i++ {
		digit := easyDigits[rand.Intn(len(easyDigits))]
		code = code*10 + digit
	}

	return code
}

func ConvertUser(user db.User) *auth.User {
	return &auth.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}
}
