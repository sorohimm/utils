package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.wildberries.ru/balance-pay/balance-pay/mobile/internal/utils/log"
)

var Err404 = errors.New("404 not found")

type TransportOpts struct {
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	Timeout             time.Duration
}

func DefaultTransport() *TransportOpts {
	return &TransportOpts{
		MaxIdleConns:        100,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000,
		Timeout:             time.Second * 10,
	}
}

func NewDefaultHTTPClient(opts *TransportOpts) (*http.Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("nil opts")
	}

	var t *http.Transport
	tt, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		t = tt.Clone()
	}
	if t == nil {
		return nil, fmt.Errorf("nil transport")
	}
	t.MaxIdleConns = opts.MaxIdleConns
	t.MaxConnsPerHost = opts.MaxConnsPerHost
	t.MaxIdleConnsPerHost = opts.MaxIdleConnsPerHost

	return &http.Client{
		Transport: t,
		Timeout:   opts.Timeout,
	}, nil
}

type BaseClientParams struct {
	BaseURI string
	Before  func(r *http.Request)
	After   func(r *http.Response)
}

func NewBaseClient(httpClient *http.Client, p *BaseClientParams) *BaseClient {
	return &BaseClient{
		Client: httpClient,
		now:    time.Now,

		baseURI: p.BaseURI,
	}
}

type BaseClient struct {
	*http.Client
	baseURI string
	now     func() time.Time

	parallelLimit int

	before func(r *http.Request)
	after  func(r *http.Response)
}

func (o *BaseClient) BaseURI() string {
	return o.baseURI
}

func (o *BaseClient) GetBytes(ctx context.Context, req *http.Request, opts ...Option) ([]byte, error) {
	var (
		bb     = bytes.NewBuffer(nil)
		logger = log.FromContext(ctx).Sugar()
		res    *http.Response
		err    error
	)

	res, err = doRequest(o.Client, req, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Errorf("failed to close response body:%v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			// maybe there is some useful info in response
			_, _ = io.Copy(bb, res.Body)
			if bb.Len() > 0 {
				logger.Errorf("404 error while doing request: %v : %s", req.URL, bb.String())
			}
			return bb.Bytes(), Err404
		}

		return nil, fmt.Errorf("http status: %d", res.StatusCode)
	}
	if _, err = io.Copy(bb, res.Body); err != nil {
		return nil, fmt.Errorf("failed to copy body: %v", err)
	}

	return bb.Bytes(), nil
}

type StatusError struct {
	Expected int
	Actual   int
}

func (o StatusError) Error() string {
	return fmt.Sprintf("the status %d was expected but the status %d was received", o.Expected, o.Actual)
}

func (o *BaseClient) DoRequest(req *http.Request, expectedStatus int, opts ...Option) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	resp, err = doRequest(o.Client, req, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		return resp, StatusError{
			Expected: expectedStatus,
			Actual:   resp.StatusCode,
		}
	}

	return resp, nil
}

func doRequest(cl *http.Client, req *http.Request, opts ...Option) (*http.Response, error) {
	o := EvaluateOptions(req, opts)

	if cl == nil {
		cl = http.DefaultClient
	}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	o.ExecAfter(resp)

	return resp, nil
}
