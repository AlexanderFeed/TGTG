package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/eshche-est/eshche-est/apps/api/internal/config"
	"github.com/eshche-est/eshche-est/apps/api/internal/mailer"
	"github.com/eshche-est/eshche-est/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg    config.Config
	store  *store.Store
	mailer mailer.Sender
	log    *slog.Logger
}

func New(cfg config.Config, dataStore *store.Store, sender mailer.Sender, log *slog.Logger) http.Handler {
	server := &Server{cfg: cfg, store: dataStore, mailer: sender, log: log}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(20 * time.Second))
	router.Use(server.securityHeaders)
	router.Use(server.limitBody)
	router.Use(server.logRequest)
	router.Use(server.requireMutationHeader)

	router.Get("/healthz", server.health)
	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register/request", server.requestRegistration)
			r.Post("/register/verify", server.verifyRegistration)
			r.Post("/login/request", server.requestLogin)
			r.Post("/login/verify", server.verifyLogin)
		})
		r.Get("/offers", server.listOffers)
		r.Get("/offers/{offerID}", server.getOffer)

		r.Group(func(r chi.Router) {
			r.Use(server.requireAuth)
			r.Get("/auth/me", server.me)
			r.Post("/auth/logout", server.logout)
			r.Patch("/users/me", server.updateProfile)
			r.Post("/offers", server.createOffer)
			r.Put("/offers/{offerID}", server.updateOffer)
			r.Delete("/offers/{offerID}", server.deleteOffer)
		})
	})

	return router
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.Ping(ctx); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database_unavailable", "База данных недоступна.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireMutationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if r.Header.Get("X-Requested-With") != "eshche-est-web" {
				writeError(w, http.StatusForbidden, "request_header_required", "Защитный заголовок запроса отсутствует.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(wrapped, r)
		s.log.Info("http request",
			"request_id", middleware.GetReqID(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.Status(),
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}

type userContextKey struct{}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.cfg.CookieName)
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			writeError(w, http.StatusUnauthorized, "authentication_required", "Войдите в аккаунт.")
			return
		}
		user, err := s.store.UserBySession(r.Context(), hashToken(cookie.Value))
		if errors.Is(err, store.ErrNotFound) {
			s.clearSessionCookie(w)
			writeError(w, http.StatusUnauthorized, "authentication_required", "Сессия недействительна или истекла.")
			return
		}
		if err != nil {
			s.internalError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentUser(r *http.Request) store.User {
	user, _ := r.Context().Value(userContextKey{}).(store.User)
	return user
}

func hashToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}

func (s *Server) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cfg.CookieName, Value: token, Path: "/",
		MaxAge: int(s.cfg.SessionTTL.Seconds()), Expires: time.Now().Add(s.cfg.SessionTTL),
		HttpOnly: true, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cfg.CookieName, Value: "", Path: "/", MaxAge: -1,
		Expires: time.Unix(1, 0), HttpOnly: true, Secure: s.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) internalError(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Error("request failed", "request_id", middleware.GetReqID(r.Context()), "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось выполнить запрос. Попробуйте ещё раз.")
}
