package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvironMap(t *testing.T) {
	environ := []string{
		"FOO=bar",
		"BAZ=",
	}
	want := map[string]string{
		"FOO": "bar",
		"BAZ": "",
	}
	got := environMap(environ)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}
