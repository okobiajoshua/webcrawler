package service

import "context"

// Producer interface
type Producer interface {
	Publish(ctx context.Context, msg []byte) error
}
