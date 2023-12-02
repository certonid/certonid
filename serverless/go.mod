module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2/config v1.25.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.27.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.49.2 // indirect
	github.com/certonid/certonid v0.0.0-20231102211808-97f23a4bdb4c
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.31.0
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/viper v1.17.0
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.16.0
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
