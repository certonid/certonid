module github.com/certonid/certonid/cli

go 1.12

require (
	github.com/ScaleFT/sshkeys v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.18.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.20.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.29.1 // indirect
	github.com/certonid/certonid v0.0.0-20221210120649-943bad14ed4b
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.29.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	golang.org/x/crypto v0.6.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
