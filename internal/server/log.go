package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

// Appends a record the the log and returns its offset
// or an error if one occurs.
func (c *Log) Append(rec Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rec.Offset = uint64(len(c.records))
	c.records = append(c.records, rec)
	return rec.Offset, nil
}

func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if offset > uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
