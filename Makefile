DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres-learning1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres-learning1 createdb --username=root --owner=root simple_bank

drobdb:
	docker exec it postgres-learning1 dropdb simple_bank

migrateuprds:
	migrate -path db/migration -database "postgresql://root:ufWuwSPmJPnyBNQPieXS@simplebank1.cd3dseywmjx4.eu-west-2.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

mysql: 
	docker run --name mysql-learning1 -p 3306:3306 -e MYSQL_DATABASE=simple_bank -e MYSQL_USER=banker -e MYSQL_PASSWORD=secret -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go backend_masterclass/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema: 
	dbml2sql doc/db.dbml -o doc/simple_bank_schema.sql

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test mysql server mock db_docs db_schema proto