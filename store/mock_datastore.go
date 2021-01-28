package store

import (
	"github.com/stretchr/testify/mock"
)

// MockDataStore struct
type MockDataStore struct {
	mock.Mock
}

// NewMockDataStore return a MockDataStore struct
func NewMockDataStore() *MockDataStore {
	return &MockDataStore{}
}

// Save mock method
func (m *MockDataStore) Save(urlVal string, value string) error {
	args := m.Called(urlVal, value)
	return args.Error(0)
}

// Fetch mock method
func (m *MockDataStore) Fetch(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}
