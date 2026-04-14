package main

import (
	"log"
	"time"

	"github.com/felipeeguia03/vol7/internal/auth"
	"github.com/felipeeguia03/vol7/internal/db"
	"github.com/felipeeguia03/vol7/internal/env"
	"github.com/felipeeguia03/vol7/internal/mailer"
	"github.com/felipeeguia03/vol7/internal/store"
	"github.com/felipeeguia03/vol7/internal/store/cache"
	"go.uber.org/zap"
)

const version = ""

//	@title	Swagger V7 API curso Backend

//	@description	V7 de la API para gophersocial
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  "development",
		dbConfig: dbConfig{
			dsn:          env.GetString("DB_ADDR", "postgres://root:root@localhost/vol7?sslmode=disable"),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 20),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "3m"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 20),
		},
		mail: mailConfig{
			fromEmail: env.GetString("FROM_EMAIL", "onboarding@resend.dev"),
			exp:       time.Hour * 24 * 3, //3 days
			mailtrap: mailtrapConfig{
				APIKey: env.GetString("MAIL_TRAP_API", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("USERNAME_BASIC", "admin"),
				password: env.GetString("PASSWORD_BASIC", "password"),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    env.GetString("TOKEN_ISS", "gophersocial"),
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		apiURL:      env.GetString("API_URL", "localhost:8080"),
		redisConfig: redisConfig{
			addr:    env.GetString("REDIS_ADDR", ""),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetString("REDIS_ENABLED", "false") == "true",
		},
	}

	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.dbConfig.dsn,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infow("Base de datos conectada correctamente")
	logger.Infow("config", "fromEmail", cfg.mail.fromEmail, "apiKeyLen", len(cfg.mail.mailtrap.APIKey), "frontendURL", cfg.frontendURL)

	store := store.NewStorage(db)

	var cacheStorage cache.Storage
	if cfg.redisConfig.enabled {
		rdb := cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.pw, cfg.redisConfig.db)
		cacheStorage = cache.NewRedisStorage(rdb)
		logger.Infow("Redis conectado correctamente")
	}

	mailerClient, err := mailer.NewMailTrapClient(cfg.mail.mailtrap.APIKey, cfg.mail.fromEmail)
	if err != nil {
		log.Fatal(err)
	}

	jwtAuth := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	api := application{
		config:       cfg,
		store:        store,
		cacheStorage: cacheStorage,
		logger:       logger,
		mailer:       mailerClient,
		auth:         jwtAuth,
	}

	if err := api.run(api.mount()); err != nil {
		log.Fatal(err)
	}

}
