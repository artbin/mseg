module github.com/artbin/mseg/pflist

go 1.25

replace github.com/artbin/mseg/segment => ../segment

require (
	github.com/artbin/mseg/segment v0.0.2
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v0.0.1 // published tag does not match the checksum database
