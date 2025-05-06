out = go-telecom-2025

build: clean
	go build -o $(out) ./cmd
test:
	go test ./internal/...
run:
	./$(out) "./sunny_5_skiers/config.json" "./sunny_5_skiers/events"

test-input-1:
	./$(out) "./sunny_5_skiers/sample/config.json" "./sunny_5_skiers/sample/events"
test-input-2:
	./$(out) "./sunny_5_skiers/sample/config.json" "./sunny_5_skiers/sample/disqual" 

clean:
	[ ! -f $(out) ] || rm $(out)

all: build test run

docker-build:
	docker build -t go-telecom-2025-docker .
