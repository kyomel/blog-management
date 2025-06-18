package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now

	query := `
        INSERT INTO tags (name, slug, color, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		tag.Name,
		tag.Slug,
		tag.Color,
		tag.CreatedAt,
		tag.UpdatedAt,
	).Scan(&tag.ID)

	return err
}

func (r *TagRepository) GetByID(id uuid.UUID) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `
        SELECT id, name, slug, color, created_at, updated_at, deleted_at
        FROM tags
        WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.Color,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&tag.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tag, err
}

func (r *TagRepository) GetByName(name string) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `
        SELECT id, name, slug, color, created_at, updated_at, deleted_at
        FROM tags
        WHERE name = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, name).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.Color,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&tag.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tag, err
}

func (r *TagRepository) GetBySlug(slug string) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `
        SELECT id, name, slug, color, created_at, updated_at, deleted_at
        FROM tags
        WHERE slug = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, slug).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.Color,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&tag.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tag, err
}

func (r *TagRepository) GetAll(limit, offset int) ([]*models.Tag, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM tags WHERE deleted_at IS NULL`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get tags
	query := `
        SELECT id, name, slug, color, created_at, updated_at
        FROM tags
        WHERE deleted_at IS NULL
        ORDER BY name ASC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.Color,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, tag)
	}

	return tags, total, nil
}

func (r *TagRepository) Update(tag *models.Tag) error {
	tag.UpdatedAt = time.Now()

	query := `
        UPDATE tags
        SET name = $2, 
            slug = $3, 
            color = $4,
            updated_at = $5
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		tag.ID,
		tag.Name,
		tag.Slug,
		tag.Color,
		tag.UpdatedAt,
	).Scan(&tag.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("tag not found")
	}

	return err
}

func (r *TagRepository) Delete(id uuid.UUID) error {
	query := `
        UPDATE tags
        SET deleted_at = $2
        WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

// GetTagsByPostID retrieves all tags associated with a specific post
func (r *TagRepository) GetTagsByPostID(postID uuid.UUID) ([]*models.Tag, error) {
	query := `
        SELECT t.id, t.name, t.slug, t.color, t.created_at, t.updated_at
        FROM tags t
        JOIN post_tags pt ON t.id = pt.tag_id
        WHERE pt.post_id = $1 AND t.deleted_at IS NULL
        ORDER BY t.name ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.Color,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// AddTagsToPost associates tags with a post
func (r *TagRepository) AddTagsToPost(postID uuid.UUID, tagIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First, remove existing associations
	deleteQuery := `DELETE FROM post_tags WHERE post_id = $1`
	_, err = tx.Exec(deleteQuery, postID)
	if err != nil {
		return err
	}

	// Then add new associations
	insertQuery := `INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`
	for _, tagID := range tagIDs {
		_, err = tx.Exec(insertQuery, postID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPostsByTagID retrieves all posts associated with a specific tag
func (r *TagRepository) GetPostsByTagID(tagID uuid.UUID, limit, offset int) ([]*models.Post, int, error) {
	// Get total count
	var total int
	countQuery := `
        SELECT COUNT(p.id)
        FROM posts p
        JOIN post_tags pt ON p.id = pt.post_id
        WHERE pt.tag_id = $1 AND p.deleted_at IS NULL`

	if err := r.db.QueryRow(countQuery, tagID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
        SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.excerpt, 
               p.featured_image_url, p.status, p.view_count, p.is_featured, 
               p.published_at, p.created_at, p.updated_at
        FROM posts p
        JOIN post_tags pt ON p.id = pt.post_id
        WHERE pt.tag_id = $1 AND p.deleted_at IS NULL
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, tagID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.CategoryID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.FeaturedImageURL,
			&post.Status,
			&post.ViewCount,
			&post.IsFeatured,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	return posts, total, nil
}
