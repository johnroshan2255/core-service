package document

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/document/models"
	"github.com/johnroshan2255/core-service/internal/document/service"
)

type Handler struct {
	service *service.Service
	uploadPath string
}

func NewHandler(documentService *service.Service, uploadPath string) *Handler {
	if uploadPath == "" {
		uploadPath = "./uploads/documents"
	}
	
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		fmt.Printf("Warning: Failed to create upload directory: %v\n", err)
	}
	
	return &Handler{
		service:    documentService,
		uploadPath: uploadPath,
	}
}

func (h *Handler) UploadDocument(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found in token"})
		return
	}

	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID format"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	var req struct {
		Name        string `form:"name"`
		Description string `form:"description"`
		Category    string `form:"category"`
		IssueDate   string `form:"issue_date"`
		ExpiryDate  string `form:"expiry_date"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		req.Name = file.Filename
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if err := h.service.ValidateDocument(file.Filename, mimeType, file.Size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docType := h.service.DetermineDocumentType(mimeType)

	fileName := fmt.Sprintf("%s_%d_%s", uuid, time.Now().Unix(), file.Filename)
	filePath := filepath.Join(h.uploadPath, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	doc := &models.Document{
		Name:        req.Name,
		Description: req.Description,
		Category:    models.DocumentCategory(req.Category),
		Type:        docType,
		FileName:    file.Filename,
		FilePath:    filePath,
		FileSize:    file.Size,
		MimeType:    mimeType,
	}

	if req.IssueDate != "" {
		issueDate, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			doc.IssueDate = &issueDate
		}
	}

	if req.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			doc.ExpiryDate = &expiryDate
		}
	}

	if doc.Category == "" {
		doc.Category = models.DocumentCategoryOther
	}

	createdDoc, err := h.service.CreateDocument(c.Request.Context(), uuid, doc)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    createdDoc,
	})
}

func (h *Handler) GetDocument(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found in token"})
		return
	}

	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID format"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), uuid, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    doc,
	})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found in token"})
		return
	}

	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID format"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	docs, err := h.service.GetUserDocuments(c.Request.Context(), uuid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    docs,
	})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found in token"})
		return
	}

	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID format"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ExpiryDate  string `json:"expiry_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			updates["expiry_date"] = &expiryDate
		}
	}

	if err := h.service.UpdateDocument(c.Request.Context(), uuid, uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document updated successfully",
	})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found in token"})
		return
	}

	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID format"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), uuid, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), uuid, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.Remove(doc.FilePath); err != nil {
		log.Printf("Failed to delete file: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document deleted successfully",
	})
}
