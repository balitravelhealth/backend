package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const accessTokenTTL  = 15 * time.Minute
const refreshTokenTTL = 30 * 24 * time.Hour

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthService struct {
	userRepo  *repository.UserRepo
	tokenRepo *repository.TokenRepo
}

func NewAuthService(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *AuthService) GoogleLogin(ctx context.Context, idToken, deviceInfo string) (*TokenPair, *models.User, error) {
	payload, err := validateGoogleIDToken(ctx, idToken)
	if err != nil {
		log.Printf("google token validation failed: %v", err)
		return nil, nil, ErrInvalidCredentials
	}

	email, _ := payload.Claims["email"].(string)
	if email == "" {
		return nil, nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindOrCreateGoogleUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrProviderMismatch) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("upsert user: %w", err)
	}

	pair, err := s.generateTokenPair(ctx, user.ID, deviceInfo)
	if err != nil {
		return nil, nil, err
	}

	return pair, user, nil
}

// RefreshTokens validates the raw refresh token, rotates it, and returns a new pair.
// On reuse of a revoked token, all sessions for the user are invalidated (theft detection).
func (s *AuthService) RefreshTokens(ctx context.Context, rawToken, deviceInfo string) (*TokenPair, error) {
	raw, hash, err := decodeAndHash(rawToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	_ = raw

	stored, err := s.tokenRepo.FindByHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if stored.Revoked {
		// Token reuse detected — revoke all sessions for this user
		_ = s.tokenRepo.RevokeAllForUser(ctx, stored.UserID)
		return nil, ErrInvalidCredentials
	}
	if time.Now().After(stored.ExpiredAt) {
		return nil, ErrInvalidCredentials
	}

	if err := s.tokenRepo.Revoke(ctx, hash); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	return s.generateTokenPair(ctx, stored.UserID, deviceInfo)
}

// Logout revokes the given refresh token (single-device logout).
func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	_, hash, err := decodeAndHash(rawToken)
	if err != nil {
		return ErrInvalidCredentials
	}

	if err := s.tokenRepo.Revoke(ctx, hash); err != nil {
		if errors.Is(err, repository.ErrTokenNotFound) {
			return ErrInvalidCredentials
		}
		return err
	}
	return nil
}

// decodeAndHash decodes a base64url refresh token and returns its SHA-256 hex hash.
func decodeAndHash(rawToken string) ([]byte, string, error) {
	raw, err := base64.URLEncoding.DecodeString(rawToken)
	if err != nil {
		return nil, "", err
	}
	sum := sha256.Sum256(raw)
	return raw, hex.EncodeToString(sum[:]), nil
}

func (s *AuthService) generateTokenPair(ctx context.Context, userID int64, deviceInfo string) (*TokenPair, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(accessTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// 32 random bytes → base64url for client, SHA-256 hash for DB storage
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	rawStr := base64.URLEncoding.EncodeToString(raw)
	sum := sha256.Sum256(raw)
	hashStr := hex.EncodeToString(sum[:])

	if err := s.tokenRepo.Create(ctx, userID, hashStr, deviceInfo, time.Now().Add(refreshTokenTTL)); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawStr,
		ExpiresIn:    int(accessTokenTTL.Seconds()),
	}, nil
}
func validateGoogleIDToken(ctx context.Context, idToken string) (*idtoken.Payload, error) {
	audiences := strings.Split(os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"), ",")

	for _, audience := range audiences {
		audience = strings.TrimSpace(audience)
		if audience == "" {
			continue
		}

		payload, err := idtoken.Validate(ctx, idToken, audience)
		if err == nil {
			return payload, nil
		}
	}

	return nil, ErrInvalidCredentials
}
