package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/user/models"
	"github.com/johnroshan2255/core-service/internal/user/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(userService *service.Service) *Handler {
	return &Handler{
		service: userService,
	}
}

func (h *Handler) GetProfile(c *gin.Context) {
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

	profile, err := h.service.GetProfile(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
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

	var req struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Username    string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.PhoneNumber != "" {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}

	if err := h.service.UpdateProfile(c.Request.Context(), uuid, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})
}

func (h *Handler) GetCompanyDetails(c *gin.Context) {
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

	company, err := h.service.GetCompanyDetails(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    company,
	})
}

func (h *Handler) UpdateCompanyDetails(c *gin.Context) {
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

	var company models.CompanyDetails
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateCompanyDetails(c.Request.Context(), uuid, &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Company details updated successfully",
	})
}

func (h *Handler) GetPaymentDetails(c *gin.Context) {
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

	payment, err := h.service.GetPaymentDetails(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

func (h *Handler) UpdatePaymentDetails(c *gin.Context) {
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

	var payment models.PaymentDetails
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePaymentDetails(c.Request.Context(), uuid, &payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment details updated successfully",
	})
}

func (h *Handler) GetPaymentHistory(c *gin.Context) {
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

	history, err := h.service.GetPaymentHistory(c.Request.Context(), uuid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
	})
}

func (h *Handler) CreatePaymentHistory(c *gin.Context) {
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

	var payment models.PaymentHistory
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePaymentHistory(c.Request.Context(), uuid, &payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment history created successfully",
	})
}
