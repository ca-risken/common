package dlp

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"path"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/encoding/japanese"
)

func TestScanFile(t *testing.T) {
	t.Parallel()

	const csvText = "氏名,メール\n髙橋,person@example.test\n"
	csv, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(csvText))
	if err != nil {
		t.Fatal(err)
	}
	var jpg, pngImage bytes.Buffer
	pixels := image.NewRGBA(image.Rect(0, 0, 2, 2))
	if err := jpeg.Encode(&jpg, pixels, nil); err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(&pngImage, pixels); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig("")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		data []byte
		want *Classification
	}{
		{name: "folder/updated.csv", data: csv, want: &Classification{Input: "text"}},
		{name: "table.tsv", data: []byte("name\temail"), want: &Classification{Input: "text"}},
		{name: "record.json", data: []byte(`{"name":"example"}`), want: &Classification{Input: "text"}},
		{name: "excluded.js", data: []byte("const x = 1")},
		{name: "image.jpg", data: jpg.Bytes(), want: &Classification{Input: "image"}},
		{name: "image.png", data: pngImage.Bytes(), want: &Classification{Input: "image"}},
		{name: "unsupported archive-based file.docx", data: []byte{'P', 'K', 3, 4, 0, 0}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var calls int
			classifier := classifyPDFPageFunc(func(_ context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
				calls++
				if req.FileName != path.Base(test.name) {
					t.Errorf("file name = %q, want %q", req.FileName, path.Base(test.name))
				}
				if test.name == "folder/updated.csv" && req.Text != csvText {
					t.Errorf("expected decoded UTF-8 CSV, got %q", req.Text)
				}
				if test.want != nil && test.want.Input == "image" && !bytes.Equal(req.Image, test.data) {
					t.Error("expected original image bytes")
				}
				return &ai.ClassifySensitiveContentResponse{}, nil
			})
			var want *ScanResponse
			if test.want != nil {
				want = &ScanResponse{Complete: true, Classification: test.want}
			}
			got, err := NewScanner(classifier, config).Scan(t.Context(), ScanRequest{FileName: test.name, Reader: bytes.NewReader(test.data)})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("scan response mismatch (-want +got):\n%s", diff)
			}
			if calls != 0 && want == nil || calls != 1 && want != nil {
				t.Errorf("Core calls = %d for result %+v", calls, got)
			}
		})
	}
}

func TestDecodeText(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		data []byte
		want string
		ok   bool
	}{
		{name: "BOM-prefixed UTF-8", data: []byte("\ufeff氏名,住所"), want: "氏名,住所", ok: true},
		{name: "BOM-prefixed UTF-16", data: []byte{0xff, 0xfe, 'a', 0, ',', 0, 'b', 0}, want: "a,b", ok: true},
		{name: "binary mislabeled as text", data: []byte("A\x00B")},
		{name: "blank input", data: []byte(strings.Repeat(" ", 2))},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			text, ok := decodeText("data.csv", test.data)
			if text != test.want || ok != test.ok {
				t.Errorf("decodeText() = %q, %t; want %q, %t", text, ok, test.want, test.ok)
			}
		})
	}
}
