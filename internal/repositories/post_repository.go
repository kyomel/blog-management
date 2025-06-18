// internal/repositories/post_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kyomel/blog-management/internal/models"

	"github.com/google/uuid"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post, tagIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO posts (author_id, category_id, title, slug, content, excerpt, 
                           featured_image_url, status, is_featured, metadata, published_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at`

	err = tx.QueryRow(
		query,
		post.AuthorID,
		post.CategoryID,
		post.Title,
		post.Slug,
		post.Content,
		post.Excerpt,
		post.FeaturedImageURL,
		post.Status,
		post.IsFeatured,
		post.Metadata,
		post.PublishedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	if len(tagIDs) > 0 {
		tagQuery := `INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`
		for _, tagID := range tagIDs {
			if _, err := tx.Exec(tagQuery, post.ID, tagID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostRepository) GetByID(id uuid.UUID) (*models.Post, error) {
	post := &models.Post{
		Author:   &models.User{},
		Category: &models.Category{},
	}
	var metadataJSON []byte

	query := `
        SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.content, 
               p.excerpt, p.featured_image_url, p.status, p.view_count, 
               p.is_featured, p.metadata, p.published_at, p.created_at, 
               p.updated_at, p.deleted_at,
               u.username, u.full_name, u.avatar_url,
               c.name, c.slug
        FROM posts p
        JOIN users u ON p.author_id = u.id
        JOIN categories c ON p.category_id = c.id
        WHERE p.id = $1 AND p.deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.CategoryID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Excerpt,
		&post.FeaturedImageURL,
		&post.Status,
		&post.ViewCount,
		&post.IsFeatured,
		&metadataJSON,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.DeletedAt,
		&post.Author.Username,
		&post.Author.Fullname,
		&post.Author.AvatarURL,
		&post.Category.Name,
		&post.Category.Slug,
	)

	if metadataJSON != nil {
		post.Metadata = metadataJSON
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	tagsErr := error(nil)
	post.Tags, tagsErr = r.getPostTags(post.ID)
	if tagsErr != nil {
		fmt.Printf("Error fetching tags for post %s: %v\n", post.ID, tagsErr)
	}

	return post, nil
}

func (r *PostRepository) GetBySlug(slug string) (*models.Post, error) {
	post := &models.Post{
		Author:   &models.User{},
		Category: &models.Category{},
	}
	var metadataJSON []byte

	query := `
        SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.content, 
               p.excerpt, p.featured_image_url, p.status, p.view_count, 
               p.is_featured, p.metadata, p.published_at, p.created_at, 
               p.updated_at, p.deleted_at,
               u.username, u.full_name, u.avatar_url,
               c.name, c.slug
        FROM posts p
        JOIN users u ON p.author_id = u.id
        JOIN categories c ON p.category_id = c.id
        WHERE p.slug = $1 AND p.deleted_at IS NULL`

	err := r.db.QueryRow(query, slug).Scan(
		&post.ID,
		&post.AuthorID,
		&post.CategoryID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Excerpt,
		&post.FeaturedImageURL,
		&post.Status,
		&post.ViewCount,
		&post.IsFeatured,
		&metadataJSON,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.DeletedAt,
		&post.Author.Username,
		&post.Author.Fullname,
		&post.Author.AvatarURL,
		&post.Category.Name,
		&post.Category.Slug,
	)

	if metadataJSON != nil {
		post.Metadata = metadataJSON
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	tagsErr := error(nil)
	post.Tags, tagsErr = r.getPostTags(post.ID)
	if tagsErr != nil {
		fmt.Printf("Error fetching tags for post %s: %v\n", post.ID, tagsErr)
	}

	return post, nil
}

func (r *PostRepository) GetAll(filter *models.PostFilter) ([]*models.Post, int, error) {
	whereConditions := []string{"p.deleted_at IS NULL"}
	args := []interface{}{}
	argCount := 0

	if filter.Status != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.status = $%d", argCount))
		args = append(args, filter.Status)
	}

	if filter.CategoryID != nil && *filter.CategoryID != uuid.Nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.category_id = $%d", argCount))
		args = append(args, filter.CategoryID)
	}

	if filter.AuthorID != nil && *filter.AuthorID != uuid.Nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.author_id = $%d", argCount))
		args = append(args, filter.AuthorID)
	}

	if filter.IsFeatured != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.is_featured = $%d", argCount))
		args = append(args, *filter.IsFeatured)
	}

	if filter.Search != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("(p.title ILIKE $%d OR p.content ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+filter.Search+"%")
	}

	whereClause := strings.Join(whereConditions, " AND ")

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM posts p WHERE %s`, whereClause)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	argCount++
	args = append(args, filter.Limit)
	argCount++
	args = append(args, filter.Offset)

	query := fmt.Sprintf(`
        SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.excerpt, 
               p.featured_image_url, p.status, p.view_count, p.is_featured, 
               p.metadata, p.published_at, p.created_at, p.updated_at,
               u.username, u.full_name, u.avatar_url,
               c.name, c.slug
        FROM posts p
        JOIN users u ON p.author_id = u.id
        JOIN categories c ON p.category_id = c.id
        WHERE %s
        ORDER BY p.created_at DESC
        LIMIT $%d OFFSET $%d`, whereClause, argCount-1, argCount)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{
			Author:   &models.User{},
			Category: &models.Category{},
		}
		var metadataJSON []byte
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
			&metadataJSON,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author.Username,
			&post.Author.Fullname,
			&post.Author.AvatarURL,
			&post.Category.Name,
			&post.Category.Slug,
		)

		if metadataJSON != nil {
			post.Metadata = metadataJSON
		}
		if err != nil {
			return nil, 0, err
		}

		var tagsErr error
		post.Tags, tagsErr = r.getPostTags(post.ID)
		if tagsErr != nil {
			fmt.Printf("Error fetching tags for post %s: %v\n", post.ID, tagsErr)
		}
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) Update(post *models.Post, tagIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	post.UpdatedAt = time.Now()

	query := `
        UPDATE posts
        SET category_id = $2, title = $3, slug = $4, content = $5, excerpt = $6,
            featured_image_url = $7, status = $8, is_featured = $9, metadata = $10,
            updated_at = $11
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at`

	err = tx.QueryRow(
		query,
		post.ID,
		post.CategoryID,
		post.Title,
		post.Slug,
		post.Content,
		post.Excerpt,
		post.FeaturedImageURL,
		post.Status,
		post.IsFeatured,
		post.Metadata,
		post.UpdatedAt,
	).Scan(&post.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}

	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM post_tags WHERE post_id = $1`, post.ID); err != nil {
		return err
	}

	if len(tagIDs) > 0 {
		tagQuery := `INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`
		for _, tagID := range tagIDs {
			if _, err := tx.Exec(tagQuery, post.ID, tagID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostRepository) Delete(id uuid.UUID) error {
	query := `
        UPDATE posts
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
		return fmt.Errorf("post not found")
	}

	return nil
}

func (r *PostRepository) IncrementViewCount(id uuid.UUID) error {
	query := `UPDATE posts SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostRepository) Publish(id uuid.UUID) error {
	query := `
        UPDATE posts
        SET status = 'published', published_at = $2
        WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, id, time.Now())
	return err
}

func (r *PostRepository) getPostTags(postID uuid.UUID) ([]*models.Tag, error) {
	query := `
        SELECT t.id, t.name, t.slug, t.color
        FROM tags t
        JOIN post_tags pt ON t.id = pt.tag_id
        WHERE pt.post_id = $1`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Color)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
