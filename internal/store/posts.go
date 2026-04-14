package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

type Post struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	UserID     int64      `json:"user_id"`
	Tags       []string   `json:"tags"`
	Comments   []*Comment `json:"comments"`
	Version    int        `json:"version"`
	Created_at string     `json:"created_at"`
	Updated_at string     `json:"updated_at"`
	User       User       `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.tags, p.created_at, u.username, COUNT(c.id) as comments_count, p.version
	FROM posts p
	LEFT JOIN comments c ON p.id = c.post_id
	INNER JOIN users u ON u.id = p.user_id
	LEFT JOIN followers f ON f.follower_id = p.user_id AND f.user_id = $1
	WHERE (f.user_id IS NOT NULL OR p.user_id = $1)
	  AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%' OR p.tags::text ILIKE '%' || $4 || '%')
	  AND (p.tags @> $5 OR cardinality($5) = 0 OR $5 IS NULL)
	GROUP BY p.id, u.username
	ORDER BY p.created_at ` + fq.Sort + `
	LIMIT $2
	OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	var feed []PostWithMetadata

	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(

			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
			&post.Created_at,
			&post.User.Username,
			&post.CommentsCount,
			&post.Version,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, post)

	}

	return feed, nil

}

func (s *PostsStore) GetPostsByUserID(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.tags, p.created_at, u.username, COUNT(c.id) as comments_count, p.version
	FROM posts p
	LEFT JOIN comments c ON p.id = c.post_id
	INNER JOIN users u ON u.id = p.user_id
	WHERE p.user_id = $1
	GROUP BY p.id, u.username
	ORDER BY p.created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
			&post.Created_at,
			&post.User.Username,
			&post.CommentsCount,
			&post.Version,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT into posts (title, content, user_id, tags) VALUES ($1,$2,$3,$4) returning id, created_at, updated_at `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.Created_at,
		&post.Updated_at,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetPostByID(ctx context.Context, id int64) (*Post, error) {

	query := `SELECT p.id, p.title, p.content, p.user_id, p.tags, p.created_at, p.updated_at, p.version, u.id, u.username
	FROM posts p
	INNER JOIN users u ON u.id = p.user_id
	WHERE p.id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	post := new(Post)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.Created_at,
		&post.Updated_at,
		&post.Version,
		&post.User.ID,
		&post.User.Username,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE 
	FROM posts
	WHERE id = $1 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	row, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := row.RowsAffected()
	if err != nil {
		return nil
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 AND version = $4 returning version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(
		&post.Version,
	)
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
