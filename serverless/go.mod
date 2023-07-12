module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2/config v1.18.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.23.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.37.0 // indirect
	github.com/certonid/certonid v0.0.0-20230609102143-4e2d45cf5237
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.29.1
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.16.0
	golang.org/x/crypto v0.11.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
