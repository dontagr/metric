package transport

import "bytes"

type (
	Transport interface {
		NewRequest(compressedBody *bytes.Buffer, HashSHA256 []string, w int) error
	}
)
