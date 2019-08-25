package main

import "fmt"

func main() {
  fmt.Println("  .global main")
  fmt.Println("main:")
  fmt.Println("  movl $42, %eax")
  fmt.Println("  ret")
}
