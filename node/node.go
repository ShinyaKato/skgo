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
  Args []Expr
}

type NotExpr struct {
  Expr Expr
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

type EqualExpr struct {
  Lhs, Rhs Expr
}

type NotEqualExpr struct {
  Lhs, Rhs Expr
}

type LessExpr struct {
  Lhs, Rhs Expr
}

type LessEqualExpr struct {
  Lhs, Rhs Expr
}

type GreaterExpr struct {
  Lhs, Rhs Expr
}

type GreaterEqualExpr struct {
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

type ReturnStmt struct {
  ReturnExpr Expr
}

type Block struct {
  StmtList []Stmt
}

type IfStmt struct {
  CondExpr Expr
  ThenBlock *Block
  ElseBlock *Block
}

type TopLevelDecl interface {
  GenTopLevelDecl()
}

type FunctionDecl struct {
  Name string
  ParamOffsets []int
  Body *Block
  Stack int
}

type SourceFile struct {
  TopLevelDecls []TopLevelDecl
}
