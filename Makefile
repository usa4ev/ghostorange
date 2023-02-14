CLI_SRC=./cmd/client/main.go
SRV_SRC=./cmd/github.com/usa4ev/ghostorange/main.go
CLI_BINARY_NAME=tuiGOrange
SRV_BINARY_NAME=GOrangeServer
BIN_PATH=./bin
BUILDDATE=`date +%Y.%m.%d`
LDFLAGS=-ldflags "-X 'github.com/usa4ev/ghostorange/internal/app/tui/appinfo.BuildDate=$(BUILDDATE)'"

 # Build server
build-srv-linux:
	GOARCH=amd64 GOOS=linux go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-linux $(SRV_SRC)

build-srv-darwin:
	GOARCH=amd64 GOOS=darwin go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-darwin $(SRV_SRC)

build-srv-windows:
	GOARCH=amd64 GOOS=windows go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-windows $(SRV_SRC)

# Build client
build-tui-windows:
	GOARCH=amd64 GOOS=windows go build $(LDFLAGS) -o $(BIN_PATH)/${CLI_BINARY_NAME}-windows $(CLI_SRC)

build-tui-darwin:
	GOARCH=amd64 GOOS=darwin go build $(LDFLAGS) -o $(BIN_PATH)/${CLI_BINARY_NAME}-darwin $(CLI_SRC)

build-tui-linux:
	GOARCH=amd64 GOOS=linux go build $(LDFLAGS) -o $(BIN_PATH)/${CLI_BINARY_NAME}-linux $(CLI_SRC) 

# Run server
run-srv-linux: build-srv-linux
	$(BIN_PATH)/${SRV_BINARY_NAME}-linux -c ./configs/srv.json

run-srv-darwin: build-srv-darwin
	$(BIN_PATH)/${SRV_BINARY_NAME}-darwin -c ./configs/srv.json

run-srv-windows: build-srv-windows
	$(BIN_PATH)/${SRV_BINARY_NAME}-windows -c ./configs/srv.json

# Run client
run-tui-linux: build-tui-linux
	$(BIN_PATH)/${CLI_BINARY_NAME}-linux -a localhost:8080 -l log.txt

run-tui-darwin: build-tui-darwin
	$(BIN_PATH)/${CLI_BINARY_NAME}-darwin -a localhost:8080 -l log.txt

run-tui-windows: build-tui-windows
	$(BIN_PATH)/${CLI_BINARY_NAME}-windows -a localhost:8080 -l log.txt

	
test:
	go test -v "./..."

