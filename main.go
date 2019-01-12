package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/funayman/hack-assembler/parser"
)

func main() {
	fmt.Println("The Hack Assembler")

	args := os.Args
	if len(args) < 2 {
		exitOnErr(fmt.Errorf("you must supply a source file"))
	}

	fn := args[1]
	f, err := os.Open(fn)
	exitOnErr(err)
	defer f.Close()

	p, _ := parser.New(f)

	index := strings.LastIndex(fn, filepath.Ext(fn))
	output := fn[:index] + ".hack"

	out, err := os.Create(output)
	exitOnErr(err)
	defer out.Close()

	for _, inst := range p.Parse() {
		fmt.Fprintln(out, inst)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
