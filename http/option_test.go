package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://one.two/free/four", nil)
	require.NoError(t, err)

	opts := []Option{
		WithContentType("some/type"),
		WithQueryParam("key", "value"),
		WithQueryParam("k", "v"),
	}
	EvaluateOptions(req, opts)

	require.Equal(t, []string{"some/type"}, req.Header["Content-Type"])
	require.Equal(t, "k=v&key=value", req.URL.RawQuery)
	// TODO more checks
}
