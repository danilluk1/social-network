module github.com/danilluk1/social-network/apps/auth

go 1.20

require (
	github.com/aead/chacha20poly1305 v0.0.0-20201124145622-1a5aba2a8b29
	github.com/bombsimon/logrusr/v3 v3.1.0
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v5 v5.4.1
	github.com/lib/pq v1.10.9
	github.com/o1egl/paseto v1.0.0
	github.com/pressly/goose/v3 v3.13.4
	github.com/riferrei/srclient v0.6.0
	github.com/samber/lo v1.38.1
	github.com/segmentio/kafka-go v0.4.42
	github.com/stretchr/testify v1.8.3
	go.opentelemetry.io/otel v1.16.0
	go.uber.org/zap v1.13.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230629202037-9506855d4529
	google.golang.org/grpc v1.56.1
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/tools v0.10.0 // indirect
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/danilluk1/social-network/libs/avro v0.0.0-00010101000000-000000000000
	github.com/danilluk1/social-network/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/danilluk1/social-network/libs/kafka v0.0.0-00010101000000-000000000000
	github.com/danilluk1/social-network/libs/utils v0.0.0-00010101000000-000000000000
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/linkedin/goavro/v2 v2.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/danilluk1/social-network/libs/avro => ../../libs/avro

replace github.com/danilluk1/social-network/libs/grpc => ../../libs/grpc

replace github.com/danilluk1/social-network/libs/kafka => ../../libs/kafka

replace github.com/danilluk1/social-network/libs/utils => ../../libs/utils
