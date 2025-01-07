package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	idx, err := newIndex(f, c)
	require.NoError(t, err)

	_, _, err = idx.Read(-1)
	require.ErrorIs(t, err, io.EOF)

	require.Equal(t, f.Name(), idx.Name())

	entries := []struct {
		Offset   uint32
		Position uint64
	}{
		{Offset: 0, Position: 0},
		{Offset: 1, Position: 10},
	}

	for _, want := range entries {
		err := idx.Write(want.Offset, want.Position)
		require.NoError(t, err)

		_, pos, err := idx.Read(int64(want.Offset))
		require.NoError(t, err)
		require.Equal(t, want.Position, pos)
	}

	_, _, err = idx.Read(int64(len(entries)))
	require.Error(t, err)

	idx.Close()

	// can rebuild
	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	idx, err = newIndex(f, c)
	require.NoError(t, err)

	offset, pos, err := idx.Read(-1)
	require.NoError(t, err)
	require.Equal(t, offset, entries[1].Offset)
	require.Equal(t, entries[1].Position, pos)
}
