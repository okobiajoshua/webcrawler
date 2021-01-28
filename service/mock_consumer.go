package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockConsumer struct
type MockConsumer struct {
	mock.Mock
}

// NewMockConsumer method
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{}
}

// Consume mock method
func (m *MockConsumer) Consume(ctx context.Context, f func([]byte) error) {
	m.Called(ctx, f)
}
