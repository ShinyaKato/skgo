package main

import (
  "fmt"
  "os"
)

import (
  "./lexer"
  "./parser"
)

func getSourceFromArgs() string {
  cmd := os.Args[0]
  args := os.Args[1:]
  if len(args) != 1 {
    fmt.Fprintf(os.Stderr, "usage: %s [input srouce code]\n", cmd)
    os.Exit(1)
  }

  return args[0]
}

func main() {
  src := getSourceFromArgs()
  tokens := lexer.New(src).Tokenize()
  sourceFile := parser.New(tokens).Parse()

  sourceFile.GenSourceFile()
}
