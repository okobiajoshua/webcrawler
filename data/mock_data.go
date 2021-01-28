package data

import (
	"github.com/stretchr/testify/mock"
)

// MockData struct
type MockData struct {
	mock.Mock
}

// NewMockData returns a fetch struct
func NewMockData() *MockData {
	return &MockData{}
}

// GetHTML returns html string
func (m *MockData) GetHTML(url []byte) ([]byte, error) {
	args := m.Called(url)
	content := args.Get(0)
	if content == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
