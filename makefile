migrate-up:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank?sslmode=disable" up
migrate-down:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank?sslmode=disable" down

migrate-up-test:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank_test?sslmode=disable" up
migrate-down-test:
	migrate -path migrations -database "postgres://postgres:usif775shakh@localhost:5432/minibank_test?sslmode=disable" down



