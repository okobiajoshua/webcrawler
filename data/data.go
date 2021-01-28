package data

// Data interface
type Data interface {
	GetHTML(url []byte) ([]byte, error)
}
