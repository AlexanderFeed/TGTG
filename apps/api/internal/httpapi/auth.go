package httpapi

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
	codePattern = regexp.MustCompile(`^[0-9]{6}$`)
	idPattern   = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

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

func normalizeEmail(raw string) (string, bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	parsed, err := mail.ParseAddress(email)
	return email, err == nil && parsed.Address == email && len(email) <= 320
}

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
	if s.cfg.EmailDelivery == "development" && s.cfg.ExposeDevCodes {
		response.DevCode = code
	}
	writeJSON(w, http.StatusAccepted, response)
}

func (s *Server) verifyRegistration(w http.ResponseWriter, r *http.Request) {
	s.verifyCode(w, r, "register")
}

func (s *Server) verifyLogin(w http.ResponseWriter, r *http.Request) {
	s.verifyCode(w, r, "login")
}

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
	sessionToken, err := newToken()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	role := "customer"
	if purpose == "register" && s.cfg.IsAdmin(email) {
		role = "admin"
	}
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
	s.setSessionCookie(w, sessionToken)
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (s *Server) hashCode(challengeID, code string) string {
	digest := sha256.Sum256([]byte(s.cfg.OTPPepper + ":" + challengeID + ":" + code))
	return hex.EncodeToString(digest[:])
}

func newCode() (string, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", value.Int64()+100000), nil
}

func newToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"user": currentUser(r)})
}

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
