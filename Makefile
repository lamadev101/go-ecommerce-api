GOBUILD=go build
GOTEST=go test

all:clean stop build_server
	./ecommerce-api
build_server:
	$(GOBUILD) -v .
clean:
	rm -f ./ecommerce-api
stop:
	pkill ecommerce-api || true
test:
	cd helper && $(GOTEST) -v .