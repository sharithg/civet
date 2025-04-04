.PHONY: migrate-up migrate-down sqlc

migrate-up:
	migrate -source file://db/migrations -database "postgresql://postgres:postgres@localhost:5432/receipt?sslmode=disable" up

migrate-down:
	migrate -source file://db/migrations -database "postgresql://postgres:postgres@localhost:5432/receipt?sslmode=disable" down

sqlc:
	sqlc generate