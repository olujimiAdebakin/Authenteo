package database



import (
	"context"
	"database/sql"
	"authentio/internal/models"
	"authentio/internal/repository"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_active, created_at, updated_at 
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_active, created_at, updated_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	
	err := r.db.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	
	return err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, is_active = $4, updated_at = $5
		WHERE id = $6`
	
	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}