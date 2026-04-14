package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
		Roles: &MockRoleStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	return &User{ID: userID}, nil
}

func (m *MockUserStore) GetUserByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *MockUserStore) SearchByUsername(ctx context.Context, query string) ([]*User, error) {
	return nil, nil
}

func (m *MockUserStore) GetSuggestedUsers(ctx context.Context, userID int64) ([]*User, error) {
	return nil, nil
}

type MockRoleStore struct{}

func (m *MockRoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	return &Role{}, nil
}
