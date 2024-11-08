package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/ncfex/dcart-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart-gateway/pkg/api"
	"github.com/ncfex/dcart-gateway/pkg/httputil/request"
	"github.com/ncfex/dcart-gateway/pkg/httputil/response"
)

type AuthMiddleware struct {
	cfg       config.AuthConfig
	responder response.Responder
}

func NewAuthMiddleware(cfg config.AuthConfig, responder response.Responder) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:       cfg,
		responder: responder,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), am.cfg.Timeout)
		defer cancel()

		valid, err := am.validateToken(ctx, r.Header)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				am.responder.RespondWithError(w, http.StatusServiceUnavailable, api.ErrTimeout.Error(), api.ErrTimeout)
			} else {
				am.responder.RespondWithError(w, http.StatusUnauthorized, api.ErrUnauthorized.Error(), api.ErrUnauthorized)
			}
			return
		}
		if !valid {
			am.responder.RespondWithError(w, http.StatusUnauthorized, api.ErrUnauthorized.Error(), api.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) validateToken(ctx context.Context, header http.Header) (bool, error) {
	token, err := request.GetBearerToken(header)
	if err != nil {
		return false, api.ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, am.cfg.ServiceURL, nil)
	if err != nil {
		return false, api.ErrRequestFailed
	}
	req.Header.Set(request.AuthorizationHeader, request.BearerPrefix+token)

	client := &http.Client{
		Timeout: am.cfg.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, api.ErrTimeout
		}
		return false, api.ErrRequestFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, api.ErrUnauthorized
	}

	return true, nil
}
