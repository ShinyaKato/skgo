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
  block, stack := parser.New(tokens).Parse()

  fmt.Printf("  .global main\n")
  fmt.Printf("main:\n")
  fmt.Printf("  pushq %%rbp\n")
  fmt.Printf("  movq %%rsp, %%rbp\n")
  fmt.Printf("  subq $%d, %%rsp\n", stack)
  block.GenBlock()
  fmt.Printf("  leave\n")
  fmt.Printf("  ret\n")
}
