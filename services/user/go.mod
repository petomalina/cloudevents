module github.com/flowup/petermalina/services/user

go 1.14

replace github.com/flowup/petermalina/apis/go-sdk => ../../apis/go-sdk

require (
	cloud.google.com/go/firestore v1.3.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/blendle/zapdriver v1.3.1
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/flowup/petermalina/apis/go-sdk v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/viper v1.7.0
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200728102440-3e129f6d46b1 // indirect
	golang.org/x/tools v0.0.0-20200729041821-df70183b1872 // indirect
	google.golang.org/genproto v0.0.0-20200729003335-053ba62fc06f // indirect
	google.golang.org/grpc v1.30.0
)
