package api

import (
	"context"
	"net/http"
	"regexp"

	"github.com/danilluk1/social-network/apps/gateway/internal/conf"
	"github.com/sebest/xff"
)

const (
	defaultVersion = "unknown version"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

type API struct {
	handler  http.Handler
	version  string
	services *Services
}

type Services struct {
	config *conf.GlobalConfiguration
}

func NewAPI(globalConfig *conf.GlobalConfiguration, services *Services) *API {
	return NewAPIWithVersion(context.Background(), globalConfig, services, defaultVersion)
}

func NewAPIWithVersion(ctx context.Context, globalConfig *conf.GlobalConfiguration, services *Services, version string) *API {
	api := &API{services: services, version: version}

	xffmw, _ := xff.Default()

	r := new
}
