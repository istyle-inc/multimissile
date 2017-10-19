# Specification for multimissile

multimissile is RESTful API gateway with JSON-RPC. It accepts a HTTP request based JSON-RPC.

## API

multimissile has the APIs below.

 * [POST /msl](#post-msl)
 * [GET /stat/go](#get-statgo)

### POST /msl

Accepts a HTTP request based JSON-RPC and proxies each converted HTTP request to a corresponding endpoint.
And User-Agent and X-Forwarded-For in a request header are forwarded.

The JSON below is a request-body example.

```json
[
  {
    "jsonrpc": "2.0",
    "endpoint": "ep-1",
    "path": "/user/get",
    "params": {
      "user_id": 1
    },
    "id": "1"
  },
  {
    "jsonrpc": "2.0",
    "endpoint": "ep-1",
    "http_method": "GET",
    "path": "/item/get",
    "params": {
      "item_id": 2
    },
    "id": "2"
  },
  {
    "jsonrpc": "2.0",
    "ep": "ep-2",
    "http_method": "POST",
    "path": "/item/update",
    "params": {
      "item_id": 2,
      "desc": "update"
    },
    "id": "3"
  }
]
```

The definitions of parameters are below.

|name            |type  |description                              |required|note                              |
|----------------|------|-----------------------------------------|--------|----------------------------------|
|jsonrpc         |string|version number of JSON-RPC               |o       |fixed as 2.0                      |
|ep              |string|endpoint name                            |o       |selected in Endpoints Section     |
|http_method     |string|method string for HTTP                   |o       |HTTP method string. GET by default|
|method          |string|method string                            |o       |URI                               |
|params          |object|parameters for method                    |o       |                                  |
|id              |string|ID string                                |o       |                                  |


When multimissile receives an invalid request(for example, malformed body is included), a status of response it returns is 400(Bad Request).

### GET /stat/go

Returns a statictics for golang-runtime. See [golang-stats-api-handler](https://github.com/fukata/golang-stats-api-handler) about details.
