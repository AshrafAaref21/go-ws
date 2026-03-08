package handlers

import (
	"net/http"
	"strconv"

	"github.com/AshrafAaref21/go-ws/internal/models"
	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	targetId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid user id", nil)
		return
	}

	user, err := models.GetUserByID(targetId)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, false, "user not found", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "user retrieved successfully", user)
}
