package rendering

import (
	"context"
	"errors"
	"time"
)

var ErrTimeout = errors.New("Timeout error. You can set timeout in seconds with &timeout url parameter")
var ErrNoRenderer = errors.New("No renderer plugin found nor is an external render server configured")
var ErrPhantomJSNotInstalled = errors.New("PhantomJS executable not found")

type Opts struct {
	Width    int
	Height   int
	Timeout  time.Duration
	OrgId    int64
	UserId   int64
	Path     string
	Encoding string
	Timezone string
}

type RenderResult struct {
	FilePath string
}

type renderFunc func(ctx context.Context, options Opts) (*RenderResult, error)

type Service interface {
	Render(ctx context.Context, opts Opts) (*RenderResult, error)
}

