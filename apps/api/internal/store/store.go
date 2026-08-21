package store

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrEmailExists      = errors.New("email already exists")
	ErrCooldown         = errors.New("challenge cooldown is active")
	ErrInvalidCode      = errors.New("invalid verification code")
	ErrChallengeExpired = errors.New("verification challenge expired")
	ErrTooManyAttempts  = errors.New("too many verification attempts")
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func NewID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		value[0:4], value[4:6], value[6:8], value[8:10], value[10:16]), nil
}

func (s *Store) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (s *Store) CreateChallenge(ctx context.Context, challenge Challenge, cooldown time.Duration) error {
	var lastCreated time.Time
	err := s.db.QueryRowContext(ctx, `
		SELECT created_at
		FROM email_challenges
		WHERE email = $1 AND purpose = $2
		ORDER BY created_at DESC
		LIMIT 1`, challenge.Email, challenge.Purpose).Scan(&lastCreated)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err == nil && time.Since(lastCreated) < cooldown {
		return ErrCooldown
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO email_challenges (
			id, email, purpose, pending_name, code_hash, max_attempts, expires_at
		) VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7)`,
		challenge.ID, challenge.Email, challenge.Purpose, challenge.PendingName,
		challenge.CodeHash, challenge.MaxAttempts, challenge.ExpiresAt)
	return err
}

func (s *Store) InvalidateChallenge(ctx context.Context, id string) {
	_, _ = s.db.ExecContext(ctx, `
		UPDATE email_challenges SET consumed_at = now()
		WHERE id = $1 AND consumed_at IS NULL`, id)
}

func (s *Store) CompleteChallenge(
	ctx context.Context,
	challengeID string,
	email string,
	purpose string,
	candidateHash string,
	sessionID string,
	sessionHash string,
	sessionExpiresAt time.Time,
	registrationRole string,
) (User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	var storedHash string
	var attempts, maxAttempts int
	var expiresAt time.Time
	var consumedAt sql.NullTime
	var pendingName sql.NullString
	err = tx.QueryRowContext(ctx, `
		SELECT code_hash, attempts, max_attempts, expires_at, consumed_at, pending_name
		FROM email_challenges
		WHERE id = $1 AND email = $2 AND purpose = $3
		FOR UPDATE`, challengeID, email, purpose).
		Scan(&storedHash, &attempts, &maxAttempts, &expiresAt, &consumedAt, &pendingName)
	if errors.Is(err, sql.ErrNoRows) || consumedAt.Valid {
		return User{}, ErrInvalidCode
	}
	if err != nil {
		return User{}, err
	}
	if time.Now().After(expiresAt) {
		return User{}, ErrChallengeExpired
	}
	if attempts >= maxAttempts {
		return User{}, ErrTooManyAttempts
	}

	if subtle.ConstantTimeCompare([]byte(storedHash), []byte(candidateHash)) != 1 {
		attempts++
		if _, err := tx.ExecContext(ctx, `UPDATE email_challenges SET attempts = $2 WHERE id = $1`, challengeID, attempts); err != nil {
			return User{}, err
		}
		if err := tx.Commit(); err != nil {
			return User{}, err
		}
		if attempts >= maxAttempts {
			return User{}, ErrTooManyAttempts
		}
		return User{}, ErrInvalidCode
	}

	var user User
	if purpose == "register" {
		if !pendingName.Valid {
			return User{}, ErrInvalidCode
		}
		userID, err := NewID()
		if err != nil {
			return User{}, err
		}
		err = tx.QueryRowContext(ctx, `
			INSERT INTO users (id, email, name, role, verified_at)
			VALUES ($1, $2, $3, $4, now())
			ON CONFLICT (email) DO NOTHING
			RETURNING id, name, email, city, role, verified_at, created_at`,
			userID, email, pendingName.String, registrationRole).
			Scan(&user.ID, &user.Name, &user.Email, &user.City, &user.Role, &user.VerifiedAt, &user.CreatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrEmailExists
		}
	} else {
		err = tx.QueryRowContext(ctx, `
			SELECT id, name, email, city, role, verified_at, created_at
			FROM users WHERE email = $1`, email).
			Scan(&user.ID, &user.Name, &user.Email, &user.City, &user.Role, &user.VerifiedAt, &user.CreatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrInvalidCode
		}
	}
	if err != nil {
		return User{}, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE email_challenges SET consumed_at = now() WHERE id = $1`, challengeID); err != nil {
		return User{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_sessions (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)`, sessionID, user.ID, sessionHash, sessionExpiresAt); err != nil {
		return User{}, err
	}
	if err := tx.Commit(); err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Store) UserBySession(ctx context.Context, tokenHash string) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.city, u.role, u.verified_at, u.created_at
		FROM user_sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1
		  AND s.revoked_at IS NULL
		  AND s.expires_at > now()`, tokenHash).
		Scan(&user.ID, &user.Name, &user.Email, &user.City, &user.Role, &user.VerifiedAt, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (s *Store) RevokeSession(ctx context.Context, tokenHash string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE user_sessions SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL`, tokenHash)
	return err
}

func (s *Store) UpdateProfile(ctx context.Context, userID, name, city string) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, `
		UPDATE users
		SET name = $2, city = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, email, city, role, verified_at, created_at`,
		userID, name, city).
		Scan(&user.ID, &user.Name, &user.Email, &user.City, &user.Role, &user.VerifiedAt, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func normalizeSearch(value string) string {
	return strings.TrimSpace(value)
}
