package socket

import (
	"io"
	"log"
	"os"

	"go.vxn.dev/dish/pkg/config"
)

// fetchSocketsFromFile opens a file and returns [io.ReadCloser] for reading from the stream.
func fetchSocketsFromFile(cfg *config.Config) (io.ReadCloser, error) {
	file, err := os.Open(cfg.Source)
	if err != nil {
		return nil, err
	}

	// TODO: replace with logger
	if cfg.Verbose {
		log.Printf("Fetching sockets from the source (%s)", cfg.Source)
	}

	return file, nil
}
