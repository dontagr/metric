package transport

type (
	Transport interface {
		NewRequest(income any, HashSHA256 []string, w int) error
	}
)
