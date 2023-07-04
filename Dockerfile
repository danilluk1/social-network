FROM node:18-alpine as builder
COPY --from=golang:1.20.1-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add git curl wget upx protoc libc6-compat && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
  npm i -g pnpm@8 @go-task/cli

COPY go.work go.work.sum docker-entrypoint.sh ./

COPY apps apps
COPY k8s k8s
COPY libs libs
COPY tools tools



### GOLANG MICROSERVICES

FROM alpine:latest as go_prod_base
RUN  apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*
COPY --from=builder /app/docker-entrypoint.sh /app/
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM builder as auth_builder
RUN cd apps/auth && \
  go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as auth
COPY --from=auth_builder /app/apps/auth/out /bin/auth
CMD ["/bin/auth"]