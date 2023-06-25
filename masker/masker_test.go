package masker

import (
	"testing"
)

func TestPhoneMask001(t *testing.T) {
	var input = "13245671145"
	t.Log(Phone(input))
}
func TestPhoneMask002(t *testing.T) {
	var input = "1324"
	t.Log(Phone(input))
}
