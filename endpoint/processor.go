package endpoint

import "context"

type Processor interface {
	Pre(context.Context, string, *Endpoint, ...any) (context.Context, bool, *any, error)
	Post(context.Context, string, *any, *error) (context.Context, *any, *error)
}
