// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package debug

import (
	"context"
	"fmt"
	"strings"

	gdtcontext "github.com/jaypipes/gdt-core/context"
)

// Printf writes a message with optional message arguments to the context's
// Debug output.
func Printf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	dbg := gdtcontext.Debug(ctx)
	dbg.Write([]byte(msg))
}

// Println writes a message with optional message arguments to the context's
// Debug output, ensuring there is a newline in the message line.
func Println(ctx context.Context, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	msg := fmt.Sprintf(format, args...)
	dbg := gdtcontext.Debug(ctx)
	dbg.Write([]byte(msg))
}
