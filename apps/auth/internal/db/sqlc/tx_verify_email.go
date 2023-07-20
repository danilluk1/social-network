package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
)

type VerifyEmailTxParams struct {
	SecretCode *string
	Token      *string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.VerifyEmail(ctx, VerifyEmailParams{
			SecretCode: lo.FromPtrOr(arg.SecretCode, ""),
			Token:      lo.FromPtrOr(arg.Token, ""),
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
