package service

import "context"

// Consumer interface
type Consumer interface {
	Consume(ctx context.Context, f func([]byte) error)
}
