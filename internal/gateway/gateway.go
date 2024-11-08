package gateway

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart-gateway/internal/proxy"
	"github.com/ncfex/dcart-gateway/pkg/httputil/response"
)

type Gateway struct {
	HttpServer *http.Server
	cfg        *config.Config
	router     *proxy.Router

	log       *log.Logger
	responder response.Responder
}

func NewGateway(
	cfg *config.Config,
	log *log.Logger,
	responder response.Responder,
) (*Gateway, error) {
	router, err := proxy.NewRouter(cfg, responder)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	return &Gateway{
		HttpServer: srv,
		cfg:        cfg,
		router:     router,

		log:       log,
		responder: responder,
	}, nil
}
