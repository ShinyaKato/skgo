package main

import (
  "fmt"
  "os"
  "unicode"
)

type Token struct {
  kind string
  intValue int
}

type Expr interface {
  generate()
}

type IntConstExpr struct {
  intValue int
}

type AddExpr struct {
  lhs, rhs Expr
}

type SubExpr struct {
  lhs, rhs Expr
}

type Lexer struct {
  src string
  chars []rune
  pos int
}

func (l *Lexer) hasNext() bool {
  return l.pos < len(l.chars)
}

func (l *Lexer) peek() rune {
  return l.chars[l.pos]
}

func (l *Lexer) next() rune {
  c := l.chars[l.pos]
  l.pos++
  return c
}

func (l *Lexer) read(c rune) bool {
  if l.chars[l.pos] == c {
    l.pos++
    return true
  }

  return false
}

func (l *Lexer) nextToken() *Token {
  switch {
  case !l.hasNext():
    return &Token { kind: "EndOfFile" }

  case unicode.IsSpace(l.peek()):
    for l.hasNext() && unicode.IsSpace(l.peek()) {
      l.next()
    }
    return &Token { kind: "Space" }

  case unicode.IsDigit(l.peek()):
    intValue := 0
    for l.hasNext() && unicode.IsDigit(l.peek()) {
      intValue = intValue * 10 + int(l.next() - '0')
    }
    return &Token { kind: "IntConst", intValue: intValue }

  case l.read('+'):
    return &Token { kind: "+" }

  case l.read('-'):
    return &Token { kind: "-" }

  default:
    panic(fmt.Sprintf("tokenize: unexpected character: %c.", l.peek()))
  }
}

func tokenize(src string) []*Token {
  l := Lexer {
    src: src,
    chars: []rune(src),
    pos: 0,
  }

  tokens := []*Token {}
  for {
    token := l.nextToken()
    if token.kind == "Space" {
      continue
    }

    tokens = append(tokens, token)
    if token.kind == "EndOfFile" {
      break
    }
  }

  return tokens
}

type Parser struct {
  tokens []*Token
  pos int
}

func (p *Parser) peek() *Token {
  return p.tokens[p.pos]
}

func (p *Parser) next() *Token {
  token := p.tokens[p.pos]
  p.pos++
  return token
}

func (p *Parser) read(kind string) bool {
  if p.tokens[p.pos].kind == kind {
    p.pos++
    return true
  }

  return false
}

func (p *Parser) parsePrimaryExpr() Expr {
  token := p.next()

  switch token.kind {
  case "IntConst":
    return &IntConstExpr {
      intValue: token.intValue,
    }

  default:
    panic(fmt.Sprintf("invalid primary expression: %s.", token.kind))
  }
}

func (p *Parser) parseAddExpr() Expr {
  expr := p.parsePrimaryExpr()

  for {
    switch {
    case p.read("+"):
      expr = &AddExpr {
        lhs: expr,
        rhs: p.parsePrimaryExpr(),
      }
    case p.read("-"):
      expr = &SubExpr {
        lhs: expr,
        rhs: p.parsePrimaryExpr(),
      }
    default:
      return expr
    }
  }
}

func (p *Parser) parseExpr() Expr {
  return p.parseAddExpr()
}

func parse(tokens []*Token) Expr {
  p := Parser {
    tokens: tokens,
    pos: 0,
  }

  return p.parseExpr()
}

func (expr *IntConstExpr) generate() {
  fmt.Printf("  movl $%d, %%eax\n", expr.intValue)
  fmt.Printf("  pushq %%rax\n")
}

func (expr *AddExpr) generate() {
  expr.lhs.generate()
  expr.rhs.generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  addl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}

func (expr *SubExpr) generate() {
  expr.lhs.generate()
  expr.rhs.generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  subl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}

func generate(expr Expr) {
  fmt.Printf("  .global main\n")
  fmt.Printf("main:\n")
  fmt.Printf("  pushq %%rbp\n")
  fmt.Printf("  movq %%rsp, %%rbp\n")

  expr.generate()
  fmt.Printf("  popq %%rax\n")

  fmt.Printf("  leave\n")
  fmt.Printf("  ret\n")
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
  expr := parse(tokens)
  generate(expr)
}
