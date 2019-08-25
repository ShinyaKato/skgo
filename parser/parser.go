package parser

import (
  "fmt"
)

import (
  "../token"
  "../node"
)

type Parser struct {
  tokens []*token.Token
  pos int
}

func (p *Parser) peek() *token.Token {
  return p.tokens[p.pos]
}

func (p *Parser) next() *token.Token {
  token := p.tokens[p.pos]
  p.pos++
  return token
}

func (p *Parser) read(kind string) bool {
  if p.tokens[p.pos].Kind == kind {
    p.pos++
    return true
  }

  return false
}

func (p *Parser) parsePrimaryExpr() node.Expr {
  token := p.next()

  switch token.Kind {
  case "IntConst":
    return &node.IntConstExpr {
      IntValue: token.IntValue,
    }

  default:
    panic(fmt.Sprintf("invalid primary expression: %s.", token.Kind))
  }
}

func (p *Parser) parseAddExpr() node.Expr {
  expr := p.parsePrimaryExpr()

  for {
    switch {
    case p.read("+"):
      expr = &node.AddExpr {
        Lhs: expr,
        Rhs: p.parsePrimaryExpr(),
      }
    case p.read("-"):
      expr = &node.SubExpr {
        Lhs: expr,
        Rhs: p.parsePrimaryExpr(),
      }
    default:
      return expr
    }
  }
}

func (p *Parser) parseExpr() node.Expr {
  return p.parseAddExpr()
}

func (p *Parser) Parse() node.Expr {
  return p.parseExpr()
}

func New(tokens []*token.Token) *Parser {
  return &Parser {
    tokens: tokens,
    pos: 0,
  }
}
