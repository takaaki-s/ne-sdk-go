# NextEngine SDK for Go

ne-sdk-go is the NextEngine SDK for the Go programming language.

# Getting Started

## Install

```shell
go get github.com/takaaki-s/ne-sdk-go
```

To update the SDK

```shell
go get -u github.com/takaaki-s/ne-sdk-go
```

## Usage

Please see to https://developer.next-engine.com/

# Example

## Quick Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/takaaki-s/ne-sdk-go/nextengine"
)

func main() {
	nc := nextengine.NewDefaultClient(
		"<CLIENT_ID>",
		"<CLIENT_SECRET>",
		"<REDIRECT_URI>",
		"<ACCESS_TOKEN>",
		"<REFRESH_TOKEN>")
	ctx := context.Background()
	res, err := nc.APIExecute(ctx, "/api_v1_receiveorder_base/count", map[string]string{"receive_order_id-gte": "1"})
	if err != nil {
		fmt.Printf("error: %#v", err)
		return
	}
	fmt.Printf("response: %v", res.Count)
}

```

## Demo Application

https://github.com/takaaki-s/nextengine-go-example

# Lisence

MIT