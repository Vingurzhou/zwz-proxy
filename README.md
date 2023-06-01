# zwz-proxy
自用代理

##  开发
```shell
protoc --proto_path=./proto \
  --proto_path=/Users/zhouwenzhe/go/pkg/mod/github.com/cosmos/gogoproto@v1.4.8/protobuf \
  --proto_path=/Users/zhouwenzhe/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
  --go_out=./ \
  --go_opt=paths=import \
  --go-grpc_out=./ \
  --go-grpc_opt=paths=import \
  --grpc-gateway_out=./ \
  --grpc-gateway_opt=paths=import \
  --grpc-gateway_opt=logtostderr=true \
  --openapiv2_out=./docs/static \
  --openapiv2_opt=allow_merge=true \
  --openapiv2_opt=merge_file_name=openapi \
  ./proto/helloworld.proto

go mod tidy
kill $(lsof -t -i:9090)
kill $(lsof -t -i:8081)

go run cmd/rpc/main.go & go run cmd/api/main.go
```

##  部署
```shell
env GOOS=linux GOARCH=amd64 go build -o build/zpr cmd/rpc/main.go
env GOOS=linux GOARCH=amd64 go build -o build/zpa cmd/api/main.go
scp -r /Users/zhouwenzhe/src/zwz-proxy/build root@108.160.138.133:/root/zwz
```
