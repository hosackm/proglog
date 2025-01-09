package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offsetWidth   uint64 = 4 // 4 bytes for offset
	positionWidth uint64 = 8 // 8 bytes for position
	entryWidth           = offsetWidth + positionWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

// Returns a new index initialized with the data in the provided file.
// If the file is larger than the max allowed size in bytes from the
// provided config, it is truncated.
//
// The file is then mmap-ed for quicker access.
func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())

	if err = os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}

	idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

// Close the index after syncing to mmap and disk and truncating if
// the current size is larger than the config max size.
func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}

	if err := i.file.Sync(); err != nil {
		return err
	}

	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}

	return i.file.Close()
}

// Read the offset and position of the provided index-point within the index.
func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	if in == -1 {
		out = uint32((i.size / entryWidth) - 1)
	} else {
		out = uint32(in)
	}

	pos = uint64(out) * entryWidth
	if i.size < pos+entryWidth {
		return 0, 0, io.EOF
	}

	// Entry Frame
	//  pos         pos+offsetWidth      pos+entryWidth
	//  ↓                ↓                    ↓
	//  |--------4-------|----------8---------|
	//  <--offsetWidth-->
	//  <---------------entryWidth----------->
	out = enc.Uint32(i.mmap[pos : pos+offsetWidth])
	pos = enc.Uint64(i.mmap[pos+offsetWidth : pos+entryWidth])
	return out, pos, nil
}

// Write the offset and position to the end of index.
// If there is no room in the index EOF is returned.
func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+entryWidth {
		return io.EOF
	}

	enc.PutUint32(i.mmap[i.size:i.size+offsetWidth], off)
	enc.PutUint64(i.mmap[i.size+offsetWidth:i.size+entryWidth], pos)
	i.size += uint64(entryWidth)
	return nil
}

// Returns the name of the file used by this index.
func (i *index) Name() string {
	return i.file.Name()
}
