package data

import (
	"io/ioutil"
	"net/http"
)

// Fetch struct
type Fetch struct {
}

// NewFetch returns a fetch struct
func NewFetch() *Fetch {
	return &Fetch{}
}

// GetHTML returns html string
func (f *Fetch) GetHTML(url []byte) ([]byte, error) {
	res, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return b, err
}
