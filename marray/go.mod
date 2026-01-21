module github.com/artbin/mlist/marray

go 1.25

require (
	github.com/artbin/mlist/mlist v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/artbin/mlist/segment v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/artbin/mlist/mlist => ../mlist
	github.com/artbin/mlist/segment => ../segment
)
