BINARY := pgxdb-example
DB_URL  := postgres://demo:demo@localhost:5432/demo

.PHONY: up down build run clean

up:
	docker compose up -d

down:
	docker compose down

build:
	glisp build .

run: up build
	DATABASE_URL=$(DB_URL) ./$(BINARY)

clean:
	rm -f $(BINARY)
	docker compose down -v
