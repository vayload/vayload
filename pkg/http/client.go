/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - pkg/http-client
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/vayload/vayload/pkg/crypto"
)

type RequestInterceptor func(req *http.Request) error
type ResponseInterceptor func(res *http.Response, req *http.Request) error
type EventHandler func(id string, req *http.Request, res *http.Response, err error)
type EventName string

const (
	EventRequest  EventName = "request"
	EventResponse EventName = "response"
	EventError    EventName = "error"
	EventAlways   EventName = "always"
)

type BasicAuth struct {
	Username string
	Password string
}

type HttpClientConfig struct {
	BaseURL   string
	Timeout   time.Duration
	Headers   map[string]string
	Params    map[string]string
	Auth      *BasicAuth
	Transport http.RoundTripper
}

type clientConfig struct {
	baseURL string
	headers map[string]string // read-only
	params  map[string]string // read-only
	auth    *BasicAuth
	timeout time.Duration
}

type interceptorChain struct {
	request  []RequestInterceptor
	response []ResponseInterceptor
}

type eventRegistry struct {
	handlers map[EventName][]EventHandler
}

type HttpClient struct {
	cfg  clientConfig // inmutable
	http *http.Client

	mu sync.RWMutex

	interceptors atomic.Pointer[interceptorChain]
	events       atomic.Pointer[eventRegistry]

	lastRequestAt atomic.Int64
}

func NewHttpClient(config ...HttpClientConfig) *HttpClient {
	cfg := clientConfig{
		timeout: 30 * time.Second,
		headers: make(map[string]string),
	}

	if len(config) > 0 {
		c := config[0]
		cfg.baseURL = c.BaseURL
		if c.Timeout > 0 {
			cfg.timeout = c.Timeout
		}
		if c.Headers != nil {
			maps.Copy(cfg.headers, c.Headers)
		}
		if c.Params != nil {
			cfg.params = make(map[string]string)
			maps.Copy(cfg.params, c.Params)
		}
		cfg.auth = c.Auth
	}

	httpClient := &http.Client{Timeout: cfg.timeout}
	if len(config) > 0 && config[0].Transport != nil {
		httpClient.Transport = config[0].Transport
	}

	c := &HttpClient{cfg: cfg, http: httpClient}

	c.interceptors.Store(&interceptorChain{})
	c.events.Store(&eventRegistry{handlers: make(map[EventName][]EventHandler)})

	return c
}

func (c *HttpClient) UseRequest(fn RequestInterceptor) {
	c.mu.Lock()
	defer c.mu.Unlock()

	old := c.interceptors.Load()
	next := &interceptorChain{
		request:  append(append([]RequestInterceptor{}, old.request...), fn),
		response: old.response,
	}
	c.interceptors.Store(next)
}

func (c *HttpClient) UseResponse(fn ResponseInterceptor) {
	c.mu.Lock()
	defer c.mu.Unlock()

	old := c.interceptors.Load()
	next := &interceptorChain{
		request:  old.request,
		response: append(append([]ResponseInterceptor{}, old.response...), fn),
	}
	c.interceptors.Store(next)
}

func (c *HttpClient) Subscribe(name EventName, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	old := c.events.Load()
	newHandlers := make(map[EventName][]EventHandler, len(old.handlers)+1)
	maps.Copy(newHandlers, old.handlers)
	handlers := append([]EventHandler{}, old.handlers[name]...)
	handlers = append(handlers, handler)
	newHandlers[name] = handlers

	c.events.Store(&eventRegistry{handlers: newHandlers})
}

func (c *HttpClient) publish(name EventName, id string, req *http.Request, res *http.Response, err error) {
	reg := c.events.Load()

	for _, h := range reg.handlers[name] {
		h(id, req, res, err)
	}
}

func (c *HttpClient) Get(path string) *RequestBuilder {
	return c.NewRequest(http.MethodGet, path)
}

func (c *HttpClient) Post(path string) *RequestBuilder {
	return c.NewRequest(http.MethodPost, path)
}

func (c *HttpClient) Put(path string) *RequestBuilder {
	return c.NewRequest(http.MethodPut, path)
}

func (c *HttpClient) Delete(path string) *RequestBuilder {
	return c.NewRequest(http.MethodDelete, path)
}

