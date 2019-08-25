package main

import (
  "fmt"
  "os"
  "unicode"
)

type Token struct {
  int_value int
}

func tokenize(src string) []*Token {
  chars := []rune(src)
  pos := 0

  tokens := []*Token {}
  for {
    if pos == len(src) {
      break
    }

    switch {
    case unicode.IsSpace(chars[pos]):
      pos++

    case unicode.IsDigit(chars[pos]):
      int_value := 0
      for pos < len(src) && unicode.IsDigit(chars[pos]) {
        int_value = int_value * 10 + int(chars[pos] - '0')
        pos++
      }

      tokens = append(tokens, &Token {
        int_value: int_value,
      })

    default:
      panic(fmt.Sprintf("tokenize: unexpected character: %c.", chars[pos]))
    }
  }

  return tokens
}

func main() {
  cmd := os.Args[0]
  args := os.Args[1:]
  if len(args) != 1 {
    fmt.Fprintf(os.Stderr, "usage: %s [input srouce code]\n", cmd)
    os.Exit(1)
  }

  src := args[0]
  tokens := tokenize(src)

  if len(tokens) != 1 {
    panic("only one integer constant is expected.")
  }

  fmt.Printf("  .global main\n")
  fmt.Printf("main:\n")
  fmt.Printf("  movl $%d, %%eax\n", tokens[0].int_value)
  fmt.Printf("  ret\n")
}
