package app

import "context"

var Empty = None{}

type None struct{}

type Interactor[I, O any] interface {
	Execute(ctx context.Context, input I) (output O, err error)
}

type Wrap[I, O any] func(next Interactor[I, O]) Interactor[I, O]
