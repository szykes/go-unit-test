package iam_test

import (
	"context"
	"testing"

	"github.com/szykes/go-unit-test/iam"
	"github.com/szykes/go-unit-test/iam/mock"
	"github.com/szykes/go-unit-test/recoverx"
	"github.com/szykes/go-unit-test/user"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mock/iam_mock.go -package=mock . identityProvider

type mocks struct {
	idp *mock.MockidentityProvider
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func bootstrap[TArg any, TWant any](
	t *testing.T,
	arg TArg,
	want TWant,
	prepare func(ctx context.Context, m *mocks, param TArg, want TWant),
) (context.Context, *iam.IAM) {
	defer recoverx.CatchPanicAndDebugPrint()

	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	ctx := t.Context()

	m := mocks{
		idp: mock.NewMockidentityProvider(ctrl),
	}

	if prepare != nil {
		prepare(ctx, &m, arg, want)
	}

	sut := iam.NewIAM(m.idp)

	return ctx, sut
}

func ListUsers_Succeeds(paramCtx context.Context, wantUsers []*user.User, idp *mock.MockidentityProvider) {
	idp.EXPECT().ListUsers(paramCtx).Return(wantUsers, nil)
}

func ListUsers_Fails(paramCtx context.Context, wantErr error, idp *mock.MockidentityProvider) {
	idp.EXPECT().ListUsers(paramCtx).Return(nil, wantErr)
}
