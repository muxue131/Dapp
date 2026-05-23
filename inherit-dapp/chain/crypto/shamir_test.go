package crypto

import (
	"bytes"
	"testing"
)

func TestShamirSplitCombine(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("This is a secret key for inheritance!")

	shares, err := sss.Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	if len(shares) != 5 {
		t.Fatalf("Expected 5 shares, got %d", len(shares))
	}

	for i, share := range shares {
		if share.Index != i+1 {
			t.Fatalf("Share %d has wrong index: %d", i, share.Index)
		}
		if len(share.Value) != len(secret) {
			t.Fatalf("Share %d has wrong length: expected %d, got %d", i, len(secret), len(share.Value))
		}
	}
}

func TestShamirReconstructWithThreshold(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("Secret inheritance key")

	shares, err := sss.Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Reconstruct with exactly k shares
	recovered, err := sss.Combine(shares[:3])
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}

	if !bytes.Equal(secret, recovered) {
		t.Fatalf("Recovered secret doesn't match.\nExpected: %v\nGot: %v", secret, recovered)
	}
}

func TestShamirReconstructWithMoreThanThreshold(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("Another secret")

	shares, err := sss.Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Use 4 shares (more than threshold of 3)
	recovered, err := sss.Combine(shares[:4])
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}

	if !bytes.Equal(secret, recovered) {
		t.Fatal("Recovered secret doesn't match")
	}
}

func TestShamirReconstructWithAllShares(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("Test secret with all shares")

	shares, err := sss.Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	recovered, err := sss.Combine(shares)
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}

	if !bytes.Equal(secret, recovered) {
		t.Fatal("Recovered secret doesn't match")
	}
}

func TestShamirDifferentSubsets(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("Subset test secret")

	shares, err := sss.Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Try different subsets of 3 shares
	subsets := [][]int{
		{0, 1, 2},
		{0, 1, 3},
		{0, 1, 4},
		{0, 2, 3},
		{0, 2, 4},
		{0, 3, 4},
		{1, 2, 3},
		{1, 2, 4},
		{1, 3, 4},
		{2, 3, 4},
	}

	for _, subset := range subsets {
		selected := make([]Share, len(subset))
		for i, idx := range subset {
			selected[i] = shares[idx]
		}

		recovered, err := sss.Combine(selected)
		if err != nil {
			t.Fatalf("Combine with subset %v failed: %v", subset, err)
		}

		if !bytes.Equal(secret, recovered) {
			t.Fatalf("Subset %v produced wrong secret", subset)
		}
	}
}

func TestSplitKey(t *testing.T) {
	key, _ := GenerateKey()

	shares, err := SplitKey(key, 5, 3)
	if err != nil {
		t.Fatalf("SplitKey failed: %v", err)
	}

	if len(shares) != 5 {
		t.Fatalf("Expected 5 shares, got %d", len(shares))
	}
}

func TestCombineKey(t *testing.T) {
	key, _ := GenerateKey()

	shares, err := SplitKey(key, 5, 3)
	if err != nil {
		t.Fatalf("SplitKey failed: %v", err)
	}

	recovered, err := CombineKey(shares[:3])
	if err != nil {
		t.Fatalf("CombineKey failed: %v", err)
	}

	if !bytes.Equal(key, recovered) {
		t.Fatal("Recovered key doesn't match original")
	}
}

func TestDistributeKeyShares(t *testing.T) {
	key, _ := GenerateKey()

	shareMap, err := DistributeKeyShares(key, 4, 3)
	if err != nil {
		t.Fatalf("DistributeKeyShares failed: %v", err)
	}

	if len(shareMap) != 4 {
		t.Fatalf("Expected 4 shares, got %d", len(shareMap))
	}

	// Check all indices are present
	for i := 1; i <= 4; i++ {
		if _, ok := shareMap[i]; !ok {
			t.Fatalf("Missing share for index %d", i)
		}
	}
}

func TestShamirInvalidParams(t *testing.T) {
	sss := NewShamirSecretSharing()
	secret := []byte("test")

	// k > n
	_, err := sss.Split(secret, 3, 5)
	if err == nil {
		t.Fatal("Expected error when k > n")
	}

	// k < 2
	_, err = sss.Split(secret, 3, 1)
	if err == nil {
		t.Fatal("Expected error when k < 2")
	}

	// n < 2
	_, err = sss.Split(secret, 1, 1)
	if err == nil {
		t.Fatal("Expected error when n < 2")
	}

	// empty secret
	_, err = sss.Split([]byte{}, 3, 2)
	if err == nil {
		t.Fatal("Expected error for empty secret")
	}
}
