package cookie

import "testing"

func TestName(t *testing.T) {
	HttpsCookie = true
	t.Logf("%t", HttpsCookie)
	HttpsCookie = false
	t.Logf("%t", HttpsCookie)
}
