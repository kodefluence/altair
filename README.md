<center>
  <h1>
    Altair
  </h1>
  <h2>Distributed, Lightweight and Robust API Gateway</h2>
</center>

<br>

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

Open source API gateway written in Go. Designed to be distributed, lightweight, simple, fast, reliable, cross platform, programming language agnostic and robust - by default.

## Notice

This software is still in *alpha version*, which may contain several hidden bugs that can cause data loss or unexpected behaviour.

## Architecture Diagram

<br>

![Altair Architecture Diagram](https://user-images.githubusercontent.com/20650401/79699757-a2337d00-82bb-11ea-8103-25e6917545bd.png)

## Documentation

### Plugin API Documentation

[Plugin API Documentation in Postman](https://documenter.getpostman.com/view/3666028/SzmcZJ79?version=latest#b870ae5a-b305-4016-8155-4899af1f26b1)

## Docker

We are on [dockerhub](https://hub.docker.com/r/codefluence/altair)! Common implementation of altair docker is to have directory where you store your config and routes folder inside it.

```
config/
routes/
.env
docker-compose.yml
```

The content of docker-compose could be like this:

```yaml
version: "3.8"
services:
  altair:
    image: codefluence/altair:latest
    volumes:
      - ./routes/:/opt/altair/routes/
      - ./config/:/opt/altair/config/
      - ./.env:/opt/altair/.env
    ports:
      - "1304:1304"
    network_mode: host
    env_file: ./.env
```

## How to Use

We recommend you to use Altair using docker-compose like above. But if you want to use the binary instead, you could download the binary from release pages.

## How to Contribute

### Installation

#### Prerequisites

1. Go >= 1.13
2. MySQL

#### How To

1. Clone this repo
2. Create databases schema based on your .env or/and configuration
3. Make sure mysql running
4. `go run altair.go migrate main_database`
5. `go run altair.go run`
6. Read [CONTRIBUTING.md](https://github.com/codefluence-x/altair/blob/master/CONTRIBUTING.md)
