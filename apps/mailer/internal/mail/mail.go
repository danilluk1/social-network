package mail

import (
	"context"
	"encoding/binary"

	"github.com/danilluk1/social-network/apps/mailer/internal/types"
	"github.com/danilluk1/social-network/libs/avro"
)

type Reader struct {
	services *types.Services
}

func New(services *types.Services) *Reader {
	return &Reader{
		services: services,
	}
}

type EmailMessage struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments"`
}

func (r *Reader) Start(ctx context.Context) {
	for {
		msg, err := r.services.Reader.ReadMessage(ctx)
		if err != nil {
			//TODO: provide rich error handling
			r.services.Logger.Sugar().Error(err)
		}

		if len(msg.Value) <= 5 {
			continue
		}

		schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
		schema, err := r.services.SchemaRegistry.GetSchema(int(schemaID))
		if err != nil {
			r.services.Logger.Sugar().Error(err)
			continue
		}

		email := &EmailMessage{}
		avro.Decode(msg.Value, schema.Codec(), email)
		r.services.Logger.Sugar().Info(email.Body)
	}

}
