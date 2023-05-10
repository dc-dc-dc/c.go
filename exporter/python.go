package exporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/dc-dc-dc/cgo"
)

type PythonExporter struct {
}

func NewPythonExporter() *PythonExporter {
	return &PythonExporter{}
}

func literalToPy(v interface{}) string {
	if t, ok := v.(string); ok {
		return fmt.Sprintf("\"%s\"", t)
	}
	return fmt.Sprintf("%v", v)
}

func (e *PythonExporter) Export(ffunc *cgo.Func, w io.Writer) {
	isEntryPoint := ffunc.Name == "main" && ffunc.Type == cgo.TypeInt
	fmt.Printf("printing function name: %s, returns: %s, args: %v, entry: %t\n", ffunc.Name, ffunc.Type, ffunc.Args, isEntryPoint)
	if isEntryPoint {
		fmt.Fprintf(w, "if __name__ == \"__main__\":\n")
	} else {
		fmt.Fprintf(w, "def %s(", ffunc.Name)
		for _, arg := range ffunc.Args {
			fmt.Fprintf(w, "%s, ", arg.Name)
		}
		fmt.Fprint(w, "):\n")
	}
	for _, s := range ffunc.Body {
		switch s := s.(type) {
		case *cgo.ReturnStmt:
			if !isEntryPoint {
				fmt.Fprintf(w, "\treturn %s\n", literalToPy(s.Value))
			}
		case *cgo.FuncCallStmt:
			{
				frmt := s.Args[0].Value.(string)
				if len(s.Args) > 1 {
					subs := " % ("
					for _, arg := range s.Args[1:] {
						subs += literalToPy(arg.Value) + ", "
					}
					subs += ")"
					fmt.Fprintf(w, "\tprint(%s%s)\n", literalToPy(frmt), subs)
					break
				}
				if strings.HasSuffix(frmt, "\\n") {
					frmt = string(frmt[:len(frmt)-2])
				}
				fmt.Fprintf(w, "\tprint(%s)\n", literalToPy(frmt))
			}
		}
	}
}
