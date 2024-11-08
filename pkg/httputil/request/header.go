package request

import (
	"net/http"
	"strings"

	"github.com/ncfex/dcart-gateway/pkg/api"
)

const (
	BearerPrefix        string = "Bearer "
	AuthorizationHeader string = "Authorization"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get(AuthorizationHeader)
	if authHeader == "" {
		return "", api.ErrNoAuthHeaderIncluded
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", api.ErrMalformedAuthHeader
	}

	token := authHeader[len(BearerPrefix):]
	if token == "" {
		return "", api.ErrMalformedAuthHeader
	}

	return token, nil
}
