package gdtcore

// Appendable allows some runnable thing to be added to it.
type Appendable interface {
	Append(Runnable)
}
