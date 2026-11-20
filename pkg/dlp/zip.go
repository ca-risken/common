package dlp

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"slices"
	"strings"
)

const (
	maxZIPEntries       = 100
	maxZIPExpandedBytes = 25 << 20
	zipEncryptedFlag    = 1
)

func (s *Scanner) scanZIP(ctx context.Context, data []byte) (*ScanResponse, error) {
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open ZIP: %w", err)
	}
	if slices.ContainsFunc(archive.File, func(entry *zip.File) bool {
		return entry.Flags&zipEncryptedFlag != 0
	}) {
		return nil, nil
	}

	withinEntryLimit := len(archive.File) <= maxZIPEntries
	result := &ScanResponse{Complete: withinEntryLimit, ZIP: &ZIPResult{Truncated: !withinEntryLimit}}
	maxFileBytes := s.config.GetMaxSingleFileSizeBytes()
	remaining := int64(maxZIPExpandedBytes)
	for _, entry := range archive.File[:min(len(archive.File), maxZIPEntries)] {
		if err := ctx.Err(); err != nil {
			result.Complete = false
			return result, err
		}
		if entry.FileInfo().IsDir() || entry.UncompressedSize64 == 0 || !s.CanScanFile(entry.Name, 0) {
			continue
		}
		if !entry.Mode().IsRegular() || !fs.ValidPath(entry.Name) || strings.Contains(entry.Name, "\\") ||
			strings.EqualFold(path.Ext(entry.Name), ".zip") {
			result.Complete = false
			continue
		}
		if entry.UncompressedSize64 > uint64(min(maxFileBytes, remaining)) {
			result.Complete = false
			result.ZIP.Truncated = true
			continue
		}

		payload, err := readZIPEntry(entry)
		if err != nil {
			result.Complete = false
			return result, fmt.Errorf("read ZIP entry: %w", err)
		}
		remaining -= int64(len(payload))
		contentType := http.DetectContentType(payload)
		if contentType == "application/zip" {
			result.Complete = false
			continue
		}
		scan, err := s.scanContents(ctx, entry.Name, contentType, payload)
		result.Complete = result.Complete && err == nil && scan != nil && scan.Complete
		if scan != nil {
			result.ZIP.Entries = append(result.ZIP.Entries, ZIPEntry{Name: entry.Name, Scan: *scan})
		}
	}
	return result, nil
}

func readZIPEntry(entry *zip.File) ([]byte, error) {
	reader, err := entry.Open()
	if err != nil {
		return nil, err
	}
	// The declared size was checked against both limits before opening the entry.
	// Allow one extra byte to detect overruns and let the reader verify its checksum.
	size := int64(entry.UncompressedSize64)
	payload, err := io.ReadAll(io.LimitReader(reader, size+1))
	if err := errors.Join(err, reader.Close()); err != nil {
		return nil, err
	}
	if int64(len(payload)) > size {
		return nil, fmt.Errorf("entry exceeds its declared size")
	}
	return payload, nil
}
