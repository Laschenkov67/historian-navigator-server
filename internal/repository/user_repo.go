package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/historian/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound raised when entity is missing.
var ErrNotFound = errors.New("not found")

// UserRepo wraps user DB operations.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo constructor.
func NewUserRepo(p *pgxpool.Pool) *UserRepo { return &UserRepo{pool: p} }

// Create inserts new user.
func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users(id,email,password,full_name,role,created_at)
		 VALUES($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Email, u.Password, u.FullName, u.Role, u.CreatedAt)
	return err
}

// GetByEmail fetches user by email.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id,email,password,full_name,role,created_at FROM users WHERE email=$1`, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

// GetByID fetches user by id.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id,email,password,full_name,role,created_at FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}
