all:
	make skgo

skgo: main.go
	go build -o skgo main.go

.PHONY: test
test: skgo
	bash tests/test.sh

.PHONY: clean
clean:
	rm main.go
