module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/aws/aws-sdk-go-v2/config v1.27.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.30.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.53.3 // indirect
	github.com/certonid/certonid v0.0.0-20240121205526-99393e80ed9a
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/zerolog v1.32.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.18.2
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.21.0
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
