// internal/repositories/category_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kyomel/blog-management/internal/models"

	"github.com/google/uuid"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	query := `
        INSERT INTO categories (name, slug, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)

	return err
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	category := &models.Category{}
	query := `
        SELECT id, name, slug, description, created_at, updated_at, deleted_at
        FROM categories
        WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return category, err
}

func (r *CategoryRepository) GetByName(name string) (*models.Category, error) {
	category := &models.Category{}
	query := `
        SELECT id, name, slug, description, created_at, updated_at, deleted_at
        FROM categories
        WHERE name = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, name).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return category, err
}

func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	category := &models.Category{}
	query := `
        SELECT id, name, slug, description, created_at, updated_at, deleted_at
        FROM categories
        WHERE slug = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return category, err
}

func (r *CategoryRepository) GetAll(limit, offset int) ([]*models.Category, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM categories WHERE deleted_at IS NULL`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get categories
	query := `
        SELECT id, name, slug, description, created_at, updated_at
        FROM categories
        WHERE deleted_at IS NULL
        ORDER BY name ASC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	category.UpdatedAt = time.Now()

	query := `
        UPDATE categories
        SET name = $2, 
            slug = $3, 
            description = $4,
            updated_at = $5
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.UpdatedAt,
	).Scan(&category.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("category not found")
	}

	return err
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	query := `
        UPDATE categories
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
		return fmt.Errorf("category not found")
	}

	return nil
}
