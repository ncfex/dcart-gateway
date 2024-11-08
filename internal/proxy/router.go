package proxy

import (
	"net/http"
	"strings"
	"sync"

	"github.com/ncfex/dcart-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart-gateway/internal/middleware"
	"github.com/ncfex/dcart-gateway/pkg/api"
	"github.com/ncfex/dcart-gateway/pkg/httputil/response"
)

type Router struct {
	cfg       *config.Config
	services  map[string]*serviceProxy
	auth      *middleware.AuthMiddleware
	responder response.Responder
	mu        sync.RWMutex
}

func NewRouter(cfg *config.Config, responder response.Responder) (*Router, error) {
	r := &Router{
		cfg:       cfg,
		services:  make(map[string]*serviceProxy),
		auth:      middleware.NewAuthMiddleware(cfg.Auth, responder),
		responder: responder,
	}

	if err := r.initializeServices(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(req.URL.Path, "/"), "/", 2)
	if len(parts) == 0 {
		r.responder.RespondWithError(w, http.StatusBadRequest, api.ErrInvalidPath.Error(), api.ErrInvalidPath)
		return
	}

	serviceName := parts[0]

	r.mu.RLock()
	service, exists := r.services[serviceName]
	r.mu.RUnlock()

	if !exists {
		r.responder.RespondWithError(w, http.StatusNotFound, api.ErrServiceNotFound.Error(), api.ErrServiceNotFound)
		return
	}

	if len(parts) > 1 {
		req.URL.Path = "/" + parts[1]
	} else {
		req.URL.Path = "/"
	}

	service.ServeHTTP(w, req)
}

func (r *Router) initializeServices() error {
	for _, svc := range r.cfg.Services {
		proxy, err := newServiceProxy(&svc, r)
		if err != nil {
			return err
		}

		r.mu.Lock()
		r.services[svc.Name] = proxy
		r.mu.Unlock()
	}
	return nil
}
