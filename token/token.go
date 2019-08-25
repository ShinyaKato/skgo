package token

type TokenType string

const (
  EOF = "EOF"
  SPACE = "SPACE"

  INT_CONST = "INT_CONST"
)

type Token struct {
  Type TokenType
  IntValue int
}
