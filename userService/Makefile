
sqlc:
	sqlc generate


protoc:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=.\
		ports/protobuf/grpc/*.proto \
		
	# protoc --go_out=./ports/protobuf/grpc \
	# 	./ports/protobuf/grpc/*.proto \
	# protoc --go_opt=module=$(GO_MODULE) --go_out=./ports/protobuf/grpc\
	# 	./ports/protobuf/grpc/*.proto \
	




.PHONY: sqlc protoc