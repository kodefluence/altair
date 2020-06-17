# Altair - Light Weight and Robust API Gateway

[![Coverage Status](https://coveralls.io/repos/github/codefluence-x/altair/badge.svg?branch=master)](https://coveralls.io/github/codefluence-x/altair?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/codefluence-x/altair)](https://goreportcard.com/report/github.com/codefluence-x/altair)
[![GitHub issues](https://img.shields.io/github/issues/codefluence-x/altair)](https://github.com/codefluence-x/altair/issues)
[![GitHub forks](https://img.shields.io/github/forks/codefluence-x/altair)](https://github.com/codefluence-x/altair/network)
[![GitHub stars](https://img.shields.io/github/stars/codefluence-x/altair)](https://github.com/codefluence-x/altair/stargazers)
[![GitHub license](https://img.shields.io/github/license/codefluence-x/altair)](https://github.com/codefluence-x/altair/blob/master/LICENSE)

Open source API gateway written in Go. Created to be lightweight, simple, fast, reliable, cross platform and programming language agnostic.

![Altair Architecture Diagram](https://user-images.githubusercontent.com/20650401/79699757-a2337d00-82bb-11ea-8103-25e6917545bd.png)

## Documentation

### Plugin API Documentation

[Plugin API Documentation in Postman](https://documenter.getpostman.com/view/3666028/SzmcZJ79?version=latest#b870ae5a-b305-4016-8155-4899af1f26b1)

## How to Use

> TBD

## How to Contribute

### Installation

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
  - [x] Downstream Plugins
    - [x] Oauth Token Checking
    - [ ] Response Caching
- [x] Metric & Monitoring
  - [x] Prometheus
- [x] Logging
  - [x] Stdout
- [ ] Plugins
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
      - [x] Revoke Access Token
  - [ ] Response Caching
    - [ ] Route Config
    - [ ] API for deleting the cache
  - [ ] JWT
