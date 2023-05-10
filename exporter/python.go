package exporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/dc-dc-dc/cgo"
)

type PythonExporter struct {
	fmap map[string]*cgo.Func
}

func NewPythonExporter(fmap map[string]*cgo.Func) *PythonExporter {
	return &PythonExporter{
		fmap: fmap,
	}
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
				if s.Name == "printf" {
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
				if f, ok := e.fmap[s.Name]; ok {
					if len(f.Args) != len(s.Args) {
						fmt.Printf("function %s expects %d arguments, but %d were provided\n", s.Name, len(f.Args), len(s.Args))
						break
					}
					fmt.Fprintf(w, "\t%s(", s.Name)
					for _, arg := range s.Args {
						fmt.Fprintf(w, "%s, ", literalToPy(arg.Value))
					}
					fmt.Fprint(w, ")\n")
				}
			}
		}
	}
}
