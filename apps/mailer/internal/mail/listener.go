package mail

import (
	"context"
	"encoding/binary"

	"github.com/danilluk1/social-network/libs/avro"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type Reader struct {
	services *Services
}

type Services struct {
	Logger         *zap.Logger
	Mail           *gomail.Dialer
	Reader         *kafka.Reader
	SchemaRegistry *srclient.SchemaRegistryClient
}

func New(services *Services) *Reader {
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
		msg, err := r.services.Reader.FetchMessage(ctx)
		if err != nil {
			//TODO: provide rich error handling
			r.services.Logger.Sugar().Error(err)
			continue
		}
		err = r.services.Reader.CommitMessages(ctx, msg)
		if err != nil {
			r.services.Logger.Sugar().Error(err)
			continue
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
		m := gomail.NewMessage()
		m.SetHeader("From", email.From)
		m.SetHeader("To", email.To...)
		m.SetAddressHeader("Cc", "info@socialnetwork.ru", "Social Network")
		m.SetHeader("Subject", "Activating new account")
		m.SetBody("text/html", email.Body)
		for _, a := range email.Attachments {
			m.Attach(a)
		}
		err = r.services.Mail.DialAndSend(m)
		if err != nil {
			r.services.Logger.Sugar().Error(err)
		}
	}

}
