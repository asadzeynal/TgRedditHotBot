postgres:
	docker run --name postgres15 --network tgreddit-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15.1-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root tgreddit

dropdb:
	docker exec -it postgres15 dropdb tgreddit

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/tgreddit?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/tgreddit?sslmode=disable" -verbose down

sqlc:
	sqlc generate

run:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown run sqlc
