package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

// FindOrCreateGoogleUser upserts a user with provider='google'.
// Returns ErrProviderMismatch if the email exists with a different provider.
func (r *UserRepo) FindOrCreateGoogleUser(ctx context.Context, email string) (*models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user := &models.User{}
	err = tx.QueryRow(ctx,
		`SELECT id, email, password_hash, provider, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Provider, &user.CreatedAt)

	if err == nil {
		// User found — reject if not a google account
		if user.Provider != "google" {
			return nil, ErrProviderMismatch
		}
	} else if errors.Is(err, pgx.ErrNoRows) {
		// New user
		if err = tx.QueryRow(ctx,
			`INSERT INTO users (email, provider) VALUES ($1, 'google')
			 RETURNING id, email, password_hash, provider, created_at`,
			email,
		).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Provider, &user.CreatedAt); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Assign traveler role (idempotent — safe for both new and existing users)
	if _, err = tx.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id)
		 SELECT $1, role_id FROM roles WHERE nama_role = 'traveler'
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		user.ID,
	); err != nil {
		return nil, err
	}

	return user, tx.Commit(ctx)
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, provider, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Provider, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *UserRepo) CreateEmailUser(ctx context.Context, email, passwordHash, roleName string) (*models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	u := &models.User{}
	if err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, provider) VALUES ($1, $2, 'email')
		 RETURNING id, email, password_hash, provider, created_at`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Provider, &u.CreatedAt); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	if _, err = tx.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id)
		 SELECT $1, role_id FROM roles WHERE nama_role = $2
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		u.ID, roleName,
	); err != nil {
		return nil, err
	}
	return u, tx.Commit(ctx)
}

func (r *UserRepo) HasRole(ctx context.Context, userID int64, roles ...string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_roles ur
		 JOIN roles r ON r.role_id = ur.role_id
		 WHERE ur.user_id = $1 AND r.nama_role = ANY($2)`,
		userID, roles,
	).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) AdminCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_roles ur
		 JOIN roles r ON r.role_id = ur.role_id
		 WHERE r.nama_role = 'admin'`,
	).Scan(&count)
	return count, err
}
