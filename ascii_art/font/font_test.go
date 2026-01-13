package font

import (
	"fmt"
	"testing"
)

func TestASCII_ART_STRING(t *testing.T) {
	str := AsciiArt("Hello World", "big")
	fmt.Println(str)
}
