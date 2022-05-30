postgres:
	 docker run --name postgres12 -p 5432:5432  --network bank-network  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	 docker exec -it postgres12  dropdb simple_bank	

migrateup:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1	 

sqlc:
	 sqlc generate
test: #run test including test coverage
	 go test -v -cover ./...

server:
	  go run main.go  

viper:
	 go get github.com/spf13/viper
	 
mock: 
	mockgen -package mockdb -destination db/mock/store.go  github.com/Franklynoble/bankapp/db/sqlc Store	
  
.PHONY: postgres createdb dropdb migrateup migratedown  migratedown1  migrateup1 sqlc test server viper