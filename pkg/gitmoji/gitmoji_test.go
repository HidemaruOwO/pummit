package gitmoji

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		handler http.HandlerFunc
		ctxFn   func() (context.Context, context.CancelFunc)
		assert  func(t *testing.T, res GitmojiResponse, err error)
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(
					`{"gitmojis":[{"emoji":"🎉","code":":tada:","name":"tada"}]}`,
				))
			},
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			assert: func(t *testing.T, res GitmojiResponse, err error) {
				if err != nil {
					t.Fatalf("Fetch returned error: %v", err)
				}
				if len(res.Gitmojis) != 1 || res.Gitmojis[0].Name != "tada" {
					t.Fatalf("unexpected result: %#v", res)
				}
			},
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			assert: func(t *testing.T, res GitmojiResponse, err error) {
				if err == nil || !errors.Is(err, ErrNonOK) {
					t.Fatalf("expected ErrNonOK, got %v", err)
				}
			},
		},
		{
			name: "timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(200 * time.Millisecond)
				_, _ = w.Write([]byte(`{"gitmojis":[]}`))
			},
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 50*time.Millisecond)
			},
			assert: func(t *testing.T, res GitmojiResponse, err error) {
				if err == nil {
					t.Fatalf("expected timeout error")
				}
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("expected deadline exceeded, got %v", err)
				}
			},
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"gitmojis":[`))
			},
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			assert: func(t *testing.T, res GitmojiResponse, err error) {
				if err == nil || !errors.Is(err, ErrDecode) {
					t.Fatalf("expected ErrDecode, got %v", err)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			ctx, cancel := tc.ctxFn()
			defer cancel()

			client := &http.Client{}
			res, err := Fetch(ctx, client, server.URL)
			tc.assert(t, res, err)
		})
	}
}
