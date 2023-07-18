package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"fmt"

	"github.com/danilluk1/social-network/apps/api/internal/graph/generated"
	"github.com/danilluk1/social-network/apps/api/internal/graph/model"
)

// VerifyEmail is the resolver for the verifyEmail field.
func (r *mutationResolver) VerifyEmail(ctx context.Context, token string, code *string) (*model.EmailVerificationResult, error) {
	panic(fmt.Errorf("not implemented: VerifyEmail - verifyEmail"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
