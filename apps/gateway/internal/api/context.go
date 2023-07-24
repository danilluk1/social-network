package api

import (
	"context"
)

type contextKey string

func (c contextKey) String() string {
	return "gateway api context key " + string(c)
}

const (
	tokenKey = contextKey("paseto")
)

func withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func getToken(ctx context.Context) *string {
	obj := ctx.Value(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*string)
}
