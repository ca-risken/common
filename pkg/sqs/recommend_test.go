package sqs

import (
	"reflect"
	"testing"
)

func TestGetRecommend(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  recommend
	}{
		{
			name:  "Exists recommend",
			input: DataSourceGoogleSCC,
			want: recommend{
				Risk: "Failed to scan google:scc, So you are not gathering the latest security threat information.",
				Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/google/overview_gcp/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
			},
		},
		{
			name:  "Unknown recommend",
			input: "typeUnknown",
			want: recommend{
				Risk:           "",
				Recommendation: "",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getRecommend(c.input)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}
