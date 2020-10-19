package systool

import (
	"fmt"
	"time"
)

// TimeTrack return a time used to complete a function and return it in a string
func TimeTrack(start time.Time, name string) string {
	elapsed := time.Since(start)
	return fmt.Sprintf("%s took %s", name, elapsed)
}
