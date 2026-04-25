package utils

import (
	"reflect"
	"testing"
)

func TestContainsString(t *testing.T) {
	tests := []struct {
		name   string
		slice  []string
		target string
		want   bool
	}{
		{name: "found", slice: []string{"go", "lang", "tool"}, target: "lang", want: true},
		{name: "not found", slice: []string{"go", "lang", "tool"}, target: "rust", want: false},
		{name: "empty slice", slice: []string{}, target: "anything", want: false},
		{name: "duplicate entries", slice: []string{"dup", "dup"}, target: "dup", want: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ContainsString(tc.slice, tc.target)
			if got != tc.want {
				t.Fatalf("ContainsString(%v, %q) = %v, want %v", tc.slice, tc.target, got, tc.want)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	tests := []struct {
		name   string
		slice  []string
		target string
		want   []string
	}{
		{name: "remove middle", slice: []string{"a", "b", "c"}, target: "b", want: []string{"a", "c"}},
		{name: "remove duplicates", slice: []string{"x", "y", "x", "x"}, target: "x", want: []string{"y"}},
		{name: "target absent", slice: []string{"a", "b"}, target: "c", want: []string{"a", "b"}},
		{name: "empty input", slice: []string{}, target: "any", want: []string{}},
		{name: "nil input", slice: nil, target: "any", want: []string{}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			original := append([]string{}, tc.slice...)

			got := RemoveString(tc.slice, tc.target)

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("RemoveString(%v, %q) = %v, want %v", tc.slice, tc.target, got, tc.want)
			}
			if tc.slice != nil && !reflect.DeepEqual(tc.slice, original) {
				t.Fatalf("RemoveString mutated input: got %v, want %v", tc.slice, original)
			}

		})
	}
}

func TestFilterString(t *testing.T) {
	tests := []struct {
		name        string
		slice       []string
		condition   func(string) bool
		want        []string
		expectCalls int
	}{
		{
			name:  "keep long strings",
			slice: []string{"go", "gopher", "play", "lang"},
			condition: func(s string) bool {
				return len(s) > 3
			},
			want:        []string{"gopher", "play", "lang"},
			expectCalls: 4,
		},
		{
			name:  "all filtered out",
			slice: []string{"a", "b", "c"},
			condition: func(string) bool {
				return false
			},
			want:        []string{},
			expectCalls: 3,
		},
		{
			name:  "empty input",
			slice: []string{},
			condition: func(string) bool {
				return true
			},
			want:        []string{},
			expectCalls: 0,
		},
		{
			name:  "nil input",
			slice: nil,
			condition: func(string) bool {
				return true
			},
			want:        []string{},
			expectCalls: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			cond := func(s string) bool {
				calls++
				return tc.condition(s)
			}

			got := FilterString(tc.slice, cond)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("FilterString(%v) = %v, want %v", tc.slice, got, tc.want)
			}
			if calls != tc.expectCalls {
				t.Fatalf("FilterString condition call count = %d, want %d", calls, tc.expectCalls)
			}
		})
	}
}
