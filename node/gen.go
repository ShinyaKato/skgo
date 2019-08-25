package node

import (
  "fmt"
)

func (e *IntConstExpr) Generate() {
  fmt.Printf("  movl $%d, %%eax\n", e.IntValue)
  fmt.Printf("  pushq %%rax\n")
}

func (e *MulExpr) Generate() {
  e.Lhs.Generate()
  e.Rhs.Generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  imull %%ecx\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *DivExpr) Generate() {
  e.Lhs.Generate()
  e.Rhs.Generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  cltd\n")
  fmt.Printf("  idivl %%ecx\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *ModExpr) Generate() {
  e.Lhs.Generate()
  e.Rhs.Generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  cltd\n")
  fmt.Printf("  idivl %%ecx\n")
  fmt.Printf("  pushq %%rdx\n")
}

func (e *AddExpr) Generate() {
  e.Lhs.Generate()
  e.Rhs.Generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  addl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}

func (e *SubExpr) Generate() {
  e.Lhs.Generate()
  e.Rhs.Generate()

  fmt.Printf("  popq %%rcx\n")
  fmt.Printf("  popq %%rax\n")
  fmt.Printf("  subl %%ecx, %%eax\n")
  fmt.Printf("  pushq %%rax\n")
}
