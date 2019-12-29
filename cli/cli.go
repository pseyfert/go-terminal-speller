package main

import (
	"flag"
	termspeller "github.com/pseyfert/go-terminal-speller/pkg"
	"golang.org/x/text/unicode/norm"
	"os"
)

func main() {
	flag.Parse()

	t := termspeller.NewTranslator(os.Stdout)
	writer := norm.NFD.Writer(t)
	for _, fl := range flag.Args() {
		writer.Write([]byte(fl))
		writer.Write([]byte(" "))
	}
	writer.Write([]byte("\n"))
}
