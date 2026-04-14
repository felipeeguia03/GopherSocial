package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/felipeeguia03/vol7/internal/store"
	"github.com/go-chi/chi/v5"
)

type contextKey string

const (
	UserKey     contextKey = "user"
	AuthUserKey contextKey = "authUser"
)

// GetUserHandler godoc
//
//	@Summary		Gets an user
//	@Description	gets an user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int			true	"userID"
//	@Success		204		{object}	store.User	"user data"
//	@Failure		400		{object}	error		"invalid payload"
//	@Failure		404		{object}	error		"user not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromContext(r)
	me := app.getAuthUserFromContext(r)

	if me != nil && me.ID != user.ID {
		isFollowing, err := app.store.Followers.IsFollowing(r.Context(), me.ID, user.ID)
		if err == nil {
			user.IsFollowing = isFollowing
		}
	}

	if err := JsonResponse(w, http.StatusOK, user); err != nil {
		app.InternalServerError(w, r, err)
	}
}

type FollowUserPayload struct {
	FollowerID int64 `json:"follower_id"`
}

// unfollowUser godoc
//
//	@Summary		Follows an user
//	@Description	Follows an user by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"userID"
//	@Success		204		{string}	string	"user followed"
//	@Failure		400		{object}	error	"invalid payload"
//	@Failure		404		{object}	error	"user not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followHandler(w http.ResponseWriter, r *http.Request) {
	// target = el usuario a quien se quiere seguir (viene del {userID} en la URL)
	target := app.getUserFromContext(r)
	// me = el usuario autenticado que hace la acción
	me := app.getAuthUserFromContext(r)

	if me.ID == target.ID {
		app.BadRequestError(w, r, errors.New("cannot follow yourself"))
		return
	}

	// followers(user_id=quien_sigue, follower_id=seguido) — según el schema existente
	err := app.store.Followers.Follow(r.Context(), me.ID, target.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.ConflictError(w, r, errors.New("already following"))
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// unfollowUser godoc
//
//	@Summary		Unfollows an user
//	@Description	Unfollows an user by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"userID"
//	@Success		204		{string}	string	"user unfollowed"
//	@Failure		400		{object}	error	"invalid payload"
//	@Failure		404		{object}	error	"user not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	target := app.getUserFromContext(r)
	me := app.getAuthUserFromContext(r)

	err := app.store.Followers.Unfollow(r.Context(), me.ID, target.ID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// activateUserHandler godoc
//
//	@Summary		Activates an user account
//	@Description	Activates an user using the invitation token
//	@Tags			user
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"user activated"
//	@Failure		404		{object}	error	"token not found or expired"
//	@Failure		500		{object}	error
//	@Router			/users/activate/{token} [put]
func (app *application) getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromContext(r)
	posts, err := app.store.Posts.GetPostsByUserID(r.Context(), user.ID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	if err := JsonResponse(w, http.StatusOK, posts); err != nil {
		app.InternalServerError(w, r, err)
	}
}

// searchUsersHandler godoc
//
//	@Summary		Search users by username
//	@Description	Returns users matching the query, including whether the authenticated user follows each one
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			q	query		string		true	"Username search query"
//	@Success		200	{array}		store.User	"list of matching users with is_following field"
//	@Failure		400	{object}	error		"q is required"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/search [get]
func (app *application) searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		app.BadRequestError(w, r, errors.New("q is required"))
		return
	}

	me := app.getAuthUserFromContext(r)

	users, err := app.store.Users.SearchByUsername(r.Context(), me.ID, q)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, users); err != nil {
		app.InternalServerError(w, r, err)
	}
}

// getSuggestedUsersHandler godoc
//
//	@Summary		Get suggested users to follow
//	@Description	Returns users not yet followed by the authenticated user, with is_following field
//	@Tags			user
//	@Produce		json
//	@Success		200	{array}		store.User	"list of suggested users with is_following field"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/suggested [get]
func (app *application) getSuggestedUsersHandler(w http.ResponseWriter, r *http.Request) {
	me := app.getAuthUserFromContext(r)

	users, err := app.store.Users.GetSuggestedUsers(r.Context(), me.ID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, users); err != nil {
		app.InternalServerError(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.NotFoundError(w, r, err)
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(id)
		if err != nil {
			app.InternalServerError(w, r, err)
			return
		}

		user, err := app.store.Users.GetUserByID(ctx, int64(userID))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.NotFoundError(w, r, err)
				return
			}
			app.InternalServerError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) getUserFromContext(r *http.Request) *store.User {
	ctx := r.Context()
	user, _ := ctx.Value(UserKey).(*store.User)
	return user
}

func (app *application) getAuthUserFromContext(r *http.Request) *store.User {
	ctx := r.Context()
	user, _ := ctx.Value(AuthUserKey).(*store.User)
	return user
}
