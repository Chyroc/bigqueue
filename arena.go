package bigqueue

import (
	"os"
	"syscall"

	"github.com/grandecola/mmap"
)

const (
	cFilePerm = 0744
)

// arena is an abstraction for a memory mapped file of a given size
type arena struct {
	mmap.IMmap
}

// newArena returns pointer to an arena. It takes a file location and mmaps it.
// If file location does not exist, it creates file of given size.
func newArena(file string, size int) (*arena, error) {
	fd, err := openOrCreateFile(file, size)
	if err != nil {
		return nil, err
	}

	m, err := mmap.NewSharedFileMmap(fd, 0, size, syscall.PROT_READ|syscall.PROT_WRITE)
	if err != nil {
		return nil, err
	}

	// We can close the file descriptor here
	if err := fd.Close(); err != nil {
		return nil, err
	}

	return &arena{
		IMmap: m,
	}, nil
}

// openOrCreateFile opens the file if it exists,
// otherwise creates a new file of given size
func openOrCreateFile(file string, size int) (*os.File, error) {
	if _, errExist := os.Stat(file); errExist == nil {
		// open file
		fd, err := os.OpenFile(file, os.O_RDWR, cFilePerm)
		if err != nil {
			return nil, err
		}

		return fd, nil
	} else if os.IsNotExist(errExist) {
		// create an empty file
		fd, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, cFilePerm)
		if err != nil {
			return nil, err
		}

		// truncate the file to required size
		if err := os.Truncate(file, int64(size)); err != nil {
			return nil, err
		}

		return fd, nil
	} else {
		return nil, errExist
	}
}
