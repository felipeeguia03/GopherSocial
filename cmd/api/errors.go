package main

import (
	"net/http"
)

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusInternalServerError, err.Error())

}

func (app *application) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusNotFound, err.Error())
}

func (app *application) ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusConflict, err.Error())
}

func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	WriteJsonError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJsonError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

	WriteJsonError(w, http.StatusForbidden, "forbidden")
}
