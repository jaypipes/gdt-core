// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

// Appendable allows some runnable thing to be added to it.
type Appendable interface {
	Append(Runnable)
}
