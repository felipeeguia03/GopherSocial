package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felipeeguia03/vol7/internal/mailer"
	"github.com/felipeeguia03/vol7/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100" `
	Email    string `json:"email" validate:"required,max=255" `
	Password string `json:"password" validate:"required,min=3" `
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUser godoc
//
//	@Summary		Registers an User
//	@Description	Register an user
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User Credentials"
//	@Success		200		{obejct}	UserWithToken		"user with token"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var payload RegisterUserPayload
	if err := ReadJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	user := new(store.User)
	user.Username = payload.Username
	user.Email = payload.Email

	if err := user.Password.Set(payload.Password); err != nil {
		log.Print("error puto")
		app.InternalServerError(w, r, err)
		return
	}

	//create token

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	//store the user

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.ConflictError(w, r, err)
			return
		case errors.Is(err, store.ErrDuplicatedEmail):
			app.BadRequestError(w, r, err)
			return
		case errors.Is(err, store.ErrDuplicatedUsername):
			app.BadRequestError(w, r, err)
			return
		default:
			log.Println("hola")
			app.InternalServerError(w, r, err)
			return
		}

	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	//send invitation

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA PATTERN)

		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.InternalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := JsonResponse(w, http.StatusOK, userWithToken); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3"`
}

// createUserToken godoc
//
//	@Summary		Creates an user token
//	@Description	Creates an user token with email and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User Credentials"
//	@Success		200		{obejct}	map[string]string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/auth/token [get]
func (app *application) createUserTokenHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var payload CreateUserTokenPayload

	if err := ReadJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	// check if user exists

	user, err := app.store.Users.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.UnauthorizedErrorResponse(w, r, err)
		case errors.Is(err, store.ErrNotActivated):
			app.UnauthorizedErrorResponse(w, r, errors.New("account is not activated, check your email"))
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}

	//check if password matches

	err = user.Password.Compare(payload.Password)
	if err != nil {
		app.UnauthorizedErrorResponse(w, r, err)
		return
	}

	// generate token --> add claims

	//3 generate token -> add claims

	if user == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.auth.GenerateToken(claims)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	//4 send token to user

	if err := JsonResponse(w, http.StatusCreated, token); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

}

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.UnauthorizedErrorResponse(w, r, errors.New("authorization header missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.UnauthorizedErrorResponse(w, r, errors.New("invalid authorization header format"))
			return
		}

		token, err := app.auth.ValidateToken(parts[1])
		if err != nil {
			app.UnauthorizedErrorResponse(w, r, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.UnauthorizedErrorResponse(w, r, errors.New("invalid token claims"))
			return
		}

		userIDStr := fmt.Sprintf("%.f", claims["sub"])
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			app.UnauthorizedErrorResponse(w, r, errors.New("invalid user id in token"))
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetUserByID(ctx, userID)
		if err != nil {
			app.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, UserKey, user)
		ctx = context.WithValue(ctx, AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
