package mail

import (
	"context"
	"encoding/binary"

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

		if len(msg.Value) <= 5 {
			return
		}

		schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
		schema, err := r.services.SchemaRegistry.GetSchema(int(schemaID))
		if err != nil {
			r.services.Logger.Sugar().Error(err)
			return
		}

		native, _, _ := schema.Codec().NativeFromBinary(msg.Value[5:])
		value, _ := schema.Codec().TextualFromNative(nil, native)
		r.services.Logger.Sugar().Info(value)
	}

}
