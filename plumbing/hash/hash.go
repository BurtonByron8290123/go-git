// Package hash provides hashing utilities for git objects.
// It supports multiple hash algorithms used in git repositories.
package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
)

// Algorithm represents a hashing algorithm used in git.
type Algorithm uint

const (
	// SHA1 is the default hashing algorithm used by git.
	SHA1 Algorithm = iota
	// SHA256 is the new hashing algorithm introduced in git 2.29.
	SHA256
)

// Hash represents a git object hash.
type Hash [32]byte

// ZeroHash is a hash with all bytes set to zero.
var ZeroHash Hash

// NewHash creates a new Hash from a hex string.
func NewHash(s string) Hash {
	b := []byte(s)
	var h Hash
	copy(h[:], b)
	return h
}

// String returns the hex representation of the hash.
func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

// IsZero returns true if the hash is the zero hash.
func (h Hash) IsZero() bool {
	return h == ZeroHash
}

// Hasher wraps a hash.Hash with the algorithm used.
type Hasher struct {
	hash.Hash
	algo Algorithm
}

// NewHasher creates a new Hasher for the given algorithm.
func NewHasher(algo Algorithm) Hasher {
	switch algo {
	case SHA256:
		return Hasher{Hash: sha256.New(), algo: SHA256}
	default:
		return Hasher{Hash: sha1.New(), algo: SHA1}
	}
}

// Sum returns the hash of all data written to the hasher.
func (h Hasher) Sum() Hash {
	var result Hash
	copy(result[:], h.Hash.Sum(nil))
	return result
}

// Algorithm returns the algorithm used by the hasher.
func (h Hasher) Algorithm() Algorithm {
	return h.algo
}

// String returns the string representation of the algorithm.
func (a Algorithm) String() string {
	switch a {
	case SHA256:
		return "sha256"
	default:
		return "sha1"
	}
}

// Size returns the byte size of the hash for the given algorithm.
func (a Algorithm) Size() int {
	switch a {
	case SHA256:
		return sha256.Size
	default:
		return sha1.Size
	}
}
