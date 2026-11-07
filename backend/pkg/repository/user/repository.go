package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/payvue/payvue-backend/pkg/domain/user"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
	)

	if err != nil {
		return nil, user.ErrDatabaseError
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, user.ErrDatabaseError
	}

	u.ID = int(id)
	return u, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, user.ErrDatabaseError
	}

	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, user.ErrDatabaseError
	}

	return &u, nil
}
