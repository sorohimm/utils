package http

import (
	"net/http"
	"net/url"
)

type After struct {
	after []func(*http.Response)
}

func (o *After) Add(f func(*http.Response)) {
	o.after = append(o.after, f)
}

func (o *After) Exec(resp *http.Response) {
	for _, f := range o.after {
		f(resp)
	}
}

type Options struct {
	request   *http.Request
	urlValues *url.Values
	after     *After
}

func (o *Options) AddAfter(after func(*http.Response)) {
	o.after.Add(after)
}

func (o *Options) ExecAfter(resp *http.Response) {
	if o.after != nil {
		o.after.Exec(resp)
	}
}

type Option func(*Options)

func EvaluateOptions(req *http.Request, opts []Option) *Options {
	optsCopy := &Options{
		request:   req,
		urlValues: &url.Values{},
	}
	for _, opt := range opts {
		opt(optsCopy)
	}
	optsCopy.request.URL.RawQuery = optsCopy.urlValues.Encode()

	return optsCopy
}

func WithContentType(ct string) Option {
	return func(o *Options) {
		o.request.Header.Set("Content-Type", ct)
	}
}

func WithQueryParam(key, value string) Option {
	return func(o *Options) {
		o.urlValues.Add(key, value)
	}
}

func WithHeader(key, value string) Option {
	return func(o *Options) {
		o.request.Header.Set(key, value)
	}
}

func WithAfterFunc(f func(*http.Response)) Option {
	return func(o *Options) {
		o.AddAfter(f)
	}
}
