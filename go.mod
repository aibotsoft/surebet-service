module github.com/aibotsoft/surebet-service

go 1.14

require (
	github.com/aibotsoft/gen v0.0.0-20200413085542-106638a26d56
	github.com/aibotsoft/micro v0.0.0-20200411114812-ccef30d833e9
	google.golang.org/grpc v1.28.0
)

replace github.com/aibotsoft/micro => ../micro
replace github.com/aibotsoft/gen => ../gen
