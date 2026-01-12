package core

import (
	"testing"
)

func TestDecodeEUCKR(t *testing.T) {
	// "삼성전자" encoded in EUC-KR (CP949)
	eucKRBytes := []byte{0xbb, 0xef, 0xbc, 0xba, 0xc0, 0xfc, 0xc0, 0xda}
	expected := "삼성전자"

	result, err := DecodeEUCKR(eucKRBytes)
	if err != nil {
		t.Fatalf("Failed to decode EUC-KR: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
