package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dc-dc-dc/cgo"
)

var (
	filePath string = "src/hello.c"
	output          = flag.Bool("o", false, "output to file")
	target          = flag.String("target", "python3", "target architecture")
)

func init() {
	flag.Parse()
}

func literalToPy(v interface{}) string {
	if t, ok := v.(string); ok {
		return fmt.Sprintf("\"%s\"", t)
	}
	return fmt.Sprintf("%v", v)
}

func genOutput(ffunc *cgo.Func, w io.Writer) {
	isEntryPoint := ffunc.Name.Value == "main" && ffunc.Type == cgo.TypeInt
	fmt.Printf("printing function name: %s, returns: %s, args: %v, entry: %t\n", ffunc.Name.Value, ffunc.Type, ffunc.Args, isEntryPoint)

	for _, s := range ffunc.Body {
		switch s := s.(type) {
		case *cgo.ReturnStmt:
			// TODO: Add return
			break
		case *cgo.FuncCallStmt:
			{
				frmt := s.Args[0].Value.(string)
				if len(s.Args) > 1 {
					subs := " % ("
					for _, arg := range s.Args[1:] {
						subs += literalToPy(arg.Value) + ", "
					}
					subs += ")"
					fmt.Fprintf(w, "print(%s%s)\n", literalToPy(frmt), subs)
					break
				}
				if strings.HasSuffix(frmt, "\\n") {
					frmt = string(frmt[:len(frmt)-2])
				}
				fmt.Fprintf(w, "print(%s)\n", literalToPy(frmt))
			}
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}

	if filePath == "" {
		fmt.Println("error: file path is empty")
		os.Exit(1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: failed to open the file \"%s\", err: %s", filePath, err.Error())
	}
	defer file.Close()

	res, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error: failed to read file, err: %s", err.Error())
		os.Exit(1)
	}
	// fmt.Print(string(res))
	fmt.Printf("%+v\n", *output)
	var w io.Writer = os.Stdout

	if *output {
		fmt.Print("output to file\n")
		f, err := os.Create("out.py")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w = f
	}

	lexer := cgo.NewLexer(filePath, string(res))
	ffunc := cgo.ParseFunction(lexer)
	for ffunc != nil {
		genOutput(ffunc, w)
		ffunc = cgo.ParseFunction(lexer)
	}
}
