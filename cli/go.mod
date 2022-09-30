module github.com/certonid/certonid/cli

go 1.12

require (
	github.com/ScaleFT/sshkeys v0.0.0-20200327173127-6142f742bca5
	github.com/aws/aws-sdk-go-v2/config v1.17.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.18.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.24.6 // indirect
	github.com/certonid/certonid v0.0.0-20220812142328-697c68d9e009
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.28.0
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/term v0.0.0-20220919170432-7a66f970e087 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace (
	github.com/certonid/certonid/adapters => ./../adapters
	github.com/certonid/certonid/proto => ./../proto
	github.com/certonid/certonid/utils => ./../utils
)
