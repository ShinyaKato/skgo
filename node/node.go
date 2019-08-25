package node

type Expr interface {
  Generate()
}

type IntConstExpr struct {
  IntValue int
}

type AddExpr struct {
  Lhs, Rhs Expr
}

type SubExpr struct {
  Lhs, Rhs Expr
}
