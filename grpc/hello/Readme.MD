Generate hello.pb.go using this command 

```
protoc --go_out=plugins=grpc:./hello --proto_path=./hello ./hello/hello.proto
```