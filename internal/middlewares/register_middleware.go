package middlewares

import (
	"net/http"
)

func RegisterMiddleWares(mux *http.ServeMux) http.Handler {
	middlewareMux := LoggingMiddleware(mux)
	middlewareMux = CorsMiddleware(middlewareMux)

	return middlewareMux
}
