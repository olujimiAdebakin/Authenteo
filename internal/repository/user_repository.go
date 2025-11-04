package repository

import (
	"context"
	"authentio/internal/models"
)

type UserRepository interface {
	// FindByEmail finds a user by email address
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id int64) (*models.User, error)
	
	// Create inserts a new user into the database
	Create(ctx context.Context, user *models.User) error
	
	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error
	
	// Delete soft deletes a user
	Delete(ctx context.Context, id int64) error
}