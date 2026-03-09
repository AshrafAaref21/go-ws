package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AshrafAaref21/go-ws/internal/dto"
	"github.com/AshrafAaref21/go-ws/internal/middlewares"
	"github.com/AshrafAaref21/go-ws/internal/models"
	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func HandleGetPrivate(w http.ResponseWriter, r *http.Request) {
	privateIDStr := r.PathValue("private_id")
	privateID, err := strconv.ParseInt(privateIDStr, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid private_id", nil)
		return
	}

	private, err := models.GetPrivateByID(privateID)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, false, "private conversation not found", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "private conversation retrieved successfully", private)
}

func HandleJoinPrivate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	var req dto.JoinPrivateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil || req.ReceiverId == 0 {
		utils.JSON(w, http.StatusBadRequest, false, "invalid requested data", nil)
		return
	}

	private, err := models.GetPrivateByUsers(userID, req.ReceiverId)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to retrieve private conversation", nil)
		return
	}

	if private != nil {
		utils.JSON(w, http.StatusOK, true, "private conversation retrieved successfully", private)
		return
	}

	private, err = models.CreatePrivate(userID, req.ReceiverId)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to create private conversation", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "private conversation created successfully", private)
}

func HandleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	privates, err := models.GetPrivatesForUser(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to retrieve conversations", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "conversations retrieved successfully", privates)
}

func HandleGetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	privateIDStr := r.PathValue("private_id")
	privateID, err := strconv.ParseInt(privateIDStr, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid private_id", nil)
		return
	}

	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		page, err = strconv.Atoi(p)
		if err != nil || page < 1 {
			utils.JSON(w, http.StatusBadRequest, false, "invalid page number", nil)
			return
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil || limit < 1 || limit > 100 {
			utils.JSON(w, http.StatusBadRequest, false, "invalid limit number (must be between 1 and 100)", nil)
			return
		}
	}

	messages, err := models.GetMessagesByPrivateID(privateID, page, limit+1)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to retrieve messages", nil)
		return
	}

	hasNextPage := false
	if len(messages) > limit {
		hasNextPage = true
		messages = messages[:limit]
	}

	utils.JSON(w, http.StatusOK, true, "messages retrieved successfully", map[string]any{
		"messages":      messages,
		"page":          page,
		"limit":         limit,
		"has_next_page": hasNextPage,
	})
}
