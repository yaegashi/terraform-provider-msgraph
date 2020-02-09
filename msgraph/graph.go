package msgraph

import (
	"context"

	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

type graph struct {
	ctx context.Context
	cli *msgraph.GraphServiceRequestBuilder
}

func newGraph(m interface{}) *graph {
	return &graph{
		ctx: context.Background(),
		cli: m.(*msgraph.GraphServiceRequestBuilder),
	}
}
