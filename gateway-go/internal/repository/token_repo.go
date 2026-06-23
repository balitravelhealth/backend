package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type TokenRepo struct {
	db *pgxpool.Pool
}

func NewTokenRepo(db *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Create(ctx context.Context, userID int64, tokenHash, deviceInfo string, expiredAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, device_info, expired_at)
		 VALUES ($1, $2, $3, $4)`,
		userID, tokenHash, deviceInfo, expiredAt,
	)
	return err
}

func (r *TokenRepo) FindByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	t := &models.RefreshToken{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, token_hash, device_info, expired_at, revoked, created_at
		 FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.DeviceInfo, &t.ExpiredAt, &t.Revoked, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTokenNotFound
	}
	return t, err
}

func (r *TokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTokenNotFound
	}
	return nil
}

func (r *TokenRepo) RevokeAllForUser(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true WHERE user_id = $1 AND revoked = false`,
		userID,
	)
	return err
}
