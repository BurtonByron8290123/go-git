package hash_test

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/hash"
)

func TestNewHash(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "simple string",
			input: "hello world",
		},
		{
			name:  "git object content",
			input: "blob 11\x00hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hash.NewHash(tt.input)
			if h.IsZero() && tt.input != "" {
				// Non-empty input should not produce zero hash in most cases
				// (collision is theoretically possible but practically impossible)
				t.Logf("got zero hash for non-empty input: %q", tt.input)
			}
		})
	}
}

func TestNewHasher(t *testing.T) {
	h := hash.NewHasher()
	if h == nil {
		t.Fatal("expected non-nil hasher")
	}

	// Write some data and verify we get a valid hash
	_, err := h.Write([]byte("blob 11\x00hello world"))
	if err != nil {
		t.Fatalf("unexpected error writing to hasher: %v", err)
	}

	result := h.Sum(nil)
	if len(result) == 0 {
		t.Fatal("expected non-empty hash result")
	}
}

func TestHashConsistency(t *testing.T) {
	// Hashing the same input twice should produce the same result
	input := "tree 0\x00"

	h1 := hash.NewHash(input)
	h2 := hash.NewHash(input)

	if h1 != h2 {
		t.Errorf("expected identical hashes for identical input, got %v and %v", h1, h2)
	}
}

func TestHashUniqueness(t *testing.T) {
	// Different inputs should (with overwhelming probability) produce different hashes
	inputs := []string{
		"blob 5\x00hello",
		"blob 5\x00world",
		"blob 3\x00foo",
		"blob 3\x00bar",
	}

	hashes := make(map[interface{}]string)
	for _, input := range inputs {
		h := hash.NewHash(input)
		if existing, ok := hashes[h]; ok {
			t.Errorf("hash collision between %q and %q", existing, input)
		}
		hashes[h] = input
	}
}

func TestHasherWriteMultiple(t *testing.T) {
	// Writing in multiple chunks should produce the same result as writing all at once
	data := []byte("commit 100\x00tree abc\nauthor foo\n\nInitial commit\n")

	// Single write
	h1 := hash.NewHasher()
	_, _ = h1.Write(data)
	sum1 := h1.Sum(nil)

	// Multiple writes
	h2 := hash.NewHasher()
	for _, b := range data {
		_, _ = h2.Write([]byte{b})
	}
	sum2 := h2.Sum(nil)

	if len(sum1) != len(sum2) {
		t.Fatalf("hash length mismatch: %d vs %d", len(sum1), len(sum2))
	}
	for i := range sum1 {
		if sum1[i] != sum2[i] {
			t.Errorf("hash mismatch at byte %d: single=%x multi=%x", i, sum1, sum2)
			break
		}
	}
}
