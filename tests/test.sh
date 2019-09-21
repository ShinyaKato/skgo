mkdir -p tests/tmp

example() {
  expr=$1
  expected=$2

  ./skgo "$expr" > tests/tmp/out.s
  if [ $? -ne 0 ]; then
    echo "NG: failed to compile: $expr."
    exit 1
  fi

  gcc -o tests/tmp/out tests/tmp/out.s
  if [ $? -ne 0 ]; then
    echo "NG: failed to link with GCC: $expr."
    exit 1
  fi

  ./tests/tmp/out
  actual=$?
  if [ $actual -ne $expected ]; then
    echo "NG: $expected is expected, but got $actual: $expr."
    exit 1
  fi
}

failed_example() {
  expr=$1

  ./skgo "$expr" > /dev/null 2> /dev/null
  if [ $? -eq 0 ]; then
    echo "NG: unexpectedly succeeded: $expr."
    exit 1
  fi
}

example 'func main() { 42; }' 42
example 'func main() { 35; }' 35
example 'func main() { 0; }' 0

example 'func main() { 11 + 22; }' 33
example 'func main() { 25 - 13; }' 12
example 'func main() { 11 + 18 - 7 + 5; }' 27

example 'func main() { 5 * 8; }' 40
example 'func main() { 48 / 16; }' 3
example 'func main() { 55 / 8; }' 6
example 'func main() { 55 % 8; }' 7
example 'func main() { 15 * 3 / 9 % 3; }' 2

example 'func main() { 2 * 3 - 1 * 4; }' 2
example 'func main() { 2 * (3 - 1) * 4; }' 16

example 'func main() { 12; 34; }' 34
example 'func main() { 1 + 2; 3 * 4; }' 12

example 'func main() { var x; x = 4 * 7; x; }' 28
example 'func main() { var x, y; x = 123; y = 100; x - y; }' 23

example 'func f() { 42; } func main() { f(); }' 42

example 'func f(a) { a * a; } func main() { f(7); }' 49
example 'func f(a, b) { a + b; } func main() { f(4, 7); }' 11
example 'func f(a, b, c, d, e, f) { a + b + c + d + e + f; } func main() { f(1, 2, 3, 4, 5, 6); }' 21

example 'func main() { var x; x = 10; if 0 { x = 20; }; x; }' 10
example 'func main() { var x; x = 10; if 1 { x = 20; }; x; }' 20
example 'func main() { var x; x = 10; if 0 { x = 20; } else { x = 30; }; x; }' 30
example 'func main() { var x; x = 10; if 1 { x = 20; } else { x = 30; }; x; }' 20

example 'func main() { return 10; return 20; }' 10

failed_example 'func main() { 123 456; }'
failed_example 'func main() { 2 * (3 + 4; }'
failed_example 'func main() 123'
failed_example 'func main() { 123;'
failed_example 'func main() { 123 }'
failed_example 'func main() { 123 = 456; }'
failed_example 'func main() { var x, x; }'
failed_example 'func main() { x = 123; }'
failed_example 'func () { 123; }'
failed_example 'func 123() { 123; }'
failed_example 'func { 123; }'
failed_example 'func main { 123; }'
failed_example 'func f(a, b, c, d, e, f, g) { 42; }'

echo 'OK'
exit 0
