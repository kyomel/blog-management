package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/services"
)

type PostHandler struct {
	postService services.PostService
}

func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Log the request for debugging
	fmt.Printf("Creating post with request: %+v\n", req)

	post, err := h.postService.Create(c.Request.Context(), &req)
	if err != nil {
		// Log the detailed error
		fmt.Printf("Error creating post: %v\n", err)

		switch err {
		case services.ErrPostSlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A post with this slug already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == services.ErrPostNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
		}
		return
	}

	// Increment view count asynchronously
	go func() {
		_ = h.postService.IncrementViewCount(c.Request.Context(), id)
	}()

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	post, err := h.postService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		if err == services.ErrPostNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
		}
		return
	}

	// Increment view count asynchronously if post is found
	if post != nil {
		go func() {
			_ = h.postService.IncrementViewCount(c.Request.Context(), post.ID)
		}()
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	filter := &models.PostFilter{}

	if status := c.Query("status"); status != "" {
		filter.Status = models.PostStatus(status)
	} else {
		filter.Status = models.StatusPublished
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err == nil {
			filter.CategoryID = &id
		}
	}

	if authorID := c.Query("author_id"); authorID != "" {
		id, err := uuid.Parse(authorID)
		if err == nil {
			filter.AuthorID = &id
		}
	}

	if featured := c.Query("featured"); featured == "true" {
		isFeatured := true
		filter.IsFeatured = &isFeatured
	}

	filter.Search = c.Query("search")

	result, err := h.postService.GetAll(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post, err := h.postService.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrPostNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case services.ErrPostSlugConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "A post with this slug already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	err = h.postService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == services.ErrPostNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) PublishPost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postService.Publish(c.Request.Context(), id)
	if err != nil {
		if err == services.ErrPostNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish post"})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}
