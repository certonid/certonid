module github.com/certonid/certonid/cli

go 1.12

require (
	github.com/ScaleFT/sshkeys v0.0.0-20200327173127-6142f742bca5
	github.com/aws/aws-sdk-go-v2/config v1.15.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.17.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.23.1 // indirect
	github.com/certonid/certonid v0.0.0-20220603151313-fc41bb1b5128
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.26.1
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.12.0
	github.com/subosito/gotenv v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
