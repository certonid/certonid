# Certonid ![Build & Release](https://github.com/certonid/certonid/workflows/Build%20&%20Release/badge.svg)

Certonid is a Serverless SSH Certificate Authority.

Consists of two parts: CLI and serverless function.

![certonid-schema](https://user-images.githubusercontent.com/98444/109483362-cdfcc300-7a87-11eb-8453-fa9d2c6d930a.png)

## Releases

[Download latest releases](https://github.com/certonid/certonid/releases)

For Mac OS or Linux you can use [Homebrew tap](https://docs.brew.sh/Taps):

```bash
brew install certonid/tap/certonid
```

or another way:

```bash
brew tap certonid/tap
brew install certonid
```

## Documentation

[All information published at Wiki page](https://github.com/certonid/certonid/wiki)

## Articles

 - [[English] Certonid — the SSH Certificate Authority Deployed on AWS Lambda](https://mailtrap.io/blog/certonid/)
 - [[Russian] Certonid — SSH центр сертификации, который работает на AWS Lambda](https://dou.ua/lenta/articles/certonid-ssh/)

## AWS Terraform module

To simplify setup on AWS, you can use [Certonid AWS Terraform module](https://registry.terraform.io/modules/certonid/certonid/aws/latest)

## Binaries security

`checksum.txt` signed by gpg key. You can check its by downloading with archive it `.sig` file and verify signature:

```bash
$ gpg --verify checksums.txt.sig checksums.txt
gpg: Signature made Fri Jul 22 17:24:40 2022 EEST
gpg:                using RSA key 36E7986334C6DE2B41A29537A77A9969BEFF93AE
gpg: Good signature from "Certonid Sign Key (Certonid Sign Key) <leopard.not.a+certonid@gmail.com>"
```

Each archive have [SBOM file](https://www.ntia.gov/SBOM).

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

