package store

// DataStore interface
type DataStore interface {
	Save(urlVal string, value string) error
	Fetch(key string) (string, error)
}
