postgres:
	docker run --name mypostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres 

createdb:
	docker exec -it mypostgres createdb --username=root --owner=root task_manager_db

dropdb:
	docker exec -it mypostgres dropdb task_manager_db

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/task_manager_db?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/task_manager_db?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown test	