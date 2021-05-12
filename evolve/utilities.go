package evolve

import (
	"fmt"
)

// Debug just prints its input arguments using `fmt.Println`.
// It's useful for `grep`ing it and deleting all its instances.
func Debug(args ...interface{}) {
	fmt.Println(args...)
}
