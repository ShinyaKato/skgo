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

example '42' 42
example '35' 35
example '0' 0

example '11 + 22' 33
example '25 - 13' 12
example '11 + 18 - 7 + 5' 27

example '5 * 8' 40
example '48 / 16' 3
example '55 / 8' 6
example '55 % 8' 7
example '15 * 3 / 9 % 3' 2

example '2 * 3 - 1 * 4' 2
example '2 * (3 - 1) * 4' 16

failed_example '123 456'
failed_example '2 * (3 + 4'

echo 'OK'
exit 0
