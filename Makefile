postgres:
	docker run --name postgres-learning1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres-learning1 createdb --username=root --owner=root simple_bank

drobdb:
	docker exec it postgres-learning1 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

mysql: 
	docker run --name mysql-learning1 -p 3306:3306 -e MYSQL_DATABASE=simple_bank -e MYSQL_USER=banker -e MYSQL_PASSWORD=secret -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mysql server