func (c *HttpClient) Patch(path string) *RequestBuilder {
	return c.NewRequest(http.MethodPatch, path)
}

type RequestBuilder struct {
	client  *HttpClient
	method  string
	path    string
	headers map[string]string
	params  map[string]string
	body    any
	err     error
}

func (c *HttpClient) NewRequest(method, path string) *RequestBuilder {
	return &RequestBuilder{
		client:  c,
		method:  method,
		path:    path,
		headers: make(map[string]string),
		params:  make(map[string]string),
	}
}

func (r *RequestBuilder) Header(k, v string) *RequestBuilder { r.headers[k] = v; return r }

func (r *RequestBuilder) Headers(kv map[string]string) *RequestBuilder {
	for k, v := range kv {
		r.headers[k] = v
	}
	return r
}

func (r *RequestBuilder) Query(k, v string) *RequestBuilder { r.params[k] = v; return r }

func (r *RequestBuilder) JSON(v any) *RequestBuilder {
	r.body = v
	r.headers["Content-Type"] = "application/json"
	return r
}

func (r *RequestBuilder) Body(v any) *RequestBuilder {
	r.body = v
	return r
}

func (r *RequestBuilder) ContentType(ct string) *RequestBuilder {
	r.headers["Content-Type"] = ct

	return r
}

func (r *RequestBuilder) Multipart(v map[string]any) *RequestBuilder {
	r.body = v
	r.headers["Content-Type"] = "multipart/form-data"
	return r
}

func (r *RequestBuilder) URLEncoded(v url.Values) *RequestBuilder {
	r.body = v
	r.headers["Content-Type"] = "application/x-www-form-urlencoded"
	return r
}

func (r *RequestBuilder) Send(ctx context.Context) (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.client.cfg.timeout)
		defer cancel()
	}

	rawURL := buildURL(r.client.cfg.baseURL, r.path)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	if len(r.client.cfg.params) > 0 || len(r.params) > 0 {
		q := parsedURL.Query()
		for k, v := range r.client.cfg.params {
			q.Set(k, v)
		}
		for k, v := range r.params {
			q.Set(k, v)
		}
		parsedURL.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	var contentType string

	if r.body != nil {
		ct := r.headers["Content-Type"]
		if ct == "" {
			ct = "application/json"
		}
		bodyReader, contentType, err = getBodyReader(ct, r.body)
		if err != nil {
			return nil, fmt.Errorf("build body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, r.method, parsedURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	for k, v := range r.client.cfg.headers {
		req.Header.Set(k, v)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if r.client.cfg.auth != nil {
		req.SetBasicAuth(r.client.cfg.auth.Username, r.client.cfg.auth.Password)
	}

	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = crypto.GenerateNanoID()
		req.Header.Set("X-Request-ID", reqID)
	}

	chain := r.client.interceptors.Load()

	for _, fn := range chain.request {
		if err := fn(req); err != nil {
			r.client.publish(EventError, reqID, req, nil, err)
			return nil, err
		}
	}

	r.client.publish(EventRequest, reqID, req, nil, nil)
	r.client.lastRequestAt.Store(time.Now().UnixNano())

	var res *http.Response
	defer func() {
		r.client.publish(EventAlways, reqID, req, res, err)
	}()

	res, err = r.client.http.Do(req)
	if err != nil {
		r.client.publish(EventError, reqID, req, nil, err)
		return nil, err
	}

	for _, fn := range chain.response {
		if err := fn(res, req); err != nil {
			r.client.publish(EventError, reqID, req, res, err)
			return res, err
		}
	}

	r.client.publish(EventResponse, reqID, req, res, nil)
	return res, nil
}

func Do[T any](r *RequestBuilder, ctx context.Context) (*APIResponse[T], error) {
	res, err := r.Send(ctx)
	return UnwrapBody[T](res, err)
}

func (c *HttpClient) Clone() *HttpClient {
	headers := make(map[string]string, len(c.cfg.headers))
	maps.Copy(headers, c.cfg.headers)
	params := make(map[string]string, len(c.cfg.params))
	maps.Copy(params, c.cfg.params)

	var auth *BasicAuth
	if c.cfg.auth != nil {
		a := *c.cfg.auth
		auth = &a
	}

	clone := &HttpClient{
		cfg: clientConfig{
			baseURL: c.cfg.baseURL,
			headers: headers,
			params:  params,
			auth:    auth,
			timeout: c.cfg.timeout,
		},
		http: &http.Client{
			Timeout:   c.http.Timeout,
			Transport: c.http.Transport,
		},
	}

	chain := c.interceptors.Load()
	clone.interceptors.Store(&interceptorChain{
		request:  append([]RequestInterceptor{}, chain.request...),
		response: append([]ResponseInterceptor{}, chain.response...),
	})

	reg := c.events.Load()
	newHandlers := make(map[EventName][]EventHandler, len(reg.handlers))
	for k, v := range reg.handlers {
		newHandlers[k] = append([]EventHandler{}, v...)
	}
	clone.events.Store(&eventRegistry{handlers: newHandlers})

	return clone
}

func (c *HttpClient) Fork() *HttpClient {
	clone := c.Clone()
	clone.http = c.http
	return clone
}

func (c *HttpClient) LastRequestAt() time.Time {
	ns := c.lastRequestAt.Load()
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}

type APIResponse[T any] struct {
	Status   string         `json:"status"`
	Error    *APIError      `json:"error,omitempty"`
	Data     *T             `json:"data,omitempty"`
	Response *http.Response `json:"-"`
}

type APIError struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
}

type APIErrorCode string

func (f *APIErrorCode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*f = APIErrorCode(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*f = APIErrorCode(n.String())
		return nil
	}
	*f = ""
	return nil
}

type HttpClientErr struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *HttpClientErr) Error() string {
	return fmt.Sprintf("[%d] %s: %s (cause: %v)", e.Status, e.Code, e.Message, e.Cause)
}

