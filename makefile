include .envrc

migrations_path= ./cmd/migrate/migration

.PHONY: migrate-create

migration: 
	@migrate create -seq -ext sql -dir $(migrations_path) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(migrations_path) -database $(DB_ADDR) up


.PHONY: migrate-down
migrate-down:
	@migrate -path $(migrations_path) -database $(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))


.PHONY: seed
seed:
	@go run ./cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt