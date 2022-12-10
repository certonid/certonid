module github.com/certonid/certonid/cli

go 1.12

require (
	github.com/ScaleFT/sshkeys v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.26.1 // indirect
	github.com/certonid/certonid v0.0.0-20221005080213-b23329b7b18c
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/rs/zerolog v1.28.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.14.0
	golang.org/x/crypto v0.4.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
