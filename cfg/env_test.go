package cfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDotEnv(t *testing.T) {
	t.Run("testLoad", func(t *testing.T) {
		testLoad(t)
	})
	t.Run("testDotEnvDoesNotOverrideEnv", func(t *testing.T) {
		testDotEnvDoesNotOverrideEnv1(t)
	})
}

func testLoad(t *testing.T) {
	t.Helper()
	l := NewLookup(".env.test", 1)
	file, err := l.FindFile()
	require.NoError(t, err)

	err = loadEnv(file)
	require.NoError(t, err)
	require.Equal(t, "value1", os.Getenv("ENV_VAR_1"))
	require.Equal(t, "https://localhost:3000/path", os.Getenv("ENV_VAR_2"))
	require.Equal(t, "10", os.Getenv("ENV_VAR_3"))
	require.Equal(t, "123_456", os.Getenv("ENV_VAR_4"))
}

func testDotEnvDoesNotOverrideEnv1(t *testing.T) {
	t.Helper()
	t.Setenv("ENV_VAR_1", "value123")

	l := NewLookup(".env.test", 1)
	file, err := l.FindFile()
	require.NoError(t, err)

	err = loadEnv(file)
	require.NoError(t, err)
	require.Equal(t, "value123", os.Getenv("ENV_VAR_1"))
}
