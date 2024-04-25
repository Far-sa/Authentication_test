postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root PGU

dropdb:
	docker exec -it postgres12 dropdb PGU

migrateup:
	migrate -path infrastructure/db/migrations -database "postgresql://root:password@localhost:5432/PGU?sslmode=disable" -verbose up

migratedown:
	migrate -path infrastructure/db/migrations -database "postgresql://root:password@localhost:5432/PGU?sslmode=disable" -verbose down

sqlc:
	sqlc generate


protoc:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=.\
		ports/protobuf/grpc/*.proto \
		
	# protoc --go_out=./ports/protobuf/grpc \
	# 	./ports/protobuf/grpc/*.proto \
	# protoc --go_opt=module=$(GO_MODULE) --go_out=./ports/protobuf/grpc\
	# 	./ports/protobuf/grpc/*.proto \
	




.PHONY: postgres createdb dropdb migrateup migratedown sqlc protoc