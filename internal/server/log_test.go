package server

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	log := NewLog()

	// can't read from empty log
	_, err := log.Read(1)
	if err != ErrOffsetNotFound {
		t.FailNow()
	}

	for i := 0; i < 20; i++ {
		val := fmt.Sprintf("hello world message: %d", i)
		offset, err := log.Append(Record{Value: []byte(val)})
		if err != nil {
			t.Errorf("Append %d returned error: %s", i, err.Error())
		}
		if offset != uint64(i) {
			t.Errorf("expected offset of first append to be %d", i)
		}
	}
}
