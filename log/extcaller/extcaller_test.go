package extcaller_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lishimeng/app-starter/log"
)

func TestAutoModuleSkipsWrapperFrames(t *testing.T) {
	var buf bytes.Buffer
	log.Config().Output(&buf).Text().Apply()
	log.Info("from subpackage test")
	out := buf.String()
	if !strings.Contains(out, "module=log/extcaller") {
		t.Fatalf("expected module=log/extcaller, got: %q", out)
	}
}
