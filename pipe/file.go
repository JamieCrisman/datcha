package pipe

import (
	"fmt"
	"io"
	"os"
)

func FileReader(filePath string) (io.ReadCloser, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	return f, nil
}
