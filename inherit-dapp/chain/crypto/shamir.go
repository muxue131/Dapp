package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// gf256Mul multiplies two elements in GF(2^8) with irreducible polynomial 0x11B
func gf256Mul(a, b byte) byte {
	var p byte
	for i := 0; i < 8; i++ {
		if b&1 != 0 {
			p ^= a
		}
		hi := a & 0x80
		a <<= 1
		if hi != 0 {
			a ^= 0x1B
		}
		b >>= 1
	}
	return p
}

// gf256Pow computes a^n in GF(2^8) using square-and-multiply
func gf256Pow(a byte, n int) byte {
	if n == 0 {
		return 1
	}
	result := byte(1)
	base := a
	for n > 0 {
		if n&1 != 0 {
			result = gf256Mul(result, base)
		}
		base = gf256Mul(base, base)
		n >>= 1
	}
	return result
}

// gf256Inverse computes the multiplicative inverse in GF(2^8)
func gf256Inverse(a byte) byte {
	if a == 0 {
		return 0
	}
	return gf256Pow(a, 254)
}

// Share represents a single Shamir secret share
type Share struct {
	Index int    `json:"index"`
	Value []byte `json:"value"`
}

// ShamirSecretSharing implements Shamir's Secret Sharing scheme over GF(2^8)
type ShamirSecretSharing struct{}

// NewShamirSecretSharing creates a new Shamir secret sharing instance
func NewShamirSecretSharing() *ShamirSecretSharing {
	return &ShamirSecretSharing{}
}

// Split splits a secret into n shares where k shares are needed to reconstruct
func (s *ShamirSecretSharing) Split(secret []byte, n, k int) ([]Share, error) {
	if k > n {
		return nil, errors.New("threshold k cannot be greater than n")
	}
	if k < 2 {
		return nil, errors.New("threshold k must be at least 2")
	}
	if n < 2 {
		return nil, errors.New("number of shares n must be at least 2")
	}
	if n > 255 {
		return nil, errors.New("number of shares n must be at most 255")
	}
	if len(secret) == 0 {
		return nil, errors.New("secret cannot be empty")
	}

	shares := make([]Share, n)
	for i := range shares {
		shares[i] = Share{
			Index: i + 1,
			Value: make([]byte, len(secret)),
		}
	}

	for byteIdx, secretByte := range secret {
		coeffs := make([]byte, k)
		coeffs[0] = secretByte

		for i := 1; i < k; i++ {
			for {
				r, err := rand.Int(rand.Reader, big.NewInt(256))
				if err != nil {
					return nil, fmt.Errorf("failed to generate random coefficient: %w", err)
				}
				coeffs[i] = byte(r.Int64())
				// Ensure the leading coefficient (highest degree) is non-zero
				// so the polynomial has degree exactly k-1
				if i != k-1 || coeffs[i] != 0 {
					break
				}
			}
		}

		for i := 0; i < n; i++ {
			x := byte(i + 1)
			y := evalPolyGF256(coeffs, x)
			shares[i].Value[byteIdx] = y
		}
	}

	return shares, nil
}

// evalPolyGF256 evaluates polynomial with coefficients in GF(2^8) using Horner's method
func evalPolyGF256(coeffs []byte, x byte) byte {
	result := coeffs[len(coeffs)-1]
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = gf256Mul(result, x) ^ coeffs[i]
	}
	return result
}

// Combine combines shares to reconstruct the secret using Lagrange interpolation at x=0
func (s *ShamirSecretSharing) Combine(shares []Share) ([]byte, error) {
	if len(shares) < 2 {
		return nil, errors.New("need at least 2 shares to reconstruct")
	}

	secretLen := len(shares[0].Value)
	for _, share := range shares[1:] {
		if len(share.Value) != secretLen {
			return nil, errors.New("all shares must have the same length")
		}
	}

	secret := make([]byte, secretLen)

	for byteIdx := 0; byteIdx < secretLen; byteIdx++ {
		var result byte

		for i, share := range shares {
			xi := byte(share.Index)
			var num, den byte = 1, 1

			for j, other := range shares {
				if i == j {
					continue
				}
				xj := byte(other.Index)
				num = gf256Mul(num, xj)
				den = gf256Mul(den, xi^xj)
			}

			denInv := gf256Inverse(den)
			y := share.Value[byteIdx]
			term := gf256Mul(gf256Mul(y, num), denInv)
			result ^= term
		}

		secret[byteIdx] = result
	}

	return secret, nil
}

// SplitKey splits an encryption key into shares for distribution to beneficiaries
func SplitKey(key []byte, numBeneficiaries, threshold int) ([]Share, error) {
	sss := NewShamirSecretSharing()
	return sss.Split(key, numBeneficiaries, threshold)
}

// CombineKey reconstructs the encryption key from shares
func CombineKey(shares []Share) ([]byte, error) {
	sss := NewShamirSecretSharing()
	return sss.Combine(shares)
}

// DistributeKeyShares creates key shares for each beneficiary
func DistributeKeyShares(masterKey []byte, numBeneficiaries, threshold int) (map[int]Share, error) {
	shares, err := SplitKey(masterKey, numBeneficiaries, threshold)
	if err != nil {
		return nil, err
	}

	result := make(map[int]Share, len(shares))
	for _, share := range shares {
		result[share.Index] = share
	}
	return result, nil
}
