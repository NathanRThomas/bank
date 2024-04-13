
# go params
GOCMD=go
GOBUILD=$(GOCMD) build -buildvcs=false
GOTEST=$(GOCMD) test -run
GOPATH=/usr/local/bin
DIR=$(shell pwd)

# normal entry points
build:
	clear
	@echo "building bank..."
	@$(GOBUILD) -o $(GOPATH)/bank ./cmd
	
update:
	clear
	@echo "updating dependencies..."
	@$(GOCMD) get -u -t ./...
	@$(GOCMD) mod tidy 

test:
	@clear 
	@echo "QA testing..."
	@$(GOTEST) QA ./...

run: build
run:
	$(GOPATH)/bank

deploy: build 
deploy:
	@systemctl restart bank
	