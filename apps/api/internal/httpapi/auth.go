package httpapi

// This file contains authentication/profile HTTP handlers. The frontend calls
// these handlers indirectly through Nuxt's /api proxy; route registration is
// in server.go and persistent database work is delegated to package store.

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/eshche-est/eshche-est/apps/api/internal/store"
)

var (
	// Compile regular expressions once at package startup, not for each request.
	codePattern = regexp.MustCompile(`^[0-9]{6}$`)
	idPattern   = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// Input structs describe JSON received from the frontend. JSON tags map the
// frontend's camelCase field names to Go's exported struct fields.
type requestCodeInput struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type verifyCodeInput struct {
	ChallengeID string `json:"challengeId"`
	Email       string `json:"email"`
	Code        string `json:"code"`
}

type challengeResponse struct {
	ChallengeID      string `json:"challengeId"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
	Delivery         string `json:"delivery"`
	DevCode          string `json:"devCode,omitempty"`
}

// normalizeEmail makes email lookup deterministic: PostgreSQL always receives
// a lowercase, whitespace-free form. mail.ParseAddress checks basic syntax.
func normalizeEmail(raw string) (string, bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	parsed, err := mail.ParseAddress(email)
	return email, err == nil && parsed.Address == email && len(email) <= 320
}

// requestRegistration handles POST /v1/auth/register/request. It validates the
// first registration form and creates an email challenge; it does not create a
// user yet. A user is created only after the code is verified.
func (s *Server) requestRegistration(w http.ResponseWriter, r *http.Request) {
	var input requestCodeInput
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	email, valid := normalizeEmail(input.Email)
	if !valid || len([]rune(input.Name)) < 2 || len([]rune(input.Name)) > 80 {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "Укажите имя и корректный email.")
		return
	}
	exists, err := s.store.UserExists(r.Context(), email)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "email_already_registered", "Аккаунт с таким email уже существует.")
		return
	}
	s.requestCode(w, r, email, input.Name, "register", true)
}

// requestLogin handles POST /v1/auth/login/request for existing users.
func (s *Server) requestLogin(w http.ResponseWriter, r *http.Request) {
	var input requestCodeInput
	if !decodeJSON(w, r, &input) {
		return
	}
	email, valid := normalizeEmail(input.Email)
	if !valid {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "Укажите корректный email.")
		return
	}
	exists, err := s.store.UserExists(r.Context(), email)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	// Store the same challenge shape and apply the same cooldown for unknown
	// addresses, but do not send mail to arbitrary inboxes. Verification can only
	// succeed when a matching user exists.
	s.requestCode(w, r, email, "", "login", exists)
}

// requestCode is shared by register and login. It generates a random code,
// stores only its hash, optionally sends the plaintext code, and returns a
// challenge ID the frontend must present during verification.
func (s *Server) requestCode(w http.ResponseWriter, r *http.Request, email, name, purpose string, send bool) {
	challengeID, err := store.NewID()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	code, err := newCode()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	// The challenge ID is part of the hash. Thus identical six-digit codes for
	// two different challenges do not create the same stored hash.
	challenge := store.Challenge{
		ID: challengeID, Email: email, Purpose: purpose, PendingName: name,
		CodeHash: s.hashCode(challengeID, code), MaxAttempts: s.cfg.OTPMaxAttempts,
		ExpiresAt: time.Now().Add(s.cfg.OTPExpiry),
	}
	if err := s.store.CreateChallenge(r.Context(), challenge, s.cfg.OTPCooldown); err != nil {
		if errors.Is(err, store.ErrCooldown) {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(s.cfg.OTPCooldown.Seconds())))
			writeError(w, http.StatusTooManyRequests, "code_cooldown", "Новый код можно запросить через минуту.")
			return
		}
		s.internalError(w, r, err)
		return
	}
	if send {
		// DevelopmentSender does nothing; SMTPSender sends real mail. The handler
		// calls the same interface in either environment.
		if err := s.mailer.SendVerificationCode(r.Context(), email, code, purpose); err != nil {
			s.store.InvalidateChallenge(r.Context(), challengeID)
			s.internalError(w, r, err)
			return
		}
	}
	response := challengeResponse{
		ChallengeID:      challengeID,
		ExpiresInSeconds: int64(s.cfg.OTPExpiry.Seconds()),
		Delivery:         "email",
	}
	// devCode exists only for local learning/testing. SMTP configuration rejects
	// EXPOSE_DEV_CODES, so a real deployment never returns the secret code.
	if s.cfg.EmailDelivery == "development" && s.cfg.ExposeDevCodes {
		response.DevCode = code
	}
	writeJSON(w, http.StatusAccepted, response)
}

// verifyRegistration and verifyLogin are thin route-specific adapters around
// the same verification workflow.
func (s *Server) verifyRegistration(w http.ResponseWriter, r *http.Request) {
	s.verifyCode(w, r, "register")
}

func (s *Server) verifyLogin(w http.ResponseWriter, r *http.Request) {
	s.verifyCode(w, r, "login")
}

// verifyCode validates the submitted code and asks Store.CompleteChallenge to
// consume it, create/find the user, and create a session in one transaction.
func (s *Server) verifyCode(w http.ResponseWriter, r *http.Request, purpose string) {
	var input verifyCodeInput
	if !decodeJSON(w, r, &input) {
		return
	}
	email, validEmail := normalizeEmail(input.Email)
	if !validEmail || !codePattern.MatchString(input.Code) || !idPattern.MatchString(input.ChallengeID) {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "Проверьте email и шестизначный код.")
		return
	}
	sessionID, err := store.NewID()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	// The random session token goes to the browser. Only hashToken(sessionToken)
	// crosses into PostgreSQL, which is the same pattern used for password hashes.
	sessionToken, err := newToken()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	role := "customer"
	if purpose == "register" && s.cfg.IsAdmin(email) {
		role = "admin"
	}
	// This long call passes prepared values across the HTTP -> store boundary.
	// The store owns the transaction because it owns the related SQL operations.
	user, err := s.store.CompleteChallenge(
		r.Context(), input.ChallengeID, email, purpose,
		s.hashCode(input.ChallengeID, input.Code), sessionID, hashToken(sessionToken),
		time.Now().Add(s.cfg.SessionTTL), role,
	)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrChallengeExpired):
			writeError(w, http.StatusGone, "code_expired", "Код истёк. Запросите новый.")
		case errors.Is(err, store.ErrTooManyAttempts):
			writeError(w, http.StatusTooManyRequests, "too_many_attempts", "Слишком много попыток. Запросите новый код.")
		case errors.Is(err, store.ErrInvalidCode):
			writeError(w, http.StatusUnauthorized, "invalid_code", "Неверный код подтверждения.")
		case errors.Is(err, store.ErrEmailExists):
			writeError(w, http.StatusConflict, "email_already_registered", "Аккаунт с таким email уже существует.")
		default:
			s.internalError(w, r, err)
		}
		return
	}
	// Set-Cookie travels back through the Nuxt proxy and is stored by the browser.
	s.setSessionCookie(w, sessionToken)
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

// hashCode mixes a server-only pepper, public challenge ID, and short OTP. A
// pepper is important because a six-digit value is cheap to brute-force if an
// attacker obtains only the database.
func (s *Server) hashCode(challengeID, code string) string {
	digest := sha256.Sum256([]byte(s.cfg.OTPPepper + ":" + challengeID + ":" + code))
	return hex.EncodeToString(digest[:])
}

// newCode uses crypto/rand rather than math/rand because OTPs are secrets.
func newCode() (string, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", value.Int64()+100000), nil
}

// newToken returns 256 random bits encoded in a cookie-safe URL format.
func newToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// me returns the user already loaded by requireAuth middleware.
func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"user": currentUser(r)})
}

// logout revokes the server-side session and also removes the browser cookie.
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(s.cfg.CookieName); err == nil && cookie.Value != "" {
		if err := s.store.RevokeSession(r.Context(), hashToken(cookie.Value)); err != nil {
			s.internalError(w, r, err)
			return
		}
	}
	s.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

type profileInput struct {
	Name string `json:"name"`
	City string `json:"city"`
}

// updateProfile handles PATCH /v1/users/me. The authenticated user ID comes
// from middleware, never from editable frontend input.
func (s *Server) updateProfile(w http.ResponseWriter, r *http.Request) {
	var input profileInput
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.City = strings.TrimSpace(input.City)
	if len([]rune(input.Name)) < 2 || len([]rune(input.Name)) > 80 || len([]rune(input.City)) < 2 || len([]rune(input.City)) > 80 {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "Имя и город должны содержать от 2 до 80 символов.")
		return
	}
	user, err := s.store.UpdateProfile(r.Context(), currentUser(r).ID, input.Name, input.City)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}
