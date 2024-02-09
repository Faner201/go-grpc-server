include .env

PROJECTNAME=$(shell basename "$(PWD)")

##- proto_gen : generate proto files
proto_gen:
	protoc -I $(MAIN_DIRECTORY) $(PROTO_FILES_DIRECTORY) --go_out=$(MAIN_DIRECTORY)/$(OUT_GO_DIRECTORY)  --go-grpc_out=$(MAIN_DIRECTORY)/$(OUT_GO_DIRECTORY) 
.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo