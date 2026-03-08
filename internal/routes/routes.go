package routes

import (
	"net/http"

	"github.com/AshrafAaref21/go-ws/internal/handlers"
	"github.com/AshrafAaref21/go-ws/internal/middlewares"
)

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health-check-http", handlers.HandleHealthCheckHTTP)

	// Authentication routes
	mux.HandleFunc("POST /api/auth/register-email", handlers.HandleEmailRegistration)
	mux.HandleFunc("POST /api/auth/login-email", handlers.HandleEmailLogin)
	mux.HandleFunc("POST /api/auth/refresh-session", handlers.HandleRefreshSession)
	mux.Handle("POST /api/auth/logout", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleEmailLogout)))
	mux.Handle("GET /api/auth/current-user", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleCurrentUser)))

	// Users routes
	mux.Handle("GET /api/users/{id}", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleGetUserByID)))

	middlewareHandlerMux := registerMiddleWares(mux)

	return middlewareHandlerMux
}

func registerMiddleWares(mux *http.ServeMux) http.Handler {
	middlewareMux := middlewares.LoggingMiddleware(mux)
	middlewareMux = middlewares.CorsMiddleware(middlewareMux)

	return middlewareMux
}
