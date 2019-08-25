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

example '42' 42
example '35' 35
example '0' 0

echo 'OK'
exit 0
