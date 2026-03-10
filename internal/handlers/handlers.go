package handlers

import (
	"net/http"

	"github.com/AshrafAaref21/go-ws/internal/utils"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func HandleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, true, "API is runnung", nil)
}

func HandleHealthCheckWs(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	ctx := r.Context()

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			break
		}

		response := map[string]any{
			"data":    string(message),
			"from":    "server",
			"success": true,
		}

		err = wsjson.Write(ctx, conn, response)
		if err != nil {
			break
		}
	}

}
