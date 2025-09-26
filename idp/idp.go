package idp

import (
	"context"
	"errors"
	"fmt"

	"github.com/szykes/go-unit-test/user"
)

type IDPClient struct{}

func NewIDPClient() *IDPClient {
	return &IDPClient{}
}

func (a *IDPClient) FetchUserByID(ctx context.Context, userID string) (*user.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("fetch user failed: %w", errors.New("userID cannot be empty"))
	}

	switch userID {
	case "2ef33f49-4832-4866-a754-2602c0e67417":
		return &user.User{
			ID:       "2ef33f49-4832-4866-a754-2602c0e67417",
			Username: "John Doe",
			Email:    "john@doe.com",
		}, nil
	case "e2091e80-6e31-4824-b5bc-301ec166b357":
		return &user.User{
			ID:       "e2091e80-6e31-4824-b5bc-301ec166b357",
			Username: "Jane Smith",
			Email:    "jane@smith.com",
		}, nil
	}

	return nil, user.ErrUserNotFound
}

func (a *IDPClient) ListUsers(ctx context.Context) ([]*user.User, error) {
	return []*user.User{
		{
			ID:       "2ef33f49-4832-4866-a754-2602c0e67417",
			Username: "John Doe",
			Email:    "john@doe.com",
		},
		{
			ID:       "e2091e80-6e31-4824-b5bc-301ec166b357",
			Username: "Jane Smith",
			Email:    "jane@smith.com",
		},
	}, nil
}
