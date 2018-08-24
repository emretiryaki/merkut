package api

import (
	"context"
	"net/http"
	"fmt"
	"time"
	"path"

	"gopkg.in/macaron.v1"

	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/routing"
	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/emretiryaki/merkut/pkg/registry"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/middleware"
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
	Bus           bus.Bus               `inject:""`
	Cfg           *setting.Cfg          `inject:""`
}
func (hs *HTTPServer) Init() error {
	hs.log = log.New("http.server")
	hs.macaron = hs.newMacaron()
	hs.registerRoutes()

	return nil
}

func (hs *HTTPServer) newMacaron() *macaron.Macaron {

	macaron.Env = setting.Env
	m := macaron.New()
	m.SetAutoHead(true)

	return m
}


func (hs *HTTPServer) Run(ctx context.Context) error {

	var err error

	hs.context = ctx

	hs.applyRoutes()

	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	hs.log.Info("HTTP Server Listen", "address", listenAddr, "protocol", setting.Protocol, "subUrl", setting.AppSubUrl, "socket", setting.SocketPath)

	hs.httpSrv = &http.Server{Addr: listenAddr, Handler: hs.macaron}

	go func() {
		<-ctx.Done()
		// Hacky fix for race condition between ListenAndServe and Shutdown
		time.Sleep(time.Millisecond * 100)
		if err := hs.httpSrv.Shutdown(context.Background()); err != nil {
			hs.log.Error("Failed to shutdown server", "error", err)
		}
	}()

	err = hs.httpSrv.ListenAndServe()
	if err == http.ErrServerClosed {
		hs.log.Debug("server was shutdown")
	}
	return err

}

func (hs *HTTPServer) applyRoutes() {
	// start with middlewares & static routes
	hs.addMiddlewaresAndStaticRoutes()
	// then add view routes & api routes
	hs.RouteRegister.Register(hs.macaron)
	// then custom app proxy routes
	hs.macaron.NotFound(NotFoundHandler)
}

func (hs *HTTPServer) addMiddlewaresAndStaticRoutes(){
	m := hs.macaron


	m.Use(middleware.Recovery())


	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(setting.StaticRootPath, "views"),
		IndentJSON: macaron.Env != macaron.PROD,
		Delims:     macaron.Delims{Left: "[[", Right: "]]"},
	}))
	m.Use(middleware.GetContextHandler())
	m.Use(middleware.AddDefaultResponseHeaders())

}
