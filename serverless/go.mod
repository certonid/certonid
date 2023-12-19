module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.42.0
	github.com/aws/aws-sdk-go-v2/config v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.27.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.49.5 // indirect
	github.com/certonid/certonid v0.0.0-20231202144713-e600898651a6
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/rs/zerolog v1.31.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.18.2
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.17.0
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
