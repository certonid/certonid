module github.com/certonid/certonid/serverless

go 1.12

require (
	github.com/aws/aws-lambda-go v1.32.1
	github.com/aws/aws-sdk-go-v2/config v1.15.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.23.4 // indirect
	github.com/certonid/certonid v0.0.0-20220701112256-1ce95aeae277
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/rs/zerolog v1.27.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/viper v1.12.0
	github.com/subosito/gotenv v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/sys v0.0.0-20220708085239-5a0f0661e09d // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
