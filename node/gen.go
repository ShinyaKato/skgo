package node

import (
  "fmt"
)

func (e *IntConstExpr) GenExpr() {
  fmt.Printf("  movl $%d, %%eax\n", e.IntValue)
  fmt.Printf("  pushq %%rax\n")
}

func (e *IdentExpr) GenExpr() {
  fmt.Printf("  movl %d(%%rbp), %%eax\n", e.Offset)
  fmt.Printf("  pushq %%rax\n")
}

func (e *MulExpr) GenExpr() {
  e.Lhs.GenExpr()
  e.Rhs.GenExpr()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  imull %%ecx\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *DivExpr) GenExpr() {
  e.Lhs.GenExpr()
  e.Rhs.GenExpr()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  cltd\n")
  fmt.Printf("  idivl %%ecx\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *ModExpr) GenExpr() {
  e.Lhs.GenExpr()
  e.Rhs.GenExpr()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  cltd\n")
  fmt.Printf("  idivl %%ecx\n")
  fmt.Printf("  pushq %%rdx\n")
}

func (e *AddExpr) GenExpr() {
  e.Lhs.GenExpr()
  e.Rhs.GenExpr()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  addl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *SubExpr) GenExpr() {
  e.Lhs.GenExpr()
  e.Rhs.GenExpr()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  subl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}

func (s *ExprStmt) GenStmt() {
  s.Expr.GenExpr()
  fmt.Printf("  popq %%rax\n")
}

func (s *Assign) GenStmt() {
  s.Rhs.GenExpr()
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  movl %%eax, %d(%%rbp)\n", s.Lhs.Offset)
}

func (b *Block) GenBlock() {
  for _, s := range b.StmtList {
    s.GenStmt()
  }
}
