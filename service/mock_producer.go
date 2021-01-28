package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockProducer struct
type MockProducer struct {
	mock.Mock
}

// NewMockProducer returns a MockProducer struct
func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

// Publish mock method
func (m *MockProducer) Publish(ctx context.Context, msg []byte) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
