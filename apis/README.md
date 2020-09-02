##
Download binary `protoc` and copy to system path.
`https://github.com/protocolbuffers/protobuf/releases`
```
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.12.4/protoc-3.12.4-linux-x86_64.zip
unzip protoc-3.12.4-linux-x86_64.zip
sudo mv proto/bin/protoc /usr/local/bin
```


## Go SDK
In `go-sdk` run:

```
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
go get github.com/envoyproxy/protoc-gen-validate`
```