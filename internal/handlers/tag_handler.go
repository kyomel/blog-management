package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/services"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tag, err := h.tagService.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrTagNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A tag with this name already exists"})
		case services.ErrTagSlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A tag with this slug already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		}
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) GetTagByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := h.tagService.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case services.ErrTagNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tag"})
		}
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	tag, err := h.tagService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		switch err {
		case services.ErrTagNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tag"})
		}
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.tagService.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req models.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tag, err := h.tagService.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrTagNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		case services.ErrTagNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A tag with this name already exists"})
		case services.ErrTagSlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A tag with this slug already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		}
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	err = h.tagService.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case services.ErrTagNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TagHandler) GetTagsByPost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	tags, err := h.tagService.GetTagsByPostID(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags for post"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) GetPostsByTag(c *gin.Context) {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.tagService.GetPostsByTagID(c.Request.Context(), tagID, page, pageSize)
	if err != nil {
		switch err {
		case services.ErrTagNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts for tag"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
