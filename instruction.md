### Clone repo
```bash
git clone github.com/GitProger/go-telecom-2025
```

### Run with docker
1. Build
```bash
    make docker-build
```

2. Run with mounted input files:
```bash
    ./docker.sh "./sunny_5_skiers/config.json" "./sunny_5_skiers/events"
```

### Bare run
```bash
make build
./go-telecom-2025 "./sunny_5_skiers/config.json" "./sunny_5_skiers/events"

make test-input-1 # test on examples from `./sunny_5_skiers/sample`
make test-input-2
```

#### P.S.
You can uncomment lines 64 and 65 in main.go if you want to use more basic and simple version:
```go
func main() { // interactive
	mainSimple()
	return
```
