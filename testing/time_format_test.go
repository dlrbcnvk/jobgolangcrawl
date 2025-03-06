package testing

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	now := time.Now()
	str := now.Format("2006-01-02 15:04")

	expected := len([]rune(str))
	if len([]rune(str)) != expected {
		t.Errorf("format failed")
	}
	fmt.Println(str)
}
