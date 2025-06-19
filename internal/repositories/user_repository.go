package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateAvatarURL(ctx context.Context, userID uuid.UUID, avatarURL string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUsernameAlreadyExists
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, username, fullname, password_hash, role, avatar_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.Fullname,
		user.PasswordHash,
		user.Role,
		user.AvatarURL,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, username, fullname, password_hash, role, avatar_url, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	return r.findOneByQuery(ctx, query, id)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, fullname, password_hash, role, avatar_url, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	return r.findOneByQuery(ctx, query, email)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, fullname, password_hash, role, avatar_url, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`
	return r.findOneByQuery(ctx, query, username)
}

func (r *userRepository) findOneByQuery(ctx context.Context, query string, args ...interface{}) (*models.User, error) {
	var user models.User
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Fullname,
		&user.PasswordHash,
		&user.Role,
		&user.AvatarURL,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, username = $2, fullname = $3, password_hash = $4, role = $5, avatar_url = $6, is_active = $7, updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Username,
		user.Fullname,
		user.PasswordHash,
		user.Role,
		user.AvatarURL,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateAvatarURL updates the avatar URL for a user
func (r *userRepository) UpdateAvatarURL(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	query := `
		UPDATE users
		SET avatar_url = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
