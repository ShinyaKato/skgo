package node

type Expr interface {
  GenExpr()
}

type IntConstExpr struct {
  IntValue int
}

type IdentExpr struct {
  Offset int
}

type CallExpr struct {
  Callee string
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

type Stmt interface {
  GenStmt()
}

type ExprStmt struct {
  Expr Expr
}

type Assign struct {
  Lhs *IdentExpr
  Rhs Expr
}

type Block struct {
  StmtList []Stmt
}

type TopLevelDecl interface {
  GenTopLevelDecl()
}

type FunctionDecl struct {
  Name string
  Body *Block
  Stack int
}

type SourceFile struct {
  TopLevelDecls []TopLevelDecl
}
