# Certonid [![Build Status](https://travis-ci.com/certonid/certonid.svg?branch=master)](https://travis-ci.com/certonid/certonid)

Certonid is a Serverless SSH Certificate Authority.

Consists of two parts: CLI and serverless function.

## Releases

[Download latest releases](https://github.com/certonid/certonid/releases)

## Documentation

[All information published at Wiki page](https://github.com/certonid/certonid/wiki)

## Dev build

```shell
$ cd serverless
$ GOOS=linux go build -o serverless main.go
$ cd ..
$ cd cli
$ go build -o certonid main.go
```

## Roadmap

 - [ ] Add tests
 - [ ] Improve documentation
 - [ ] Website and video
 - [ ] Use serverless framework to automate deploy and setup
 - [ ] Support AWS S3 for SSH CA key
 - [ ] Support Google Cloud functions
 - [ ] Support Google Cloud Storage for SSH CA key
 - [ ] Support Azure Functions
 - [ ] Support Apache OpenWhisk

