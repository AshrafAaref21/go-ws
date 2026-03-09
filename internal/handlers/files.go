package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AshrafAaref21/go-ws/internal/middlewares"
	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	senderID, ok := r.Context().Value(middlewares.CtxUserID).(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	privateIDStr := r.PathValue("private_id")
	privateID, err := strconv.ParseInt(privateIDStr, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid private_id", nil)
		return
	}

	err = r.ParseMultipartForm(50 << 20)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "failed to parse multipart form", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "failed to retrieve file from form", nil)
		return
	}
	defer file.Close()

	dirPath := filepath.Join("files", "chats", fmt.Sprintf("%d", privateID), fmt.Sprintf("%d", senderID))

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to create directory", nil)
		return
	}

	filePath := filepath.Join(dirPath, header.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to create file", nil)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to save file", nil)
		return
	}

	fileUrl := fmt.Sprintf("/files/chats/%d/%d/%s", privateID, senderID, header.Filename)

	utils.JSON(w, http.StatusOK, true, "file uploaded successfully", fileUrl)
}

func HandleGetFile() http.Handler {
	fs := http.FileServer(http.Dir("./files"))

	return http.StripPrefix("/api/files", fs)
}
