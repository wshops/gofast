# Gofast
[![Go Report Card](https://goreportcard.com/badge/github.com/wshops/gofast)](https://goreportcard.com/report/github.com/wshops/gofast)

⚡️ Gofast is a HTTP client based on [fasthttp](https://github.com/valyala/fasthttp) with zero memory allocation. 

Automatic struct binding let you focus on entity writing.

## Install

```console
go get -u github.com/wshops/gofast
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/wshops/gofast"
)

type Out struct {
	Hello string `json:"hello"`
}

func main() {
	fast := gofast.New()

	var out Out
	uri := "http://echo.jsontest.com/hello/world"
	err, code := fast.Get(uri, &out, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if code == 200 {
		fmt.Printf("hello %v", out.Hello)
    }
	// hello world
}
```

## Examples

### Send request with body

The default encoding is `JSON` with `application/json` header.

You can also use `map` to bind value, but the worse performance you will get. 
```go
type CreateToken struct {
    ID     string `json:"id"`
    Secret string `json:"secret"`
}

type Token struct {
    Token     string `json:"token"`
    ExpiredAt string `json:"expired_at"`
}

fast := gofast.New()

uri := "https://example.com/api/v1/token"
body := CreateToken{
    ID:     "my-id",
    Secret: "my-secret",
}
var token Token
err, code := fast.Post(uri, &body, &token, nil)
if err != nil {
    log.Fatalln(err)
}
if code != 200 {
	fmt.Printf("token: %v, expired_at: %v", token.Token, token.ExpiredAt)
}
```

### Get with header

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

fast := gofast.New()

var user User
uri := "https://example.com/api/v1/users/100"
h := gofast.Header{fasthttp.HeaderAuthorization: "Bearer My-JWT"}
if err,_ := fast.Get(uri, &user, h); err != nil {
    log.Fatalln(err)
}
fmt.Printf("id: %v, name: %v", user.ID, user.Name)
```

### URL encode

Post body with `application/x-www-form-urlencoded` header and get text.

```go
fast := gofast.New(gofast.Config{
    RequestEncoder:  gofast.URLEncoder,
    ResponseDecoder: gofast.TextDecoder,
})

uri := "https://example.com/api/v1/token"
body := gofast.Body{
    "id":     "my-id",
    "secret": "my-secret",
}
var token string
if err,_ := fast.Post(uri, body, &token, nil); err != nil {
    log.Fatalln(err)
}
```

### Customize error handler (*deprecated*)
Error handler will handle non 2xx HTTP status code.

```go
cfg := gofast.Config{
    ErrorHandler: func(resp *fasthttp.Response) error {
        return fmt.Errorf("http code = %d", resp.StatusCode())
    },
}

fast := gofast.New(cfg)
err := fast.Get(uri, nil, nil)
// http code = 400
```

## Benchmarks

```console
$ go test -bench=. -benchmem -benchtime=10s -run=none -cpu 10
BenchmarkPostJSON-10         	 8168836	      1440 ns/op	       0 B/op	       0 allocs/op
BenchmarkPostURLEncode-10    	 8357803	      1443 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/wshops/gofast	26.770s
```