func (e *HttpClientErr) Unwrap() error { return e.Cause }

func UnwrapBody[T any](response *http.Response, err error) (*APIResponse[T], error) {
	apiResponse := &APIResponse[T]{Response: response}
	if err != nil {
		return apiResponse, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(apiResponse); err != nil {
		return apiResponse, &HttpClientErr{
			Status:  response.StatusCode,
			Code:    "DECODE_ERROR",
			Message: "failed to decode response body",
			Cause:   err,
		}
	}

	if apiResponse.Error != nil {
		return apiResponse, &HttpClientErr{
			Status:  response.StatusCode,
			Code:    string(apiResponse.Error.Code),
			Message: apiResponse.Error.Message,
		}
	}

	return apiResponse, nil
}

func buildURL(base, path string) string {
	base = strings.TrimSuffix(base, "/")
	path = strings.TrimPrefix(path, "/")
	if base == "" {
		return path
	}
	return base + "/" + path
}

type FormFile struct {
	FileName string
	Content  io.Reader
}

type Stringable interface {
	String() string
}

type Encodable interface {
	Encode() string
}

func getBodyReader(contentType string, body any) (io.Reader, string, error) {
	ct := strings.TrimSpace(strings.Split(contentType, ";")[0])

	switch ct {
	case "application/x-www-form-urlencoded":
		switch v := body.(type) {
		case string:
			return strings.NewReader(v), ct, nil
		case Stringable:
			return strings.NewReader(v.String()), ct, nil
		case Encodable:
			return strings.NewReader(v.Encode()), ct, nil
		default:
			return nil, "", fmt.Errorf(
				"urlencoded body must be string or String(), got %T",
				body,
			)
		}

	case "multipart/form-data":
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		formData, ok := body.(map[string]any)
		if !ok {
			return nil, "", fmt.Errorf("multipart body must be map[string]any, got %T", body)
		}
		for key, val := range formData {
			if file, ok := val.(FormFile); ok {
				part, err := w.CreateFormFile(key, file.FileName)
				if err != nil {
					w.Close()
					return nil, "", err
				}
				if _, err := io.Copy(part, file.Content); err != nil {
					w.Close()
					return nil, "", err
				}
			} else {
				if err := w.WriteField(key, fmt.Sprintf("%v", val)); err != nil {
					w.Close()
					return nil, "", err
				}
			}
		}
		if err := w.Close(); err != nil {
			return nil, "", err
		}
		return buf, w.FormDataContentType(), nil
	case "application/json":
		b, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(b), "application/json; charset=utf-8", nil
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(b), "application/json; charset=utf-8", nil
	}
}
