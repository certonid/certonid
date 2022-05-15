# Certonid ![Build & Release](https://github.com/certonid/certonid/workflows/Build%20&%20Release/badge.svg)

Certonid is a Serverless SSH Certificate Authority.

Consists of two parts: CLI and serverless function.

![certonid-schema](https://user-images.githubusercontent.com/98444/109483362-cdfcc300-7a87-11eb-8453-fa9d2c6d930a.png)

## Releases

[Download latest releases](https://github.com/certonid/certonid/releases)

For Mac OS or Linux you can use Homebrew tap:

```bash
brew install certonid/tap/certonid
```

## Documentation

[All information published at Wiki page](https://github.com/certonid/certonid/wiki)

## Articles

 - [[English] Certonid — the SSH Certificate Authority Deployed on AWS Lambda](https://mailtrap.io/blog/certonid/)
 - [[Russian] Certonid — SSH центр сертификации, который работает на AWS Lambda](https://dou.ua/lenta/articles/certonid-ssh/)

## AWS Terraform module

To simplify setup on AWS, you can use [Certonid AWS Terraform module](https://registry.terraform.io/modules/certonid/certonid/aws/latest)

## Binaries security

All archives signed by gpg key. You can check its by downloading with archive it `.sig` file and verify signature (example with linux x86 cli archive):

```bash
$ gpg --verify certonid_0.8.2_Linux_x86_64.tar.gz.sig certonid_0.8.2_Linux_x86_64.tar.gz
gpg: Signature made Wed Mar 10 11:02:40 2021 EET
gpg:                using RSA key 6894D468143A22469D6603D1E44200219869E71E
gpg: Good signature from "leopard apps <leopard.not.a+apps@gmail.com>"
```

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

