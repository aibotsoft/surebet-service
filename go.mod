module github.com/aibotsoft/surebet-service

go 1.14

require (
	github.com/aibotsoft/decimal v0.0.0-20200424173126-4bc23d40885f
	github.com/aibotsoft/gen v0.0.0-20200531091936-c4d5d714bf82
	github.com/aibotsoft/micro v0.0.0-20200421094132-4cf4004de76e
	github.com/denisenkom/go-mssqldb v0.0.0-20200428022330-06a60b6afbbc
	github.com/dgraph-io/ristretto v0.0.2
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jmoiron/sqlx v1.2.0
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5 // indirect
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	google.golang.org/genproto v0.0.0-20200424135956-bca184e23272 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/aibotsoft/micro => ../micro

replace github.com/aibotsoft/gen => ../gen
