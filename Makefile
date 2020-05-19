.PHONY: test
PROJECTNAME="Go-CD-script"

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build application for your system
build:
		GCO_ENABLED=0 go build -o ./bin/deployer ./cmd/main.go
## test: run integration tests
test:
	cd deployment ;\
	test_status=0 ;\
	docker-compose up --build --exit-code-from tests --abort-on-container-exit|| test_status=$$? ;\
	cd .. ;\
	exit $$test_status ;\
