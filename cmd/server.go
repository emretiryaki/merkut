package main

import (
	"golang.org/x/sync/errgroup"
	"context"
	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/routing"
	"github.com/emretiryaki/merkut/pkg/api"
	"fmt"
	"github.com/emretiryaki/merkut/pkg/login"
)

type MerkutServerImpl struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	log                log.Logger
	cfg                *setting.Cfg
	shutdownReason     string
	shutdownInProgress bool
	RouteRegister routing.RouteRegister `inject:""`
	HttpServer    *api.HTTPServer       `inject:""`
}

func NewMerkutServer() *MerkutServerImpl {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &MerkutServerImpl{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		log:           log.New("server"),
		cfg:           setting.NewCfg(),

	}
}

func (m *MerkutServerImpl) Shutdown(reason string)  {
	m.log.Info("Shutdown started", "reason", reason)
	m.shutdownReason = reason
	m.shutdownInProgress = true

	m.shutdownFn()

	m.childRoutines.Wait()

}
func (m *MerkutServerImpl) Run()  error{

	m.loadConfiguration()
	login.Init()

}

func (g *MerkutServerImpl) Exit(reason error) int {
	// default exit code is 1
	code := 1

	if reason == context.Canceled && g.shutdownReason != "" {
		reason = fmt.Errorf(g.shutdownReason)
		code = 0
	}

	g.log.Error("Server shutdown", "reason", reason)
	return code
}