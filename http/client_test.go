package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBaseClient_GetBytes(t *testing.T) {
	// Test with a 200 response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer ts.Close()

	baseClient := &BaseClient{
		Client:  http.DefaultClient,
		baseURI: ts.URL,
		now:     time.Now,
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	b, err := baseClient.GetBytes(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to get bytes: %v", err)
	}
	if !bytes.Equal(b, []byte("OK")) {
		t.Fatalf("unexpected response: %s", string(b))
	}

	// Test with a 404 response
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	}))
	defer ts.Close()

	baseClient = &BaseClient{
		Client:  http.DefaultClient,
		baseURI: ts.URL,
		now:     time.Now,
	}

	req, err = http.NewRequestWithContext(context.Background(), "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	b, err = baseClient.GetBytes(context.Background(), req)
	if err != Err404 {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte("Not Found")) {
		t.Fatalf("unexpected response: %s", string(b))
	}

	// Test with a 500 response
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server ValidationError"))
	}))
	defer ts.Close()

	baseClient = &BaseClient{
		Client:  http.DefaultClient,
		baseURI: ts.URL,
		now:     time.Now,
	}

	req, err = http.NewRequestWithContext(context.Background(), "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	_, err = baseClient.GetBytes(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "http status: 500") {
		t.Fatalf("expected error with status code 500, got: %v", err)
	}
}

func TestBaseURI(t *testing.T) {
	baseClient := &BaseClient{baseURI: "https://example.com"}
	expectedURI := "https://example.com"

	if baseClient.BaseURI() != expectedURI {
		t.Errorf("Expected base URI to be %s, but got %s", expectedURI, baseClient.BaseURI())
	}
}

func TestDoRequestDefaultClientError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("failed to create http.Request: %v", err)
	}
	req.RequestURI = "invalid"

	_, got := doRequest(nil, req)
	if got == nil {
		t.Errorf("doRequest() error = %v, wantErr %v", got, "Request.RequestURI can't be set in client requests")
	}
}

func TestDoRequestCustomClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response body"))
	}))
	defer ts.Close()

	cl := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create http.Request: %v", err)
	}

	resp, err := doRequest(cl, req)
	if err != nil {
		t.Fatalf("doRequest() error = %v, wantErr %v", err, nil)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("doRequest() response status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
