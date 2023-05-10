package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dc-dc-dc/cgo"
	"github.com/dc-dc-dc/cgo/exporter"
)

var (
	filePath string = "src/hello.c"
	output          = flag.Bool("o", false, "output to file")
	target          = flag.String("target", "python3", "target architecture")
)

func init() {
	flag.Parse()
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
	var w io.Writer = os.Stdout

	if *output {
		fmt.Print("output to file\n")
		f, err := os.Create("out/main.py")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w = f
	}
	fmap := map[string]*cgo.Func{}
	lexer := cgo.NewLexer(filePath, string(res))
	ffunc := cgo.ParseFunction(lexer)
	exporter := exporter.NewPythonExporter()
	for ffunc != nil {
		fmap[ffunc.Name] = ffunc
		ffunc = cgo.ParseFunction(lexer)
	}
	for _, ffunc := range fmap {
		exporter.Export(ffunc, w)
	}
}
