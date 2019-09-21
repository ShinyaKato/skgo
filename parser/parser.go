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

func (p *Parser) insertVariable(ident string) int {
  if _, ok := p.variables[ident]; !ok {
    p.stack += 4
    offset := -p.stack
    p.variables[ident] = offset
    return offset
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
  switch t := p.next(); t.Type {
  case token.INT_CONST:
    return &node.IntConstExpr {
      IntValue: t.IntValue,
    }

  case token.IDENT:
    if p.read("(") {
      args := []node.Expr {}
      if p.peek().Type != ")" {
        for {
          args = append(args, p.parseExpr())
          if !p.read(",") {
            break
          }
        }
      }
      p.expect(")")
      return &node.CallExpr {
        Callee: t.Ident,
        Args: args,
      }
    } else {
      return &node.IdentExpr {
        Offset: p.lookupVariable(t.Ident),
      }
    }

  case "(":
    expr := p.parseExpr()
    p.expect(")")
    return expr

  default:
    panic(fmt.Sprintf("invalid primary expression: %s.", t.Type))
  }
}

func (p *Parser) parseUnaryExpr() node.Expr {
  switch {
    case p.read("!"):
      var expr node.NotExpr
      expr.Expr = p.parseUnaryExpr()
      return &expr
    default:
      return p.parsePrimaryExpr()
  }
}

func (p *Parser) parseMulExpr() node.Expr {
  expr := p.parseUnaryExpr()

LOOP:
  for {
    switch {
    case p.read("*"):
      expr = &node.MulExpr {
        Lhs: expr,
        Rhs: p.parseUnaryExpr(),
      }
    case p.read("/"):
      expr = &node.DivExpr {
        Lhs: expr,
        Rhs: p.parseUnaryExpr(),
      }
    case p.read("%"):
      expr = &node.ModExpr {
        Lhs: expr,
        Rhs: p.parseUnaryExpr(),
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

func (p *Parser) parseComparisonExpr() node.Expr {
  expr := p.parseAddExpr()

LOOP:
  for {
    switch {
    case p.read("=="):
      expr = &node.EqualExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    case p.read("!="):
      expr = &node.NotEqualExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    case p.read("<"):
      expr = &node.LessExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    case p.read("<="):
      expr = &node.LessEqualExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    case p.read(">"):
      expr = &node.GreaterExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    case p.read(">="):
      expr = &node.GreaterEqualExpr {
        Lhs: expr,
        Rhs: p.parseComparisonExpr(),
      }
    default:
      break LOOP
    }
  }

  return expr
}

func (p *Parser) parseExpr() node.Expr {
  return p.parseComparisonExpr()
}

func (p *Parser) parseStmt() node.Stmt {
  switch {
  case p.read("var"):
    for {
      t := p.expect(token.IDENT)
      p.insertVariable(t.Ident)
      if !p.read(",") {
        break
      }
    }
    return nil

  case p.read("return"):
    var stmt node.ReturnStmt
    if p.peek().Type != ";" {
      stmt.ReturnExpr = p.parseExpr()
    }
    return &stmt

  case p.read("if"):
    var stmt node.IfStmt
    stmt.CondExpr = p.parseExpr()
    stmt.ThenBlock = p.parseBlock()
    if p.read("else") {
      stmt.ElseBlock = p.parseBlock()
    }
    return &stmt

  default:
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
}

func (p *Parser) parseBlock() *node.Block {
  list := []node.Stmt {}

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

func (p *Parser) parseFunctionDecl() *node.FunctionDecl {
  p.variables = map[string]int {}
  p.stack = 0

  p.expect("func")
  name := p.expect(token.IDENT).Ident
  p.expect("(")
  paramOffsets := []int {}
  if p.peek().Type != ")" {
    for {
      param := p.expect(token.IDENT).Ident
      paramOffsets = append(paramOffsets, p.insertVariable(param))
      if !p.read(",") {
        break
      }
    }
  }
  p.expect(")")
  body := p.parseBlock()

  if len(paramOffsets) > 6 {
    panic(fmt.Sprintf("too many parameters: %s", name))
  }

  return &node.FunctionDecl {
    Name: name,
    ParamOffsets: paramOffsets,
    Body: body,
    Stack: p.stack,
  }
}

func (p *Parser) Parse() *node.SourceFile {
  topLevelDecls := []node.TopLevelDecl {}
  for {
    if p.peek().Type == token.EOF {
      break
    }
    topLevelDecls = append(topLevelDecls, p.parseFunctionDecl())
  }

  return &node.SourceFile {
    TopLevelDecls: topLevelDecls,
  }
}

func New(tokens []*token.Token) *Parser {
  return &Parser {
    tokens: tokens,
    pos: 0,
  }
}
