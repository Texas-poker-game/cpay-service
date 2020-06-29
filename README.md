#### release

```bash
go build -o bin/cpay-mac
env GOOS=linux GOARCH=amd64 go build -o bin/cpay-linux
```

#### run

```
./bin/cpay-mac

GIN_MODE=release GO_ENV=test ./bin/cpay-linux &
```
