packages =  \
  ./infrastructure/config \
  ./internal/post \

.PHONY: test
test: 
	@$(foreach package,$(packages), \
	  set -e; \
	  go test -coverprofile $(package)/cover.out -covermode=count $(package);)

.PHONY: cover
cover: test
	@echo "mode: count" > cover-all.out 
	@$(foreach package,$(packages), \
  	  tail -n +2 $(package)/cover.out >> cover-all.out;)
	@gocover-cobertura < cover-all.out > cover-cobertura.xml
	@go tool cover -func=cover-all.out | tail -n 1

# Start all services
.PHONY: up
up:
	docker-compose up -d

.PHONY: lint
lint:
	cd ./src && golangci-lint run ./... --verbose --no-config --out-format checkstyle > golangci-lint.out;