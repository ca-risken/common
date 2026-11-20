package dlp

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *Scanner) scanGZIPLog(ctx context.Context, name string, compressed []byte, maxBytes int64) (*ScanResponse, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("open gzip log: %w", err)
	}

	data, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err := errors.Join(err, reader.Close()); err != nil {
		return nil, fmt.Errorf("read gzip log: %w", err)
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("expanded gzip log exceeds the %d-byte limit", maxBytes)
	}

	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "text/") {
		return nil, nil
	}
	return s.scanContents(ctx, name, contentType, data)
}
