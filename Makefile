PROJECTNAME="Go-CD-script"

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build application for your system
build:
		GCO_ENABLED=0 go build -o ./bin/deployer ./cmd/main.go
