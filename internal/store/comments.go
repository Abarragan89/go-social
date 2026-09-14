package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CreatedAt string `json:"created_at"`
	User
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (content, user_id, post_id)
		VALUES ( $1, $2, $3)
		RETURNING id, content, created_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.UserID,
		comment.PostID,
	).Scan(
		&comment.ID,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postId int64) (*[]Comment, error) {
	query := `
		SELECT c.id, c.content, c.created_at, u.username, c.user_id FROM comments c
		JOIN users u On u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC; 
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		postId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	// returns false when no more rows
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}
