package lexer

import (
  "fmt"
  "unicode"
)

import (
  "../token"
)

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

func (l *Lexer) nextToken() *token.Token {
  switch {
  case !l.hasNext():
    return &token.Token { Type: token.EOF }

  case unicode.IsSpace(l.peek()):
    for l.hasNext() && unicode.IsSpace(l.peek()) {
      l.next()
    }
    return &token.Token { Type: token.SPACE }

  case unicode.IsDigit(l.peek()):
    intValue := 0
    for l.hasNext() && unicode.IsDigit(l.peek()) {
      intValue = intValue * 10 + int(l.next() - '0')
    }
    return &token.Token { Type: token.INT_CONST, IntValue: intValue }

  case l.read('{'):
    return &token.Token { Type: "{" }

  case l.read('}'):
    return &token.Token { Type: "}" }

  case l.read('('):
    return &token.Token { Type: "(" }

  case l.read(')'):
    return &token.Token { Type: ")" }

  case l.read('*'):
    return &token.Token { Type: "*" }

  case l.read('/'):
    return &token.Token { Type: "/" }

  case l.read('%'):
    return &token.Token { Type: "%" }

  case l.read('+'):
    return &token.Token { Type: "+" }

  case l.read('-'):
    return &token.Token { Type: "-" }

  case l.read(';'):
    return &token.Token { Type: ";" }

  default:
    panic(fmt.Sprintf("tokenize: unexpected character: %c.", l.peek()))
  }
}

func (l *Lexer) Tokenize() []*token.Token {
  tokens := []*token.Token {}

  for {
    t := l.nextToken()
    if t.Type == token.SPACE {
      continue
    }
    tokens = append(tokens, t)
    if t.Type == token.EOF {
      break
    }
  }

  return tokens
}

func New(src string) *Lexer {
  return &Lexer {
    src: src,
    chars: []rune(src),
    pos: 0,
  }
}
