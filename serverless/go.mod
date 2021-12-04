module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.42.19 // indirect
	github.com/certonid/certonid v0.0.0-20211028211740-e65e6f34ffef
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/rs/zerolog v1.26.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.9.0
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
