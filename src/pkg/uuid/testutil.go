package uuid

import "testing"

func SetDummy(t *testing.T, id string) (clear func()) {
	t.Helper()
	Generate = func() string {
		return id
	}
	return func() {
		Generate = generate
	}
}
