package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/abarragan89/go-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	UserID  string   `json:"user_id"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// write the response
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type SinglePost struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	UserID  string   `json:"user_id"`
	Tags    []string `json:"tags"`
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	postID := r.PathValue("postID")
	// chi package way
	// postID := chi.URLParam(r, "postID")

	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postIDInt)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}
