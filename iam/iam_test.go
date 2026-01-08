package iam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/szykes/go-unit-test/errorc"
	"github.com/szykes/go-unit-test/testx"
	"github.com/szykes/go-unit-test/user"
)

func TestUserByID(t *testing.T) {
	t.Parallel()

	type arg struct {
		userID string
	}

	type want struct {
		user *user.User
		err  error
	}

	tcs := []struct {
		name    string
		arg     arg
		prepare func(ctx context.Context, m *mocks, arg arg, want want)
		want    want
	}{
		{
			name: "user has found",
			arg: arg{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().FetchUser(ctx, arg.userID).Return(want.user, nil)
			},
			want: want{
				user: &user.User{
					ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
					Username: "John Doe",
					Email:    "john@doe.com",
				},
			},
		},
		{
			name: "FetchUser general error",
			arg: arg{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().FetchUser(ctx, arg.userID).Return(nil, want.err)
			},
			want: want{
				err: errors.New("failed to get fetch user"),
			},
		},
		{
			name: "FetchUser user not found error",
			arg: arg{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().FetchUser(ctx, arg.userID).Return(nil, user.ErrUserNotFound)
			},
			want: want{
				err: errorc.ErrNotFound,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.arg, tc.want, tc.prepare)

			gotUser, gotErr := sut.UserByID(ctx, tc.arg.userID)

			assert.Equal(t, tc.want.user, gotUser, tc.name)
			testx.AssertError(t, tc.want.err, gotErr, tc.name)
		})
	}
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	type arg struct {
	}

	type want struct {
		users []*user.User
		err   error
	}

	tcs := []struct {
		name    string
		arg     arg
		prepare func(ctx context.Context, m *mocks, arg arg, want want)
		want    want
	}{
		{
			name: "users have found",
			arg:  arg{},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().ListUsers(ctx).Return(want.users, nil)
			},
			want: want{
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
			name: "ListUsers general error",
			arg:  arg{},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().ListUsers(ctx).Return(nil, want.err)
			},
			want: want{
				err: errors.New("failed to list users"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.arg, tc.want, tc.prepare)

			gotUsers, gotErr := sut.ListUsers(ctx)

			assert.Equal(t, tc.want.users, gotUsers, tc.name)
			testx.AssertError(t, tc.want.err, gotErr, tc.name)
		})
	}
}

func TestUserByEmail(t *testing.T) {
	t.Parallel()

	type arg struct {
		email string
	}

	type want struct {
		user *user.User
		err  error
	}

	tcs := []struct {
		name    string
		arg     arg
		prepare func(ctx context.Context, m *mocks, arg arg, want want)
		want    want
	}{
		{
			name: "user has found",
			arg: arg{
				email: "jane@smith.com",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
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
			want: want{
				user: &user.User{
					ID:       "481fa333-9b6d-4a6d-b4fd-da95535ed436",
					Username: "Jane Smith",
					Email:    "jane@smith.com",
				},
			},
		},
		{
			name: "user has NOT found",
			arg: arg{
				email: "eef485ab-483b-4a56-81ae-99e8e60075c6",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
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
			want: want{
				err: errorc.ErrNotFound,
			},
		},
		{
			name: "ListUsers general error",
			arg:  arg{},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				ListUsers_Fails(ctx, want.err, m.idp)
			},
			want: want{
				err: errors.New("failed to list users"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.arg, tc.want, tc.prepare)

			gotUser, gotErr := sut.UserByEmail(ctx, tc.arg.email)

			assert.Equal(t, tc.want.user, gotUser, tc.name)
			testx.AssertError(t, tc.want.err, gotErr, tc.name)
		})
	}
}
