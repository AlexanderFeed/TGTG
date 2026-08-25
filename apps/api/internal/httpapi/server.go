package httpapi

// Package httpapi is the transport layer: it translates HTTP requests into
// calls to the store/mailer layers, then translates results back into JSON.
// It should validate input and enforce access rules, but it should not contain
// raw SQL or know how PostgreSQL connections are created.

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

// Server is a dependency container for HTTP handlers. Handler methods have a
// *Server receiver so they can share configuration, database access, email,
// and logging without global variables.
type Server struct {
	cfg    config.Config
	store  *store.Store
	mailer mailer.Sender
	log    *slog.Logger
}

// New builds the complete backend router and returns the standard-library
// http.Handler interface expected by http.Server in cmd/api/main.go.
//
// Every incoming Go API request first passes through the middleware registered
// with Use, and is then dispatched to one handler registered below.
func New(cfg config.Config, dataStore *store.Store, sender mailer.Sender, log *slog.Logger) http.Handler {
	server := &Server{cfg: cfg, store: dataStore, mailer: sender, log: log}
	router := chi.NewRouter()

	// Middleware wraps handlers like nested functions. Registration order is
	// request order; the response travels back through them in reverse order.
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(20 * time.Second))
	router.Use(server.securityHeaders)
	router.Use(server.limitBody)
	router.Use(server.logRequest)
	router.Use(server.requireMutationHeader)

	// Public operational endpoint used by Docker/monitoring.
	router.Get("/healthz", server.health)

	// These are the real Go API entry points. The browser calls /api/v1/..., and
	// Nuxt's server/api/[...path].ts proxy removes only the /api prefix before
	// forwarding here, so /api/v1/auth/me becomes /v1/auth/me.
	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register/request", server.requestRegistration)
			r.Post("/register/verify", server.verifyRegistration)
			r.Post("/login/request", server.requestLogin)
			r.Post("/login/verify", server.verifyLogin)
		})
		r.Get("/offers", server.listOffers)
		r.Get("/offers/{offerID}", server.getOffer)

		// Everything in this group runs requireAuth before its final handler.
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

// health confirms both the HTTP process and its required database are alive.
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.Ping(ctx); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database_unavailable", "База данных недоступна.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// securityHeaders adds browser hardening and prevents private API responses
// from being cached by an intermediary.
func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// limitBody caps all request bodies at 1 MiB before a handler decodes JSON.
func (s *Server) limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		next.ServeHTTP(w, r)
	})
}

// requireMutationHeader rejects state-changing browser requests that did not
// come through our frontend client. Together with SameSite cookies and the
// same-origin Nuxt proxy, this is a simple CSRF defense for the pet project.
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

// logRequest wraps ResponseWriter so it can record the final status code and
// duration after the downstream handler finishes.
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

// An unexported, zero-size key type avoids collisions with context values from
// other packages. Context is for request-scoped values, not general app state.
type userContextKey struct{}

// requireAuth is authentication middleware. It reads the opaque session token
// from the HttpOnly cookie, hashes it, and asks PostgreSQL for a live session.
// The raw token never needs to be stored in the database.
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
		// Attach the authenticated user to this request only. A later handler can
		// read it with currentUser without querying the database a second time.
		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// currentUser reads the value placed in context by requireAuth. It must only be
// called inside the protected router group.
func currentUser(r *http.Request) store.User {
	user, _ := r.Context().Value(userContextKey{}).(store.User)
	return user
}

// hashToken creates the value stored and queried in user_sessions. If the
// database leaks, its hashes cannot directly be used as browser cookies.
func hashToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}

// setSessionCookie sends the raw opaque token only to the browser. HttpOnly
// prevents frontend JavaScript from reading it; the browser attaches it to
// later same-origin requests automatically.
func (s *Server) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cfg.CookieName, Value: token, Path: "/",
		MaxAge: int(s.cfg.SessionTTL.Seconds()), Expires: time.Now().Add(s.cfg.SessionTTL),
		HttpOnly: true, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
	})
}

// clearSessionCookie expires the browser's copy during logout or after an
// invalid/expired database session is encountered.
func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cfg.CookieName, Value: "", Path: "/", MaxAge: -1,
		Expires: time.Unix(1, 0), HttpOnly: true, Secure: s.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

// internalError logs technical details on the server while returning a generic
// message to the client, so SQL/SMTP internals do not leak through the API.
func (s *Server) internalError(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Error("request failed", "request_id", middleware.GetReqID(r.Context()), "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось выполнить запрос. Попробуйте ещё раз.")
}
