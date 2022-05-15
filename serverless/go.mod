module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.31.1
	github.com/aws/aws-sdk-go v1.44.14 // indirect
	github.com/certonid/certonid v0.0.0-20220414092445-05f6274df14b
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/rs/zerolog v1.26.1
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/viper v1.11.0
	golang.org/x/crypto v0.0.0-20220513210258-46612604a0f9
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20220512140231-539c8e751b99 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
