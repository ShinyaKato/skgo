package node

import (
  "fmt"
)

var nextLabel = 0
var returnLabel int

func assignLabel() int {
  label := nextLabel
  nextLabel++
  return label
}

func (e *IntConstExpr) GenExpr() {
  fmt.Printf("  movl $%d, %%eax\n", e.IntValue)
  fmt.Printf("  pushq %%rax\n")
}

func (e *IdentExpr) GenExpr() {
  fmt.Printf("  movl %d(%%rbp), %%eax\n", e.Offset)
  fmt.Printf("  pushq %%rax\n")
}

func (e *CallExpr) GenExpr() {
  argRegs := []string { "rdi", "rsi", "rdx", "rcx", "r8", "r9" }

  fmt.Printf("  movq %%rsp, %%rax\n")
  fmt.Printf("  andq $-0x08, %%rsp\n")
  fmt.Printf("  pushq %%rax\n")
  for _, arg := range e.Args {
    arg.GenExpr()
  }
  for i := len(e.Args) - 1; i >= 0; i-- {
    fmt.Printf("  popq %%%s\n", argRegs[i])
  }
  fmt.Printf("  call %s\n", e.Callee)
  fmt.Printf("  popq %%rsp\n")
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

func (s *ReturnStmt) GenStmt() {
  s.ReturnExpr.GenExpr()
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  jmp .L%d\n", returnLabel)
}

func (b *Block) GenBlock() {
  for _, s := range b.StmtList {
    s.GenStmt()
  }
}

func (s *IfStmt) GenStmt() {
  elseLabel := assignLabel()
  endLabel := assignLabel()
  s.CondExpr.GenExpr()
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  cmpl $0, %%eax\n")
  fmt.Printf("  je .L%d\n", elseLabel)
  s.ThenBlock.GenBlock()
  fmt.Printf("  jmp .L%d\n", endLabel)
  fmt.Printf(".L%d:\n", elseLabel)
  if s.ElseBlock != nil {
    s.ElseBlock.GenBlock()
  }
  fmt.Printf(".L%d:\n", endLabel)
}

func (f *FunctionDecl) GenTopLevelDecl() {
  argRegs := []string { "edi", "esi", "edx", "ecx", "r8d", "r9d" }

  returnLabel = assignLabel()
  fmt.Printf("  .global %s\n", f.Name)
  fmt.Printf("%s:\n", f.Name)
  fmt.Printf("  pushq %%rbp\n")
  fmt.Printf("  movq %%rsp, %%rbp\n")
  fmt.Printf("  subq $%d, %%rsp\n", (f.Stack + 15) / 16 * 16)
  for i, offset := range f.ParamOffsets {
    fmt.Printf("  movl %%%s, %d(%%rbp)\n", argRegs[i], offset)
  }
  f.Body.GenBlock()
  fmt.Printf(".L%d:\n", returnLabel)
  fmt.Printf("  leave\n")
  fmt.Printf("  ret\n")
}

func (s *SourceFile) GenSourceFile() {
  for _, t := range s.TopLevelDecls {
    t.GenTopLevelDecl()
  }
}
