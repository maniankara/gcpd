GOARCH = amd64
SERVER_BINARY = gcpd
CLIENT_BINARY = gcp
BUILD_DIR=$(shell pwd)
BIN_DIR=${BUILD_DIR}/bin

# Build the project
all: linux darwin windows

clean:
	@rm -rf ${BIN_DIR}

init: clean
	@mkdir -p ${BIN_DIR}/win
	@mkdir -p ${BIN_DIR}/linux
	@mkdir -p ${BIN_DIR}/darwin

linux: init
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BIN_DIR}/linux/${SERVER_BINARY} gcp_server/server.go ; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BIN_DIR}/linux/${CLIENT_BINARY} gcp_client/client.go ; \
	cd - >/dev/null

darwin: init
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BIN_DIR}/darwin/${SERVER_BINARY} gcp_server/server.go ; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BIN_DIR}/darwin/${CLIENT_BINARY} gcp_client/client.go ; \
	cd - >/dev/null

windows: init
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BIN_DIR}/win/${SERVER_BINARY}.exe gcp_server/server.go ; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BIN_DIR}/win/${CLIENT_BINARY}.exe gcp_client/client.go ; \
	cd - >/dev/null