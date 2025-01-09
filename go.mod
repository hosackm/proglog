module github.com/hosackm/proglog

go 1.23.2

replace github.com/hosackm/proglog/api => ./api

require (
	github.com/go-chi/chi/v5 v5.2.0
	github.com/hosackm/proglog/api v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	github.com/tysonmote/gommap v0.0.3
	google.golang.org/protobuf v1.36.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
