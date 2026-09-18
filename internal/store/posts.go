package store

import (
	"context"
	"database/sql"
	"errors"

	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts(content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, tags, version
		FROM posts	
		WHERE id = $1
	`

	var post Post
	err := s.db.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostsStore) DeleteByID(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(
		ctx,
		query,
		postID,
	)
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) PatchPostByID(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW(), version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING updated_at, version
	`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).
		Scan(&post.UpdatedAt, &post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
