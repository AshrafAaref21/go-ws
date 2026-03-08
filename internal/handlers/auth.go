package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AshrafAaref21/go-ws/internal/dto"
	"github.com/AshrafAaref21/go-ws/internal/middlewares"
	"github.com/AshrafAaref21/go-ws/internal/models"
	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func HandleEmailRegistration(w http.ResponseWriter, r *http.Request) {
	var req dto.EmailRegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid requested data", nil)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, false, "name, email and password are required", nil)
		return
	}

	existingUser, _ := models.GetUserByEmail(req.Email)

	if existingUser != nil {
		utils.JSON(w, http.StatusConflict, false, "email already exists", nil)
		return
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to hash password", nil)
		return
	}

	user, err := models.CreateUserByEmail(req.Email, req.Name, hashedPass)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to create user", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "user created successfully", user)

}

func HandleEmailLogin(w http.ResponseWriter, r *http.Request) {
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(middlewares.CtxPlatform)))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req dto.EmailLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid requested data", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, false, "email and password are required", nil)
		return
	}

	user, err := models.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "invalid email or password", nil)
		return
	}

	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "invalid email or password", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Name, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to generate token", nil)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to generate refresh token", nil)
		return
	}

	err = models.UpdateUserRefreshToken(user.ID, platform, refreshToken)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to save refresh token", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "login successful", map[string]any{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func HandleEmailLogout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.CtxUserID).(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "unauthorized", nil)
		return
	}

	platform, ok := r.Context().Value(middlewares.CtxPlatform).(string)
	if !ok {
		utils.JSON(w, http.StatusBadRequest, false, "invalid platform", nil)
		return
	}

	err := models.DeleteUserRefreshToken(userId, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to logout", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "logout successful", nil)
}

func HandleRefreshSession(w http.ResponseWriter, r *http.Request) {

	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(middlewares.CtxPlatform)))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req dto.RefreshSessionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid requested data", nil)
		return
	}

	if req.RefreshToken == "" {
		utils.JSON(w, http.StatusBadRequest, false, "refresh token is required", nil)
		return
	}

	user, err := models.GetUserByRefreshToken(req.RefreshToken, platform)
	if err != nil || user == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "invalid refresh token", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Name, platform)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to generate token", nil)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to generate refresh token", nil)
		return
	}

	err = models.UpdateUserRefreshToken(user.ID, platform, refreshToken)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to save refresh token", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "session refreshed successfully", map[string]any{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

func HandleCurrentUser(w http.ResponseWriter, r *http.Request) {
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get(middlewares.CtxPlatform)))
	if platform != middlewares.PlatformWeb && platform != middlewares.PlatformMobile {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid platform", nil)
		return
	}

	var req dto.RefreshSessionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid requested data", nil)
		return
	}

	if req.RefreshToken == "" {
		utils.JSON(w, http.StatusBadRequest, false, "refresh token is required", nil)
		return
	}

	user, err := models.GetUserByRefreshToken(req.RefreshToken, platform)
	if err != nil || user == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "invalid refresh token", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "current user retrieved successfully", map[string]any{
		"user": user,
	})
}
