module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.41.11 // indirect
	github.com/certonid/certonid v0.0.0-20210924154530-c86d465bf9df
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/rs/zerolog v1.25.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.9.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace github.com/certonid/certonid/adapters => ./../adapters

replace github.com/certonid/certonid/proto => ./../proto

replace github.com/certonid/certonid/utils => ./../utils
