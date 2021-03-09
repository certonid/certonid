# Certonid ![Build & Release](https://github.com/certonid/certonid/workflows/Build%20&%20Release/badge.svg)

Certonid is a Serverless SSH Certificate Authority.

Consists of two parts: CLI and serverless function.

![certonid-schema](https://user-images.githubusercontent.com/98444/109483362-cdfcc300-7a87-11eb-8453-fa9d2c6d930a.png)

## Releases

[Download latest releases](https://github.com/certonid/certonid/releases)

## Documentation

[All information published at Wiki page](https://github.com/certonid/certonid/wiki)

## Articles

 - [[English] Certonid — the SSH Certificate Authority Deployed on AWS Lambda](https://blog.mailtrap.io/certonid/)
 - [[Russian] Certonid — SSH центр сертификации, который работает на AWS Lambda](https://dou.ua/lenta/articles/certonid-ssh/)

## AWS Terraform module

To simplify setup on AWS, you can use [Certonid AWS Terraform module](https://registry.terraform.io/modules/certonid/certonid/aws/latest)

## Dev build

```shell
$ cd serverless && GOOS=linux go build -o serverless main.go
$ cd ..
$ cd cli && go build -o certonid main.go
```

## Roadmap

 - [ ] Add tests
 - [ ] Improve documentation
 - [ ] Website and video
 - [ ] Use serverless framework to automate deploy and setup
 - [ ] Support AWS S3 for SSH CA key

