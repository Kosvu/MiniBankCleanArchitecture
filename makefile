migrate-up:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank?sslmode=disable" up
migrate-down:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank?sslmode=disable" down

