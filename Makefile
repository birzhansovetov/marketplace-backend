DB_URL=postgres://postgres:postgres@localhost:5432/marketplace_db?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)