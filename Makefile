CLI_BINARY_NAME=tuiGhostOrange
SRV_BINARY_NAME=GhostOrangeServer

build:
# Client
 GOARCH=amd64 GOOS=darwin go build -o ${CLI_BINARY_NAME}-darwin ./ghostorange/cmd/client/main.go
 GOARCH=amd64 GOOS=linux go build -o ${CLI_BINARY_NAME}-linux ./ghostorange/cmd/client/main.go
 GOARCH=amd64 GOOS=windows go build -o ${CLI_BINARY_NAME}-windows ./ghostorange/cmd/client/main.go

#  # Server
#  GOARCH=amd64 GOOS=darwin go build -o ${SRV_BINARY_NAME}-darwin ./ghostorange/cmd/client/main.go
#  GOARCH=amd64 GOOS=linux go build -o ${SRV_BINARY_NAME}-linux ./ghostorange/cmd/client/main.go
#  GOARCH=amd64 GOOS=windows go build -o ${SRV_BINARY_NAME}-windows ./cmd/ghostorange/main.go/main.go

run: build
	./${CLI_BINARY_NAME}
# build
# ./${SRV_BINARY_NAME}

clean:
	go clean
	rm ${CLI_BINARY_NAME}-darwin
	rm ${CLI_BINARY_NAME}-linux
	rm ${CLI_BINARY_NAME}-windows

# rm ${SRV_BINARY_NAME}-darwin
#  rm ${SRV_BINARY_NAME}-linux
#  rm ${SRV_BINARY_NAME}-windows