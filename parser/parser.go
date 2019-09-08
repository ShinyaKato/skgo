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

  variables map[string]int
  stack int
}

func (p *Parser) peek() *token.Token {
  return p.tokens[p.pos]
}

func (p *Parser) next() *token.Token {
  token := p.tokens[p.pos]
  p.pos++
  return token
}

func (p *Parser) read(tokenType token.TokenType) bool {
  if p.tokens[p.pos].Type == tokenType {
    p.pos++
    return true
  }

  return false
}

func (p *Parser) expect(tokenType token.TokenType) *token.Token {
  t := p.tokens[p.pos]
  if t.Type == tokenType {
    p.pos++
    return t
  }

  panic(fmt.Sprintf("%s is expected, but got %s.", tokenType, t.Type))
}

func (p *Parser) insertVariable(ident string) {
  if _, ok := p.variables[ident]; !ok {
    p.stack += 4
    p.variables[ident] = -p.stack
  } else {
    panic(fmt.Sprintf("duplicated variable declaration: %s.", ident))
  }
}

func (p *Parser) lookupVariable(ident string) int {
  if offset, ok := p.variables[ident]; ok {
    return offset
  } else {
    panic(fmt.Sprintf("undefined variable: %s.", ident))
  }
}

func (p *Parser) parsePrimaryExpr() node.Expr {
  t := p.next()

  switch t.Type {
  case token.INT_CONST:
    return &node.IntConstExpr {
      IntValue: t.IntValue,
    }

  case token.IDENT:
    return &node.IdentExpr {
      Offset: p.lookupVariable(t.Ident),
    }

  case "(":
    expr := p.parseExpr()
    p.expect(")")
    return expr

  default:
    panic(fmt.Sprintf("invalid primary expression: %s.", t.Type))
  }
}

func (p *Parser) parseMulExpr() node.Expr {
  expr := p.parsePrimaryExpr()

LOOP:
  for {
    switch {
    case p.read("*"):
      expr = &node.MulExpr {
        Lhs: expr,
        Rhs: p.parsePrimaryExpr(),
      }
    case p.read("/"):
      expr = &node.DivExpr {
        Lhs: expr,
        Rhs: p.parsePrimaryExpr(),
      }
    case p.read("%"):
      expr = &node.ModExpr {
        Lhs: expr,
        Rhs: p.parsePrimaryExpr(),
      }
    default:
      break LOOP
    }
  }

  return expr
}

func (p *Parser) parseAddExpr() node.Expr {
  expr := p.parseMulExpr()

LOOP:
  for {
    switch {
    case p.read("+"):
      expr = &node.AddExpr {
        Lhs: expr,
        Rhs: p.parseMulExpr(),
      }
    case p.read("-"):
      expr = &node.SubExpr {
        Lhs: expr,
        Rhs: p.parseMulExpr(),
      }
    default:
      break LOOP
    }
  }

  return expr
}

func (p *Parser) parseExpr() node.Expr {
  return p.parseAddExpr()
}

func (p *Parser) parseStmt() node.Stmt {
  if p.read("var") {
    for {
      t := p.expect(token.IDENT)
      p.insertVariable(t.Ident)
      if !p.read(",") {
        break
      }
    }
    return nil
  }

  expr := p.parseExpr()

  if p.read("=") {
    if ident, ok := expr.(*node.IdentExpr); ok {
      return &node.Assign {
        Lhs: ident,
        Rhs: p.parseExpr(),
      }
    } else {
      panic("currently, only identifier is supported for left hand side of assignment.")
    }
  }

  return &node.ExprStmt { Expr: expr }
}

func (p *Parser) parseBlock() *node.Block {
  var list = []node.Stmt {}

  p.expect("{")
  for p.peek().Type != "}" {
    stmt := p.parseStmt()
    if stmt != nil {
      list = append(list, stmt)
    }
    p.expect(";")
  }
  p.expect("}")

  return &node.Block {
    StmtList: list,
  }
}

func (p *Parser) Parse() (*node.Block, int) {
  block := p.parseBlock()

  if p.peek().Type != token.EOF {
    panic("invalid program.")
  }

  return block, p.stack
}

func New(tokens []*token.Token) *Parser {
  return &Parser {
    tokens: tokens,
    pos: 0,

    variables: map[string]int {},
    stack: 0,
  }
}
