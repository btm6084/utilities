package handler

import (
	"context"

	"github.com/btm6084/utilities/metrics"
)

type ContextAware interface {
	Context() context.Context
}

type RecorderAware interface {
	Recorder() metrics.Recorder
}

type ContextRecorder interface {
	ContextAware
	RecorderAware
}
