package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/felipeeguia03/vol7/docs"
	"github.com/felipeeguia03/vol7/internal/auth"
	"github.com/felipeeguia03/vol7/internal/env"
	"github.com/felipeeguia03/vol7/internal/mailer"
	"github.com/felipeeguia03/vol7/internal/ratelimiter"
	"github.com/felipeeguia03/vol7/internal/store"
	"github.com/felipeeguia03/vol7/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config       config
	store        store.Storage
	mailer       mailer.Client
	auth         auth.Authenticator
	logger       *zap.SugaredLogger
	cacheStorage cache.Storage
	rateLimiter  ratelimiter.Limiter
}

type config struct {
	addr        string
	env         string
	dbConfig    dbConfig
	mail        mailConfig
	auth        authConfig
	apiURL      string
	frontendURL string
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	enabled bool
}

type dbConfig struct {
	dsn          string
	maxIdleConns int
	maxOpenConns int
	maxIdleTime  string
}

type mailConfig struct {
	fromEmail string
	exp       time.Duration
	mailtrap  mailtrapConfig
}
type mailtrapConfig struct {
	APIKey string
}

type authConfig struct {
	token tokenConfig
	basic basicConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	username string
	password string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.config.rateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	r.Route("/v1", func(r chi.Router) {

		//health
		r.Get("/health", app.handleHealth)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		//swaggo
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		//auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/token", app.createUserTokenHandler)
		})

		//users
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Group(func(r chi.Router) {
				r.Use(app.authTokenMiddleware)
				r.Get("/search", app.searchUsersHandler)
				r.Route("/{userID}", func(r chi.Router) {
					r.Use(app.userContextMiddleware)
					r.Get("/", app.getUserHandler)
					r.Get("/posts", app.getUserPostsHandler)
					r.Put("/follow", app.followHandler)
					r.Put("/unfollow", app.unfollowHandler)
				})
			})
		})

		// rutas protegidas (requieren JWT)
		r.Group(func(r chi.Router) {
			r.Use(app.authTokenMiddleware)

			//posts
			r.Route("/posts", func(r chi.Router) {
				r.Route("/{postID}", func(r chi.Router) {
					r.Use(app.postContextMiddleware)
					r.Get("/", app.getPostHandler)
					r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
					r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
					r.Post("/comments", app.createCommentHandler)
				})
				r.Post("/", app.createPostHandler)
			})

			//feed
			r.Route("/feed", func(r chi.Router) {
				r.Get("/", app.getUserFeedHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 20,
		IdleTimeout:  time.Second * 10,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
