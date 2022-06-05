module refurbedchallenge/executable

go 1.18

replace refurbedchallenge/notifier => ../notifier

require (
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.7.1
	refurbedchallenge/notifier v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
