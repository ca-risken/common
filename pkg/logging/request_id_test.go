package logging

import (
	"strings"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantPrefix string
		loops      int
	}{
		{
			name:       "OK-100",
			input:      "prefix",
			wantPrefix: "prefix-",
			loops:      100,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotIDs := []string{}
			for i := 0; i < c.loops; i++ {
				got, err := GenerateRequestID(c.input)
				if err != nil {
					t.Fatalf("Unexpected GenerateRequestID error: err=%+v", err)
				}
				if !strings.HasPrefix(got, c.wantPrefix) {
					t.Fatalf("Unexpected not match ID format: want=%s, got=%s", c.wantPrefix, got)
				}
				// fmt.Println(got)
				gotIDs = append(gotIDs, got)
			}
			for _, target := range gotIDs {
				for _, id := range gotIDs {
					if target == id {
						t.Fatalf("Detect same id pattern on the %d loops, ID=%s", c.loops, target)
					}
				}
			}
		})
	}
}
