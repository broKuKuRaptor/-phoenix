.PHONY: build-accounts clean

build-accounts:
	go build -o bin/accounts ./cmd/accounts

clean:
	rm -rf bin/
