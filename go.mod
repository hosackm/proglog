module github.com/hosackm/proglog

go 1.23.2

replace github.com/hosackm/proglog/api => ./api

require (
	github.com/casbin/casbin v1.9.1
	github.com/go-chi/chi/v5 v5.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/hosackm/proglog/api v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	github.com/tysonmote/gommap v0.0.3
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.2
)

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
