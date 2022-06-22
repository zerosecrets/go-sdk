# Zero Go SDK

## Short intro.
Go SDK for [Zero](https://tryzero.com). Provides a clear and simple interface for the secrets manager GraphQL API.

## Installation
```bash
go get github.com/zerosecrets/go-sdk`
```

## Usage
Fetch secrets for AWS by passing your `zero` token

```go
package main

import (
	"log"
	"os"

	zero "github.com/zerosecrets/go-sdk"
)

func main() {
	api, err := zero.Zero(os.Getenv("ZERO_TOKEN"), []string{"aws"})

	if err != nil {
		panic(err)
	}

	result, err := api.Fetch()

	if err != nil {
		panic(err)
	}

	log.Println(result) // map[aws:map[secret:value secret2:value2]]
}
```
