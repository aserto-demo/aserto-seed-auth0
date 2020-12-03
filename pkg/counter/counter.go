package counter

import (
	"fmt"
	"os"
	"sync/atomic"
)

// Counter - accumulator for row, skipped and error counts.
type Counter struct {
	rowCounter  int32
	skipCounter int32
	errCounter  int32
}

// IncrRows - increment row counter.
func (c *Counter) IncrRows() {
	atomic.AddInt32(&c.rowCounter, 1)
}

// IncrSkipped - increment skipped row counter.
func (c *Counter) IncrSkipped() {
	atomic.AddInt32(&c.skipCounter, 1)
}

// IncrError - increment error counter.
func (c *Counter) IncrError() {
	atomic.AddInt32(&c.errCounter, 1)
}

// Print - print counter at interval % m.
func (c *Counter) Print(m int32) {
	linefeed := ""

	if m == 0 {
		linefeed = "\n"
		m = 1 // avoid divide by zero
	}

	if d := c.rowCounter % m; d == 0 {
		fmt.Fprintf(os.Stdout, "\033[2K\rrow count: %d skip count %d error count: %d%s",
			c.rowCounter,
			c.skipCounter,
			c.errCounter,
			linefeed,
		)
	}
}
