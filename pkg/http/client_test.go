/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientGET(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET")
		}

		if r.URL.Path != "/users" {
			t.Fatalf("wrong path")
		}

		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok","data":{"name":"john"}}`))
	}))
	defer s.Close()

	client := NewHttpClient(HttpClientConfig{
		BaseURL: s.URL,
	})

	res, err := client.Get("/users").Send(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatal("bad status")
	}
}

func TestQueryParams(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Fatal("missing query param")
		}

		w.WriteHeader(200)
	}))
	defer s.Close()

	client := NewHttpClient(HttpClientConfig{
		BaseURL: s.URL,
	})

	_, err := client.
		Get("/users").
		Query("limit", "10").
		Send(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONBody(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatal("wrong content type")
		}

		var buf bytes.Buffer
		io.Copy(&buf, r.Body)

		if !strings.Contains(buf.String(), "name") {
			t.Fatal("json body missing")
		}

		w.WriteHeader(200)
	}))
	defer s.Close()

	client := NewHttpClient(HttpClientConfig{
		BaseURL: s.URL,
	})

	_, err := client.
		Post("/users").
		JSON(map[string]string{"name": "john"}).
		Send(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}

func TestInterceptor(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "1" {
			t.Fatal("interceptor failed")
		}

		w.WriteHeader(200)
	}))
	defer s.Close()

	client := NewHttpClient(HttpClientConfig{
		BaseURL: s.URL,
	})

	client.UseRequest(func(req *http.Request) error {
		req.Header.Set("X-Test", "1")
		return nil
	})

	_, err := client.Get("/").Send(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}

func TestEvents(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer s.Close()

	client := NewHttpClient(HttpClientConfig{
		BaseURL: s.URL,
	})

	called := false

	client.Subscribe(EventResponse, func(id string, req *http.Request, res *http.Response, err error) {
		called = true
	})

	_, err := client.Get("/").Send(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("event not fired")
	}
}
