all:
	make skgo

skgo: token/*.go node/*.go lexer/*.go parser/*.go main.go
	go build -o skgo

.PHONY: test
test: skgo
	bash tests/test.sh

.PHONY: clean
clean:
	rm skgo
