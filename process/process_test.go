package process

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/monzo/webcrawler/data"
	"github.com/monzo/webcrawler/service"
	"github.com/monzo/webcrawler/store"
	"github.com/stretchr/testify/mock"
)

const (
	homePage = `<html>
<a href="/">Home</a>
<a href="/about">About</a>
<a href="/contact">Contact</a>
</html>`
)

func TestProcess(t *testing.T) {
	mockData := data.NewMockData()
	mockPub := service.NewMockProducer()
	mockCache := store.NewMockDataStore()

	mockData.On("GetHTML", []byte("http://xyz.com/")).Return([]byte(homePage), nil)
	mockPub.On("Publish", mock.Anything, mock.Anything).Return(nil)
	mockCache.On("Save", mock.Anything, mock.Anything).Return()
	mockCache.On("Fetch", mock.Anything).Return("", nil)

	process := NewProcess(mockData, mockPub, mockCache)
	err := process.Process([]byte("http://xyz.com/"))
	assert.Equal(t, err, nil)
}
