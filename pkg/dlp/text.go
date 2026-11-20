package dlp

import (
	"bytes"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	textunicode "golang.org/x/text/encoding/unicode"
)

func hasTextExtension(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".csv", ".tsv", ".txt":
		return true
	default:
		return false
	}
}

func decodeText(name string, data []byte) (string, bool) {
	var err error
	var converted bool
	if bytes.HasPrefix(data, []byte{0xff, 0xfe}) || bytes.HasPrefix(data, []byte{0xfe, 0xff}) {
		data, err = textunicode.UTF16(textunicode.LittleEndian, textunicode.UseBOM).NewDecoder().Bytes(data)
		converted = true
	} else if !utf8.Valid(data) && hasTextExtension(name) {
		data, err = japanese.ShiftJIS.NewDecoder().Bytes(data)
		converted = true
	}
	if err != nil || !utf8.Valid(data) {
		return "", false
	}

	text := strings.TrimPrefix(string(data), "\ufeff")
	if strings.TrimSpace(text) == "" || converted && strings.ContainsRune(text, utf8.RuneError) {
		return "", false
	}
	if strings.IndexFunc(text, func(r rune) bool {
		return unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' && r != '\f'
	}) >= 0 {
		return "", false
	}
	return text, true
}
