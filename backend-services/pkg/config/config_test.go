package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetHelpers(t *testing.T) {
	t.Setenv("COVERAGE_TEST_INT", "7")
	t.Setenv("COVERAGE_TEST_BOOL", "true")
	t.Setenv("COVERAGE_TEST_FLOAT", "2.5")
	Load()
	viper.AutomaticEnv()

	require.Equal(t, 7, GetInt("COVERAGE_TEST_INT", 1))
	require.True(t, GetBool("COVERAGE_TEST_BOOL", false))
	require.InDelta(t, 2.5, GetFloat64("COVERAGE_TEST_FLOAT", 0), 1e-9)
	require.Equal(t, 99, GetInt("COVERAGE_TEST_MISSING_INT", 99))
}
