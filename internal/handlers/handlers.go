package handlers

import (
	"net/http"

	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func HandleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, true, "API is runnung", nil)
}
