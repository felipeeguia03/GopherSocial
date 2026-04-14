package store

import (
	"context"
	"database/sql"
)

type CommentsStore struct {
	db *sql.DB
}

type Comment struct {
	ID         int64  `json:"id"`
	PostID     int64  `json:"post_id"`
	UserID     int64  `json:"user_id"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
	User       User   `json:"user"`
}

func (s *CommentsStore) GetCommentsByPostID(ctx context.Context, postID int64) ([]*Comment, error) {

	comments := make([]*Comment, 0)
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.id, users.username
	FROM comments as c inner join users on users.id = c.user_id
	WHERE post_id = $1 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cmt, err := scanRowsIntoComments(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, cmt)
	}
	return comments, nil
}

func scanRowsIntoComments(rows *sql.Rows) (*Comment, error) {
	cmt := new(Comment)
	err := rows.Scan(
		&cmt.ID,
		&cmt.PostID,
		&cmt.UserID,
		&cmt.Content,
		&cmt.Created_at,
		&cmt.User.ID,
		&cmt.User.Username,
	)
	if err != nil {
		return nil, err
	}
	return cmt, nil
}

func (s *CommentsStore) Create(ctx context.Context, cmt *Comment) error {
	query := `INSERT into comments (user_id, post_id, content) VALUES ($1,$2,$3) returning id, created_at `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, cmt.UserID, cmt.PostID, cmt.Content).Scan(
		&cmt.ID,
		&cmt.Created_at,
	)

	if err != nil {
		return err
	}
	return nil
}
