package trace

import (
	"fmt"
	"io"
)

// Tracer is an interface that describes an object capable of tracing
// events throughout the code

type Tracer interface{
	Trace(...interface{})
}

type tracer struct{
	out io.Writer
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}){}

func (t *tracer) Trace(a ...interface{}){
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// Creates a new instance of a Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Gets a Tracer that does nothing
func Off() Tracer{
	return &nilTracer{}
}