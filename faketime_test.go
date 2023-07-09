package faketime

import (
	"fmt"
	"testing"
	"time"
)

func TestFakeTime_FixTime(t *testing.T) {
	FixTime(2000, 12, 30, 14, 00, 00)
	fmt.Println(time.Now())

	FixTime(2001, 12, 30, 14, 00, 00)
	fmt.Println(time.Now())
}
