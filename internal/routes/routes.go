package routes

import (
	"net/http"

	"github.com/AshrafAaref21/go-ws/internal/handlers"
	"github.com/AshrafAaref21/go-ws/internal/middlewares"
	"github.com/AshrafAaref21/go-ws/internal/realtime"
)

func RegisterRoutes(hub *realtime.Hub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health-check-http", handlers.HandleHealthCheckHTTP)
	mux.HandleFunc("GET /api/health-check-ws", handlers.HandleHealthCheckWs)

	// Authentication routes
	mux.HandleFunc("POST /api/auth/register-email", handlers.HandleEmailRegistration)
	mux.HandleFunc("POST /api/auth/login-email", handlers.HandleEmailLogin)
	mux.HandleFunc("POST /api/auth/refresh-session", handlers.HandleRefreshSession)
	mux.Handle("POST /api/auth/logout", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleEmailLogout)))
	mux.Handle("POST /api/auth/current-user", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleCurrentUser)))

	// Users routes
	mux.Handle("GET /api/users/{id}", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleGetUserByID)))

	// Conversations routes
	mux.Handle("GET /api/conversations/privates/{private_id}", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleGetPrivate)))
	mux.Handle("POST /api/conversations/privates/create", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleJoinPrivate)))
	mux.Handle("GET /api/conversations", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleGetConversations)))
	mux.Handle("GET /api/conversations/privates/{private_id}/messages", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleGetPrivateMessages)))

	// Files Routes
	mux.Handle("POST /api/files/{private_id}", middlewares.AuthenticateMiddleware(http.HandlerFunc(handlers.HandleFileUpload)))
	mux.Handle("GET /api/files/", middlewares.AuthenticateHandler(handlers.HandleGetFile()))

	// WebSockets Route
	mux.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(hub, w, r)
	})

	middlewareHandlerMux := middlewares.RegisterMiddleWares(mux)

	return middlewareHandlerMux
}
