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
	mustWrite := func(t *testing.T, w http.ResponseWriter, b []byte) {
		t.Helper()
		if _, err := w.Write(b); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}

	cases := []struct {
		name         string
		handler      func(t *testing.T) http.HandlerFunc
		ctxFn        func() (context.Context, context.CancelFunc)
		useNilClient bool
		assert       func(t *testing.T, res GitmojiResponse, err error)
	}{
		{
			name: "success",
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					mustWrite(t, w, []byte(
						`{"gitmojis":[{"emoji":"🎉","code":":tada:","name":"tada"}]}`,
					))
				}
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
			name: "nil client uses default",
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					mustWrite(t, w, []byte(`{"gitmojis":[{"emoji":"🔧","code":":wrench:","name":"wrench"}]}`))
				}
			},
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			useNilClient: true,
			assert: func(t *testing.T, res GitmojiResponse, err error) {
				if err != nil {
					t.Fatalf("Fetch with nil client returned error: %v", err)
				}
				if len(res.Gitmojis) != 1 || res.Gitmojis[0].Name != "wrench" {
					t.Fatalf("unexpected result: %#v", res)
				}
			},
		},
		{
			name: "server error",
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
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
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(200 * time.Millisecond)
					mustWrite(t, w, []byte(`{"gitmojis":[]}`))
				}
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
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					mustWrite(t, w, []byte(`{"gitmojis":[`))
				}
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
			server := httptest.NewServer(tc.handler(t))
			defer server.Close()

			ctx, cancel := tc.ctxFn()
			defer cancel()

			var client *http.Client
			if !tc.useNilClient {
				client = &http.Client{}
			}

			res, err := Fetch(ctx, client, server.URL)
			tc.assert(t, res, err)
		})
	}
}

func TestFindByCode(t *testing.T) {
	t.Parallel()

	gitmojis := []Gitmoji{
		{Emoji: "🎉", Code: ":tada:", Name: "tada", Description: "Begin a project"},
		{Emoji: "✨", Code: ":sparkles:", Name: "sparkles", Description: "Introduce new features"},
		{Emoji: "🐛", Code: ":bug:", Name: "bug", Description: "Fix a bug"},
	}

	tests := []struct {
		name      string
		code      string
		wantFound bool
		wantName  string
	}{
		{"found", ":tada:", true, "tada"},
		{"found sparkles", ":sparkles:", true, "sparkles"},
		{"not found", ":unknown:", false, ""},
		{"empty code", "", false, ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, found := FindByCode(tt.code, gitmojis)
			if found != tt.wantFound {
				t.Fatalf("found=%v, want %v", found, tt.wantFound)
			}
			if found && got.Name != tt.wantName {
				t.Fatalf("name=%s, want %s", got.Name, tt.wantName)
			}
		})
	}
}

func TestFindByName(t *testing.T) {
	t.Parallel()

	gitmojis := []Gitmoji{
		{Emoji: "🎉", Code: ":tada:", Name: "tada", Description: "Begin a project"},
		{Emoji: "✨", Code: ":sparkles:", Name: "sparkles", Description: "Introduce new features"},
	}

	tests := []struct {
		name      string
		lookup    string
		wantFound bool
		wantCode  string
	}{
		{"found", "tada", true, ":tada:"},
		{"not found", "unknown", false, ""},
		{"empty name", "", false, ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, found := FindByName(tt.lookup, gitmojis)
			if found != tt.wantFound {
				t.Fatalf("found=%v, want %v", found, tt.wantFound)
			}
			if found && got.Code != tt.wantCode {
				t.Fatalf("code=%s, want %s", got.Code, tt.wantCode)
			}
		})
	}
}

func TestGetAllGitmojisWithConfig_OfflineMode(t *testing.T) {
	t.Parallel()

	_, err := GetAllGitmojisWithConfig(true)
	if err == nil {
		t.Fatal("expected error for offline mode")
	}
	if err.Error() != "offline mode is enabled" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestFindByCode_EmptySlice(t *testing.T) {
	t.Parallel()
	_, found := FindByCode(":tada:", []Gitmoji{})
	if found {
		t.Fatal("expected not found for empty slice")
	}
}

func TestFindByName_EmptySlice(t *testing.T) {
	t.Parallel()
	_, found := FindByName("tada", []Gitmoji{})
	if found {
		t.Fatal("expected not found for empty slice")
	}
}
