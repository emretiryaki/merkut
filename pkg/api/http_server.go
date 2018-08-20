package api

import (
	"context"
	"net/http"
	"gopkg.in/macaron.v1"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/routing"
	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/emretiryaki/merkut/pkg/registry"
)

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "HTTPServer",
		Instance:     &HTTPServer{},
		InitPriority: registry.High,
	})
}

type HTTPServer struct {
	log           log.Logger
	macaron       *macaron.Macaron
	context       context.Context
	httpSrv       *http.Server

	RouteRegister routing.RouteRegister `inject:""`
	Cfg           *setting.Cfg          `inject:""`
}
func (hs *HTTPServer) Init() error {
	hs.log = log.New("http.server")

	//hs.streamManager = live.NewStreamManager()
	hs.macaron = hs.newMacaron()
	hs.registerRoutes()

	return nil
}

func (hs *HTTPServer) newMacaron() *macaron.Macaron {
	macaron.Env = setting.Env
	m := macaron.New()

	// automatically set HEAD for every GET
	m.SetAutoHead(true)

	return m
}