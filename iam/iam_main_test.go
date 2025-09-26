package iam_test

import (
	"context"
	"testing"

	"github.com/szykes/go-unit-test/iam"
	"github.com/szykes/go-unit-test/iam/mock"
	"github.com/szykes/go-unit-test/user"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mock/mock.go -package=mock . idper

type mocks struct {
	idp *mock.Mockidper
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func bootstrap[TParams any, TWants any](
	t *testing.T,
	param TParams,
	want TWants,
	prepare func(ctx context.Context, m *mocks, param TParams, want TWants),
) (context.Context, *iam.IAM) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()

	m := mocks{
		idp: mock.NewMockidper(ctrl),
	}

	if prepare != nil {
		prepare(ctx, &m, param, want)
	}

	sut := iam.NewIAM(m.idp)

	return ctx, sut
}

func ListUsers_Succeeds(paramCtx context.Context, wantUsers []*user.User, idp *mock.Mockidper) {
	idp.EXPECT().ListUsers(paramCtx).Return(wantUsers, nil)
}

func ListUsers_Fails(paramCtx context.Context, wantErr error, idp *mock.Mockidper) {
	idp.EXPECT().ListUsers(paramCtx).Return(nil, wantErr)
}
