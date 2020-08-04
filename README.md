![Altair - Lightweight and Robust API Gateway](https://user-images.githubusercontent.com/1132660/85085782-0bc00100-b203-11ea-8e89-bbeb7d03a418.png)

<p align="center">
  <a href="https://coveralls.io/github/codefluence-x/altair?branch=master"><img src="https://coveralls.io/repos/github/codefluence-x/altair/badge.svg?branch=master" alt="Coverage Status"></a>
  <a href="https://goreportcard.com/report/github.com/codefluence-x/altair"><img src="https://goreportcard.com/badge/github.com/codefluence-x/altair" alt="Go Report Card"></a>
  <a href="https://github.com/codefluence-x/altair/issues"><img src="https://img.shields.io/github/issues/codefluence-x/altair" alt="GitHub Issues"></a>
  <a href="https://github.com/codefluence-x/altair/network"><img src="https://img.shields.io/github/forks/codefluence-x/altair" alt="GitHub Forks"></a>
  <a href="https://github.com/codefluence-x/altair/stargazers"><img src="https://img.shields.io/github/stars/codefluence-x/altair" alt="GitHub Stars"></a>
  <a href="https://github.com/codefluence-x/altair/blob/master/LICENSE"><img src="https://img.shields.io/github/license/codefluence-x/altair" alt="GitHub License"></a>
</p>

<br><br>

## Introduction

Open source API gateway written in Go. Created to be lightweight, simple, fast, reliable, cross platform, and programming language agnostic.

## Notice

This software is still in *alpha version*, which may contain several hidden bugs that can cause data loss or unexpected behaviour.

## Architecture Diagram

<br>

![Altair Architecture Diagram](https://user-images.githubusercontent.com/20650401/79699757-a2337d00-82bb-11ea-8103-25e6917545bd.png)

## Documentation

### Plugin API Documentation

[Plugin API Documentation in Postman](https://documenter.getpostman.com/view/3666028/SzmcZJ79?version=latest#b870ae5a-b305-4016-8155-4899af1f26b1)

## How to Use

> TBD

## How to Contribute

### Installation

#### Prerequisites

1. Go >= 1.13
2. MySQL

#### How To

1. Clone this repo
2. `go run altair.go migrate main_database`
3. `go run altair.go run`
4. Read [CONTRIBUTING.md](https://github.com/insomnius/code-geek/blob/master/CONTRIBUTING.md)

## Feature

- [ ] Request Forwarder
  - [x] Route Compiler
  - [x] Route Generator Forwader
  - [x] Downstream Plugins
    - [x] Oauth Token Checking
    - [x] Oauth Scope Checking
    - [ ] Response Caching
  - [ ] Persistent HTTP client implementation
- [x] Metric & Monitoring
  - [x] Prometheus
- [x] Logging
  - [x] Stdout
- [ ] Plugins
  - [ ] Plugin dynamic database migration
  - [ ] Oauth Authorization
    - [x] CRUD Oauth Application
      - [x] Create
      - [x] List
      - [x] One
      - [x] Update
    - [x] Authorization
      - [x] Authorization Code Grant
    - [ ] Access Token
      - [x] Access Token Implicit Request for Confidential Application
      - [x] Access Token Code Grant Flow
      - [ ] Refresh Token
        - [ ] Access Token Refresh Token Flow
        - [ ] Refresh Token Generation
      - [x] Revoke Access Token
  - [ ] Response Caching
    - [ ] Route Config
    - [ ] API for deleting the cache
  - [ ] JWT
