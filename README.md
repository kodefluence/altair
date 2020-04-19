# Altair - Light Weight and Robust API Gateway

Open source API gateway written in Go.

![Altair Architecture Diagram](https://user-images.githubusercontent.com/20650401/79699757-a2337d00-82bb-11ea-8103-25e6917545bd.png)

## Documentation

> TBD

## How to Use

> TBD

## How to Contribute

### Instalation

#### Prerequisites

1. Go version 1.13 or higher
2. Mysql

#### How To

1. Clone this repo
2. `go run altair.go migrate`
3. `go run altair.go run`
4. Read [CONTRIBUTING.md](https://github.com/insomnius/code-geek/blob/master/CONTRIBUTING.md)

## Feature

- [ ] Request Forwarder
  - [x] Route Compiler
  - [x] Route Generator Forwader
  - [ ] Authorization Checking
    - [ ] Oauth
- [ ] Oauth Authorization
  - [ ] CRUD Oauth Application
    - [x] Create
    - [x] List
    - [x] One
    - [ ] Update
  - [x] Authorization
    - [x] Authorization Code Grant
  - [ ] Access Token
    - [x] Access Token Implicit Request for Confidential Application
    - [ ] Access Token Code Grant Flow
    - [ ] Refresh Token
      - [ ] Access Token Refresh Token Flow
      - [ ] Refresh Token Generation
  - [ ] Request Validation
- [ ] Response Caching
  - [ ] Route Config
  - [ ] API for deleting the cache
- [ ] JWT
