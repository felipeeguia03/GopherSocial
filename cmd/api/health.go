package main

import "net/http"

// GetHealth godoc
//
//	@Summary		Fetches Health
//	@Description	Fetches Health
//	@Tags			DevOps
//	@Accept			json
//	@Produce		json
//	@Success		200	{obejct}	map[string]string
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app *application) handleHealth(w http.ResponseWriter, r *http.Request) {
	envelope := map[string]string{
		"status":      "ok",
		"version":     version,
		"environment": app.config.env,
	}

	err := JsonResponse(w, http.StatusOK, envelope)
	if err != nil {
		app.InternalServerError(w, r, err)
	}
}
