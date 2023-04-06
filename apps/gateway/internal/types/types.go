package types

import (
	"github.com/danilluk1/social-network/apps/gateway/internal/services/redis"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Services struct {
	RedisStorage        *redis.RedisStorage
	Validator           *validator.Validate
	ValidatorTranslator ut.Translator
}
