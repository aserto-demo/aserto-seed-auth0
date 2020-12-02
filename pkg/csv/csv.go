package csv

import (
	enccsv "encoding/csv"
	"os"
	"strings"
)

// Reader -.
type Reader struct {
	cols      map[string]int
	row       []string
	csvReader *enccsv.Reader
}

// NewCsvReader -.
func NewCsvReader() *Reader {
	return &Reader{}
}

// Open -.
func (r *Reader) Open(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	r.csvReader = enccsv.NewReader(f)

	// read first row for metadata
	r.cols = make(map[string]int)
	if r.row, err = r.csvReader.Read(); err != nil {
		return err
	}

	for index, name := range r.row {
		r.cols[name] = index
	}

	return nil
}

// Read -.
func (r *Reader) Read() (err error) {
	if r.row, err = r.csvReader.Read(); err != nil {
		return err
	}
	return nil
}

// Get -.
func (r *Reader) Get(col string) string {
	i, ok := r.cols[col]
	if ok {
		return r.row[i]
	}
	return ""
}

// GetToLower -.
func (r *Reader) GetToLower(col string) string {
	i, ok := r.cols[col]
	if ok {
		return strings.ToLower(r.row[i])
	}
	return ""
}
