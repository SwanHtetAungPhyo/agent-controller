include .env
export

DB_URL := postgresql://$(APP_DATABASE_USERNAME):$(APP_DATABASE_PASSWORD)@$(APP_DATABASE_HOST):$(APP_DATABASE_PORT)/$(APP_DATABASE_NAME)?sslmode=$(APP_DATABASE_SSL_MODE)

# Migration commands
migrate-up:
	@echo "Running migrations up..."
	@migrate -path ./db/migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path ./db/migrations -database "$(DB_URL)" down

migrate-down-one:
	@echo "Rolling back one migration..."
	@migrate -path ./db/migrations -database "$(DB_URL)" down 1

migrate-up-one:
	@echo "Running one migration..."
	@migrate -path ./db/migrations -database "$(DB_URL)" up 1

migrate-force:
	@echo "Forcing migration version to $(VERSION)..."
	@migrate -path ./db/migrations -database "$(DB_URL)" force $(VERSION)

migrate-version:
	@echo "Current migration version:"
	@migrate -path ./db/migrations -database "$(DB_URL)" version

migrate-create:
	@echo "Creating new migration: $(NAME)..."
	@migrate create -ext sql -dir ./db/migrations -seq $(NAME)