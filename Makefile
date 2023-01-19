migration-up:
	migrate -path ./migration/postgres/ -database 'postgres://samandar:saman107@localhost:5432/catalog?sslmode=disable' up 


run:
	go run cmd/main.go
swag:
	swag init -g api/api.go -o api/docs

