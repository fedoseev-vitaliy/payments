# Payments

Simple test payments service

## Install

Clone git repo locally and run `make start` to start docker container

Or `go get github.com/fedoseev-vitaliy/payments` and run bin `payments server` or `$GOPATH/bin/paymenys server`

## Project filesystem structure

```
payments
├── cmd                          # commands
│   ├── server                   # server command
├── internal                     # project internal sources
│   ├── controller               # controller to handle bussiness logic
│   ├── mocks                    # generated mocks (https://github.com/mockery/mockery)
│   ├── provider                 # providers clients
│   │   ├── apay                 # ApplePay client
│   │   └── gpay                 # GooglePay client
│   ├── server                   # server implementation
│   └── utils                    # utils (e.g. http client)
├── tools                        # indirect import for extenal tools like golangci-lint, mockery
└── vendor                       # vednor folder
```

## Dependencies
Golang version at least 1.14

Docker - `brew install docker`

Makefile - `brew install make`

## Some basic command
To get all available cmds run `make help`

###### Env commands

Start service - `make start`, after that service will be available on `localhost:8080`

Stop service - `make stop`

Restart running service - `make restart`

## Available endpoints
After running `make start` payments service will be available on `localhost:8080`

  
GET /api/v1/payments/urls?productID=<productID to get urls>

For testing purposes the following products will cause diff errors:
- `panic` - will cause panic in service
- `fatal` - endpoint will return response with 500
- `badGoogle` - will fail to get GPay url
- `badApple` - will fail to get ApplePay url
