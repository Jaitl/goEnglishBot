package utils

import (
	"fmt"
	"math"
	"time"
)

func DurationPretty(d time.Duration) string {
	hours := int64(math.Mod(d.Hours(), 24))
	minutes := int64(math.Mod(d.Minutes(), 60))
	seconds := int64(math.Mod(d.Seconds(), 60))

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
