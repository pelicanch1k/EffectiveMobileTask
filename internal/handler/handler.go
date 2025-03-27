package handler

import (
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/service"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"

	_ "github.com/pelicanch1k/EffectiveMobileTestTask/docs"
)

type Handler struct {
	services *service.Service
	logger   *logging.Logger
}

func NewHandler(services *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}