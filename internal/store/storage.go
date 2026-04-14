package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeDuration     = time.Second * 10
	ErrDuplicatedEmail    = errors.New("duplicated user email")
	ErrDuplicatedUsername = errors.New("duplicated usernames")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrNotActivated       = errors.New("user account is not activated")
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetUserByID(context.Context, int64) (*User, error)
		Delete(context.Context, int64) error
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(context.Context, string) error
		GetUserByEmail(context.Context, string) (*User, error)
		SearchByUsername(ctx context.Context, query string) ([]*User, error)
	}
	Posts interface {
		Create(context.Context, *Post) error
		GetPostByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
		GetPostsByUserID(context.Context, int64) ([]PostWithMetadata, error)
	}
	Comments interface {
		GetCommentsByPostID(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, userID, followerID int64) error
		Unfollow(ctx context.Context, userID, followerID int64) error
	}
	Roles interface {
		GetByName(ctx context.Context, slug string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:     &UsersStore{db},
		Posts:     &PostsStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()

}
