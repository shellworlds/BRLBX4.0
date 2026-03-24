package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDev(t *testing.T) {
	l, err := New(true)
	require.NoError(t, err)
	require.NotNil(t, l)
	_ = l.Sync()
}

func TestNewProd(t *testing.T) {
	l, err := New(false)
	require.NoError(t, err)
	require.NotNil(t, l)
	_ = l.Sync()
}
