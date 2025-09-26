package iam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/szykes/go-unit-test/errorc"
	"github.com/szykes/go-unit-test/user"
)

func TestUserByID(t *testing.T) {
	t.Parallel()

	type params struct {
		userID string
	}

	type wants struct {
		user *user.User
		err  error
	}

	tcs := []struct {
		name    string
		param   params
		prepare func(ctx context.Context, m *mocks, param params, want wants)
		want    wants
	}{
		{
			name: "user has found",
			param: params{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				m.idp.EXPECT().FetchUser(ctx, param.userID).Return(want.user, nil)
			},
			want: wants{
				user: &user.User{
					ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
					Username: "John Doe",
					Email:    "john@doe.com",
				},
			},
		},
		{
			name: "FetchUser general error",
			param: params{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				m.idp.EXPECT().FetchUser(ctx, param.userID).Return(nil, want.err)
			},
			want: wants{
				err: errors.New("failed to get fetch user"),
			},
		},
		{
			name: "FetchUser user not found error",
			param: params{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				m.idp.EXPECT().FetchUser(ctx, param.userID).Return(nil, user.ErrUserNotFound)
			},
			want: wants{
				err: errorc.ErrNotFound,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.param, tc.want, tc.prepare)

			gotUser, gotErr := sut.UserByID(ctx, tc.param.userID)

			assert.Equal(t, tc.want.user, gotUser, tc.name)
			assert.ErrorIs(t, gotErr, tc.want.err, tc.name)
		})
	}
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	type params struct {
	}

	type wants struct {
		users []*user.User
		err   error
	}

	tcs := []struct {
		name    string
		param   params
		prepare func(ctx context.Context, m *mocks, param params, want wants)
		want    wants
	}{
		{
			name:  "users have found",
			param: params{},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				m.idp.EXPECT().ListUsers(ctx).Return(want.users, nil)
			},
			want: wants{
				users: []*user.User{
					{
						ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
						Username: "John Doe",
						Email:    "john@doe.com",
					},
					{
						ID:       "481fa333-9b6d-4a6d-b4fd-da95535ed436",
						Username: "Jane Smith",
						Email:    "jane@smith.com",
					},
				},
			},
		},
		{
			name:  "ListUsers general error",
			param: params{},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				m.idp.EXPECT().ListUsers(ctx).Return(nil, want.err)
			},
			want: wants{
				err: errors.New("failed to list users"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.param, tc.want, tc.prepare)

			gotUsers, gotErr := sut.ListUsers(ctx)

			assert.Equal(t, tc.want.users, gotUsers, tc.name)
			assert.ErrorIs(t, gotErr, tc.want.err, tc.name)
		})
	}
}

func TestUserByEmail(t *testing.T) {
	t.Parallel()

	type params struct {
		email string
	}

	type wants struct {
		user *user.User
		err  error
	}

	tcs := []struct {
		name    string
		param   params
		prepare func(ctx context.Context, m *mocks, param params, want wants)
		want    wants
	}{
		{
			name: "user has found",
			param: params{
				email: "jane@smith.com",
			},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				users := []*user.User{
					{
						ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
						Username: "John Doe",
						Email:    "john@doe.com",
					},
					want.user,
				}
				ListUsers_Succeeds(ctx, users, m.idp)
			},
			want: wants{
				user: &user.User{
					ID:       "481fa333-9b6d-4a6d-b4fd-da95535ed436",
					Username: "Jane Smith",
					Email:    "jane@smith.com",
				},
			},
		},
		{
			name: "user has NOT found",
			param: params{
				email: "eef485ab-483b-4a56-81ae-99e8e60075c6",
			},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				users := []*user.User{
					{
						ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
						Username: "John Doe",
						Email:    "john@doe.com",
					},
					{
						ID:       "481fa333-9b6d-4a6d-b4fd-da95535ed436",
						Username: "Jane Smith",
						Email:    "jane@smith.com",
					},
				}
				ListUsers_Succeeds(ctx, users, m.idp)
			},
			want: wants{
				err: errorc.ErrNotFound,
			},
		},
		{
			name:  "ListUsers general error",
			param: params{},
			prepare: func(ctx context.Context, m *mocks, param params, want wants) {
				ListUsers_Fails(ctx, want.err, m.idp)
			},
			want: wants{
				err: errors.New("failed to list users"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.param, tc.want, tc.prepare)

			gotUser, gotErr := sut.UserByEmail(ctx, tc.param.email)

			assert.Equal(t, tc.want.user, gotUser, tc.name)
			assert.ErrorIs(t, gotErr, tc.want.err, tc.name)
		})
	}
}
