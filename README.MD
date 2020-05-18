# Deployer  
[![Go Report Card](https://goreportcard.com/badge/github.com/lenniDespero/go-cd)](https://goreportcard.com/report/github.com/lenniDespero/go-cd)
[![Build Status](https://travis-ci.org/lenniDespero/go-cd.svg?branch=master)](https://travis-ci.org/lenniDespero/go-cd)
[![codecov](https://codecov.io/gh/lenniDespero/go-cd/branch/master/graph/badge.svg)](https://codecov.io/gh/lenniDespero/go-cd)  

Simple golang deploy script for web-apps(like php).
It will get source code from git, run work pipes from config and make symlink to current deploy folder.

## Installation
You can:
- `go get github.com/lenniDespero/go-cd`
- Build from source `GOOS={target os} GOARCH={target arch} go build -o ./bin/deployer ./cmd/main.go` 
- Use make for build
- Take prebuild versions from bin/  

Then you can move script to `/usr/local/bin/` for easy usage.  

## Usage
By default, Deployer will take deploy.yml from folder where it run. If you want to use another config - use `--config path/to/config.yml`  
#### Default usage  
```sh
$deployer -target prod
```
Deployer will test config file, and then will deploy app for target 'prod' from config  
#### Test config
```sh
$deployer -test
```
Deployer will only test config file and exit.  

### Config file example
```yaml
projectName: project
git: "some.git.url.git"
count: 2 #count releases to keep in release path
targets:
  devel:
    type: local # type of usage (local or host(ssh))
    path: "testdata/releases" 
    pipe:
      - name: make links
        type: links # type of pipe (links or command)
        args:
          - from: "pathFrom"
            to: "PathTo"
          - from: "pathFrom2"
            to: "PathTo2"
      - name: run cmd
        type: command
        args:
          - command: "ls"
            args:
              - "-la"
  test:
      type: host
      host:
        host: prod.host.ru
        auth: key
        user: user
        key: "path/to/key"
      path: "/var/www/releases/"
      pipe:
        - name: make some work
          type: command
          args:
            - command: "composer"
              args:
                - "install"
        - name: make links
          type: links
          args:
            - from: "pathFrom"
              to: "PathTo"
            - from: "pathFrom2"
              to: "PathTo2"
  prod:
    type: host
    host:
      host: prod.host.ru
      auth: password
      user: user
      password: pass
    path: "/var/www/releases/"
    pipe:
      - name: make some work
        type: command
        args:
          - command: "composer"
            args:
              - "install"
      - name: make links
        type: links
        args:
          - from: "pathFrom"
            to: "PathTo"
          - from: "pathFrom2"
            to: "PathTo2"
```