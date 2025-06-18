package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrCategoryNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A category with this name already exists"})
		case services.ErrCategorySlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A category with this slug already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		}
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category"})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	category, err := h.categoryService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category"})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.categoryService.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrCategoryNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		case services.ErrCategoryNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A category with this name already exists"})
		case services.ErrCategorySlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A category with this slug already exists"})
		default:
			log.Printf("Error updating category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = h.categoryService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
