BINARY_NAME=subscriptionapp.exe

## build: builds all binaries
build:
	@go build -o ${BINARY_NAME} ./cmd/web
	@echo back end built!

run: build
	@echo Starting...
	@${BINARY_NAME}
	@echo back end started!

clean:
	@echo Cleaning...
	@DEL ${BINARY_NAME}
	@go clean
	@echo Cleaned!

start: run

stop:
	@echo "Stopping..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped back end

restart: stop start

test:
	@echo "Testing..."
	go test -v ./...
