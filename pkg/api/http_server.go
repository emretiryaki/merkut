package api

import (
	"context"
	"net/http"

	"gopkg.in/macaron.v1"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/routing"
	"github.com/emretiryaki/merkut/pkg/setting"

)

type HTTPServer struct {
	log           log.Logger
	macaron       *macaron.Macaron
	context       context.Context
	httpSrv       *http.Server

	RouteRegister routing.RouteRegister `inject:""`
	Cfg           *setting.Cfg          `inject:""`
}
