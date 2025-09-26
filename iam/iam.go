package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/szykes/go-unit-test/errorc"
	"github.com/szykes/go-unit-test/user"
)

type idper interface {
	FetchUser(ctx context.Context, userID string) (*user.User, error)
	ListUsers(ctx context.Context) ([]*user.User, error)
}

type IAM struct {
	idp idper
}

func NewIAM(idp idper) *IAM {
	return &IAM{
		idp: idp,
	}
}

func (i *IAM) UserByID(ctx context.Context, userID string) (*user.User, error) {
	userEntity, err := i.idp.FetchUser(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, fmt.Errorf("user by ID: %w", errorc.ErrNotFound)
		}

		return nil, fmt.Errorf("user by ID: %w", err)
	}

	return userEntity, nil
}

func (i *IAM) ListUsers(ctx context.Context) ([]*user.User, error) {
	userEntities, err := i.idp.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list idp users: %w", err)
	}

	return userEntities, nil
}

func (i *IAM) UserByEmail(ctx context.Context, userEmail string) (*user.User, error) {
	userEntities, err := i.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	for _, userEntity := range userEntities {
		if userEntity.Email == userEmail {
			return userEntity, nil
		}
	}

	return nil, errorc.ErrNotFound
}
