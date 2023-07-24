FROM golang:1.20-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git && \
  go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /go/src/github.com/danilluk1/social-network/apps/auth

COPY ./Taskfile.yaml ./go.* ./
RUN task deps

COPY . /go/src/github.com/danilluk1/social-network/apps/auth
RUN task deps

FROM alpine:3.17
RUN adduser -D -u 1000 socialy

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/danilluk1/social-network/apps/auth /usr/local/bin/auth
COPY --from=build /go/src/github.com/danilluk1/social-network/apps/auth/internal/db/migration /usr/local/etc/auth/migration

ENV MIGRATION_PATH=/usr/local/etc/auth/migration

USER socialy
CMD ["auth"]
