package endpoint

import "context"

type Processor interface {
	Pre(context.Context, string, *Endpoint, *any, ...any) (context.Context, bool, error)
	Post(context.Context, string, *any, *error) context.Context
}
