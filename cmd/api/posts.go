package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/felipeeguia03/vol7/internal/store"
	"github.com/go-chi/chi/v5"
)

type postContext string

var postKey postContext = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"omitempty,max=100"`
	Content string   `json:"content" validate:"omitempty,max=100"`
	UserId  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"post data"
//	@Success		201		{object}	store.Post			"created post"
//	@Failure		400		{object}	error				"invalid payload"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	ctx := r.Context()
	if err := ReadJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}
	post := new(store.Post)
	post.Title = payload.Title
	post.Content = payload.Content
	post.UserID = payload.UserId
	post.Tags = payload.Tags

	err := app.store.Posts.Create(ctx, post)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	JsonResponse(w, http.StatusCreated, post)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"postID"
//	@Param			payload	body		UpdatePostPayload	true	"post data"
//	@Success		201		{object}	store.Post			"updated post"
//	@Failure		400		{object}	error				"invalid payload"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdatePostPayload

	post := app.getPostFromContext(r)

	ctx := r.Context()
	if err := ReadJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err := app.store.Posts.Update(ctx, post)
	if err != nil {
		switch {
		case errors.Is(err, store.NotFoundError):
			app.NotFoundError(w, r, err)
			return
		default:
			app.InternalServerError(w, r, err)
			return
		}
	}

	if err := JsonResponse(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
	}
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int			true	"postID"
//	@Success		201		{object}	store.Post	" post data"
//	@Failure		400		{object}	error		"invalid payload"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := app.getPostFromContext(r)

	ctx := r.Context()
	comments, err := app.store.Comments.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := JsonResponse(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Deletes a post
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"postID"
//	@Success		201		{string}	string	"deleted"
//	@Failure		400		{object}	error	"invalid payload"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := app.getPostFromContext(r)
	ctx := r.Context()

	err := app.store.Posts.Delete(ctx, post.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.NotFoundError):
			app.NotFoundError(w, r, err)
			return
		default:
			app.InternalServerError(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := app.getPostFromContext(r)
	me := app.getAuthUserFromContext(r)

	var payload CreateCommentPayload
	if err := ReadJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	cmt := &store.Comment{
		PostID:  post.ID,
		UserID:  me.ID,
		Content: payload.Content,
	}
	if err := app.store.Comments.Create(r.Context(), cmt); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusCreated, cmt); err != nil {
		app.InternalServerError(w, r, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id := chi.URLParam(r, "postID")
		postID, err := strconv.Atoi(id)
		if err != nil {
			app.InternalServerError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetPostByID(ctx, int64(postID))
		if err != nil {
			switch {
			case errors.Is(err, store.NotFoundError):
				app.NotFoundError(w, r, err)
				return
			}
			app.InternalServerError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) getPostFromContext(r *http.Request) *store.Post {
	ctx := r.Context()
	post, _ := ctx.Value(postKey).(*store.Post)
	return post
}
