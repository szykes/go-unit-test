package testx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertError(t *testing.T, expectedErr any, gotErr error, tcName string) {
	t.Helper()

	switch wantErr := expectedErr.(type) {
	case nil:
		assert.NoError(t, gotErr, tcName)
	case error:
		assert.ErrorIs(t, gotErr, wantErr, tcName)
	case string:
		assert.ErrorContains(t, gotErr, wantErr, tcName)
	default:
		t.Fatalf("unhandled type: %T at %v", expectedErr, tcName)
	}
}
