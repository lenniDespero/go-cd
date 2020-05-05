PROJECTNAME="Go-CD-script"

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build application for your system
build:
		GCO_ENABLED=0 go build -o ./bin/deployer ./cmd/main.go
## build-linux: Build application for linux
build-linux:
		GCO_ENABLED=0 GOOS=linux go build -o ./bin/deployer-lin ./cmd/main.go
## build-mac: Build application for mac
build-mac:
		GCO_ENABLED=0 GOOS=darwin go build -o ./bin/deployer-mac ./cmd/main.go
## build-win: Build application for windows
build-win:
		GCO_ENABLED=0 GOOS=windows go build -o ./bin/deployer-win ./cmd/main.go
