package node

type Expr interface {
  Generate()
}

type IntConstExpr struct {
  IntValue int
}

type MulExpr struct {
  Lhs, Rhs Expr
}

type DivExpr struct {
  Lhs, Rhs Expr
}

type ModExpr struct {
  Lhs, Rhs Expr
}

type AddExpr struct {
  Lhs, Rhs Expr
}

type SubExpr struct {
  Lhs, Rhs Expr
}
