package mail

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"github.com/danilluk1/social-network/libs/avro"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

func (r *Reader) Start(ctx context.Context, conf *conf.Configuration) {
	newCtx, span := otel.Tracer("mail").Start(ctx, "Start")
	defer span.End()

	for {

		msg, err := r.services.Reader.FetchMessage(newCtx)
		if err != nil {
			span.RecordError(err)
			continue
		}
		r.services.Logger.Sugar().Infof("Message %s received", msg.Value)
		span.SetStatus(codes.Ok, fmt.Sprintf("Message %s received", msg.Value))
		err = r.services.Reader.CommitMessages(newCtx, msg)
		if err != nil {
			r.services.Logger.Sugar().Error(err)
			span.RecordError(err)
			continue
		}
		if len(msg.Value) <= 5 {
			continue
		}

		schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
		schema, err := r.services.SchemaRegistry.GetSchema(int(schemaID))
		if err != nil {
			span.RecordError(err)
			r.services.Logger.Sugar().Error(err)
			continue
		}

		email := &EmailMessage{}
		avro.Decode(msg.Value, schema.Codec(), email)
		m := gomail.NewMessage()
		m.SetHeader("From", email.From)
		m.SetHeader("To", email.To...)
		m.SetAddressHeader("Cc", conf.SMTP.AdminEmail, conf.SMTP.SenderName)
		m.SetHeader("Subject", "Activating new account")
		m.SetBody("text/html", email.Body)
		for _, a := range email.Attachments {
			m.Attach(a)
		}
		err = r.services.Mail.DialAndSend(m)
		if err != nil {
			span.RecordError(err)
			r.services.Logger.Sugar().Error(err)
		}
	}

}
