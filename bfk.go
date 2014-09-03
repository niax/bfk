package main

import (
	"flag"
	"github.com/niax/bfk/brainfuck"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()
	programs := flag.Args()

	for _, filename := range programs {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
		} else {
			m := brainfuck.NewMachine(string(content), os.Stdin, os.Stdout)
			m.Run()
		}
	}
}
