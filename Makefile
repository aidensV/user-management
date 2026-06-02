# Makefile - Migration terpisah, run tidak migrasi ulang

.PHONY: help run swagger migrate seed

help:
	@echo "========================================="
	@echo "User Management Service"
	@echo "========================================="
	@echo ""
	@echo "  make run        - Run application (normal, no migration)"
	@echo "  make migrate    - Run AutoMigrate (create/update tables)"
	@echo "  make seed       - Seed admin user only"
	@echo "  make swagger    - Generate Swagger docs"
	@echo "  make setup      - First time setup (migrate + seed)"
	@echo ""

run:
	@echo ">>> Running User Management Service..."
	go run ./cmd/main.go

migrate:
	@echo ">>> Running AutoMigrate..."
	go run ./cmd/migrate/main.go
	@echo ">>> Migration complete!"

seed:
	@echo ">>> Seeding admin user..."
	go run ./cmd/seed/main.go
	@echo ">>> Seed complete! Admin: admin_user_m@yopmail.com / admin123"

swagger:
	@echo ">>> Generating Swagger documentation..."
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal --parseDepth 3 --generatedTime
	@echo ">>> Swagger generated at ./docs/index.html"

setup: migrate seed
	@echo ">>> Setup complete! Now run: make run"