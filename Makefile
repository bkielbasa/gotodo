
migrate:
	@migrate -path=migrations/ -database postgres://postgres:password@localhost:5432/todos?sslmode=disable up

