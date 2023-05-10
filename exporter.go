package cgo

import "io"

type Exporter interface {
	Export(ffunc *Func, w io.Writer) error
}
