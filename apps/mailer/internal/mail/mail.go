package mail

import (
	"context"

	"github.com/danilluk1/social-network/apps/mailer/internal/types"
)

type Reader struct {
	services *types.Services
}

func New(services *types.Services) *Reader {
	return &Reader{
		services: services,
	}
}

func (r *Reader) Start(ctx context.Context) {
	for {
		msg, err := r.services.Reader.ReadMessage(ctx)
		if err != nil {
			//TODO: provide rich error handling
			r.services.Logger.Sugar().Error(err)
		}
	}
}

func (r *Reader) sendMail {
	
}
