package fxapp

import "github.com/danilluk1/social-network/libs/go/pkg/logger"

type FxErrorHandler struct {
	logger logger.Logger
}

func NewFxErrorHandler(logger logger.Logger) *FxErrorHandler {
	return &FxErrorHandler{logger: logger}
}

func (h *FxErrorHandler) HandleError(e error) {
	h.logger.Error(e)
}
