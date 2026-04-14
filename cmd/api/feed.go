package main

import (
	"net/http"

	"github.com/felipeeguia03/vol7/internal/store"
)

// GetFeed godoc
//
//	@Summary		Fetches Feed
//	@Description	Fetches Feed to a user by ID
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		string	true	"Until"
//	@Param			offset	query		string	true	"Until"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Param			sort	query		string	true	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//get the user

	//get feed

	fq := new(store.PaginatedFeedQuery)
	feedQuery, err := fq.Parse(r)
	if err != nil {
		app.InternalServerError(w, r, err)
	}

	user := app.getUserFromContext(r)
	feed, err := app.store.Posts.GetUserFeed(ctx, user.ID, feedQuery)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, feed)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}
