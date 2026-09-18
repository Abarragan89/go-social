package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/abarragan89/go-social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
	UserID  int64  `json:"user_id"`
	PostID  int64  `json:"post_id"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		UserID:  payload.UserID,
		PostID:  payload.PostID,
	}

	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getCommentsByPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("postID")

	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	ctx := r.Context()

	post, err := app.store.Comments.GetByPostID(ctx, postIDInt)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}

}
