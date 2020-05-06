module github.com/aibotsoft/surebet-service

go 1.14

require (
	github.com/aibotsoft/decimal v0.0.0-20200424173126-4bc23d40885f
	github.com/aibotsoft/gen v0.0.0-20200421094132-60f0f2019f16
	github.com/aibotsoft/micro v0.0.0-20200421094132-4cf4004de76e
	github.com/dgraph-io/ristretto v0.0.2
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5 // indirect
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd // indirect
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	google.golang.org/genproto v0.0.0-20200424135956-bca184e23272 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/aibotsoft/micro => ../micro

replace github.com/aibotsoft/gen => ../gen
