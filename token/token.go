package token

type TokenType string

const (
  EOF = "EOF"
  SPACE = "SPACE"

  INT_CONST = "INT_CONST"
  IDENT = "IDENT"
)

type Token struct {
  Type TokenType
  IntValue int
  Ident string
}
