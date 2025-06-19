package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyomel/blog-management/internal/services"
	"github.com/kyomel/blog-management/internal/services/cloudinary"
	"github.com/kyomel/blog-management/internal/utils"
)

type UploadHandler struct {
	userService       *services.UserService
	cloudinaryService *cloudinary.CloudinaryService
}

func NewUploadHandler(userService *services.UserService, cloudinaryService *cloudinary.CloudinaryService) *UploadHandler {
	return &UploadHandler{
		userService:       userService,
		cloudinaryService: cloudinaryService,
	}
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// Get user claims from context (set by auth middleware)
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Extract user ID from claims
	userClaims, ok := claims.(*utils.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user claims"})
		return
	}
	
	userID := userClaims.UserID

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files (JPEG, PNG, GIF) are allowed"})
		return
	}

	imageURL, err := h.cloudinaryService.UploadAvatar(c.Request.Context(), file, userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	err = h.userService.UpdateAvatarURL(c.Request.Context(), userID.String(), imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar uploaded successfully",
		"avatar_url": imageURL,
	})
}
