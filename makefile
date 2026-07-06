include .env
export

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-up-test:
	migrate -path migrations -database "$(DATABASE_URL_TEST)" up

migrate-down-test:
	migrate -path migrations -database "$(DATABASE_URL_TEST)" down
