package core

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// DecodeEUCKR converts EUC-KR (CP949) encoded bytes to a UTF-8 string.
func DecodeEUCKR(b []byte) (string, error) {
	r := transform.NewReader(bytes.NewReader(b), korean.EUCKR.NewDecoder())
	decoded, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
