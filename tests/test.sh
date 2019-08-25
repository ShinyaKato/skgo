mkdir -p tests/tmp

./skgo > tests/tmp/out.s
gcc -o tests/tmp/out tests/tmp/out.s
./tests/tmp/out

if [ $? -ne 42 ]; then
  echo 'NG'
  exit 1
fi

echo 'OK'
exit 0
