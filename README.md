# Go Proxy

Simple proxy service implementation in Go

## Features

- request / response forward inc. headers to specified targets based on URL path e.g. `"/test": "https://google.com"`
- header stamp `X-Proxied-By` with the proxy service name as value
- JSON configuration

## Usage

Requires **Go 1.22+** & **make** (WSL recommended)

Simply run `make` to cleanup, test, build and run the service or run the steps individually.
