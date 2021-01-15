module github.com/aibotsoft/surebet-service

go 1.15

require (
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/aibotsoft/gen v0.0.0-20210115125951-7d0a2b681929
	github.com/aibotsoft/micro v0.0.0-20210115130122-86ba91d9a4ad
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/dgraph-io/ristretto v0.0.3
	github.com/go-resty/resty/v2 v2.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/jinzhu/copier v0.2.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	github.com/stretchr/testify v1.6.1
	github.com/vrischmann/envconfig v1.3.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f // indirect
	google.golang.org/grpc v1.35.0
)

//replace github.com/aibotsoft/micro => ../micro

//replace github.com/aibotsoft/gen => ../gen
